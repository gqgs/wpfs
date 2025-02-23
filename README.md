# WPFS (Wallpaper Filesystem)

WPFS is a FUSE-based filesystem that provides dynamic access to wallpapers. It supports multiple backends including local files and wallhaven.cc API integration.

## Features

- FUSE-based filesystem implementation
- Multiple backend support:
  - Local file serving
  - Wallhaven.cc API integration
  - HTTP server for random image serving
- Dynamic file generation

## Requirements

- Go 1.24.0 or later
- FUSE kernel module installed
- Wallhaven.cc API key (if using the wallhaven backend)

## Installation

```bash
go install github.com/gqgs/wpfs@latest
```

## Components

### 1. Filesystem Server (cmd/fs)

The core FUSE filesystem implementation that presents images as files in your filesystem.

Environment variables:
- `WPFS_MOUNTPOINT`: Directory where the filesystem will be mounted
- `WPFS_FILE_SERVER`: Endpoint of the file server to fetch images from

### 2. HTTP Server (cmd/http)

A simple HTTP server that serves random images from a local directory.

Environment variables:
- `WPFS_ROOT`: Root directory containing images to serve
- Default port: 9999

### 3. Wallhaven Integration (cmd/wallhaven)

Integration with wallhaven.cc API to fetch random wallpapers.

Environment variables:
- `WPFS_WALLHAVEN_API_KEY`: Your wallhaven.cc API key
- Default port: 9999
- Default resolution: 3840x2160
- Default aspect ratio: 16x9

## Usage

1. Start the desired backend server:

   For local files:
   ```bash
   export WPFS_ROOT=/path/to/your/wallpapers
   go run ./cmd/http
   ```

   For wallhaven:
   ```bash
   export WPFS_WALLHAVEN_API_KEY=your_api_key
   go run ./cmd/wallhaven
   ```

2. Mount the filesystem:
   ```bash
   export WPFS_MOUNTPOINT=/path/to/mount
   export WPFS_FILE_SERVER=http://localhost:9999
   go run ./cmd/fs
   ```

3. Access your wallpapers through the mount point:
   ```bash
   ls /path/to/mount
   ```
