package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

const (
	numberOfFiles = 2
)

var (
	randomEndpoint string
	counter        int64 = 0
)

// File represents an image file node with cached data.
type File struct {
	fs.Inode
	Name string

	mu   sync.Mutex
	data []byte // cached image data; if nil, not yet fetched
}

// Ensure File implements Getattr and NodeOpener.
var _ = (fs.NodeGetattrer)((*File)(nil))
var _ = (fs.NodeOpener)((*File)(nil))
var _ = (fs.NodeReleaser)((*File)(nil))

func (f *File) Release(ctx context.Context, fh fs.FileHandle) syscall.Errno {
	slog.Info("File Release")
	f.mu.Lock()
	defer f.mu.Unlock()
	f.data = nil
	return fs.OK
}

func (f *File) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	slog.Info("Getattr", "file", f.Name)
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.data == nil {
		resp, err := http.Get(randomEndpoint)
		if err != nil {
			return syscall.EIO
		}
		defer resp.Body.Close()
		f.data, err = io.ReadAll(resp.Body)
		if err != nil {
			return syscall.EIO
		}
	}
	out.Mode = fuse.S_IFREG | 0444
	out.Size = uint64(len(f.data))
	return 0
}

func (f *File) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	slog.Info("Open", "file", f.Name)
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.data == nil {
		resp, err := http.Get(randomEndpoint)
		if err != nil {
			return nil, 0, syscall.EIO
		}
		defer resp.Body.Close()

		d, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, syscall.EIO
		}
		f.data = d
	}
	// Make a copy for the file handle.
	dataCopy := make([]byte, len(f.data))
	copy(dataCopy, f.data)
	return &FileHandle{data: dataCopy}, fuse.FOPEN_DIRECT_IO, 0
}

// FileHandle holds the file's content.
type FileHandle struct {
	data []byte
}

// Read returns file data respecting offset and length.
func (fh *FileHandle) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	// slog.Info("Read", "off", off, "len", len(dest))
	if off >= int64(len(fh.data)) {
		return fuse.ReadResultData([]byte{}), 0
	}
	end := int(off) + len(dest)
	if end > len(fh.data) {
		end = len(fh.data)
	}
	return fuse.ReadResultData(fh.data[off:end]), 0
}

// Release is a no-op in this example.
func (fh *FileHandle) Release(ctx context.Context) syscall.Errno {
	slog.Info("FileHandle Release")
	return 0
}

// Dir represents the root directory.
type Dir struct {
	fs.Inode
	files   map[string]*File
	entries []fuse.DirEntry
}

// Readdir returns a list of files in the directory.
func (d *Dir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	slog.Info("Readdir")
	return fs.NewListDirStream(d.entries), 0
}

// Lookup will check the current (latest) generation of file nodes.
func (d *Dir) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	slog.Info("Lookup", "name", name)

	if !strings.HasSuffix(name, ".png") {
		return nil, syscall.ENOENT
	}

	if d.files[name] == nil {
		d.files[name] = &File{Name: name}
	}

	return d.NewInode(ctx, d.files[name], fs.StableAttr{Mode: fuse.S_IFREG}), 0
}

func handler(opts options) error {
	entries := make([]fuse.DirEntry, 0, numberOfFiles)
	for range numberOfFiles {
		entries = append(entries, fuse.DirEntry{
			Name: fmt.Sprintf("%05d.png", atomic.AddInt64(&counter, 1)),
			Mode: fuse.S_IFREG,
		})
	}
	root := &Dir{
		files:   make(map[string]*File, 0),
		entries: entries,
	}

	randomEndpoint = opts.fileServer + "/random"

	server, err := fs.Mount(opts.mountpoint, root, &fs.Options{
		MountOptions: fuse.MountOptions{
			Debug: false,
		},
	})
	if err != nil {
		return fmt.Errorf("mount fail: %w", err)
	}

	slog.Info("mounted at", "mountpoint", opts.mountpoint, "randomEndpoint", randomEndpoint)
	server.Wait()

	return nil
}
