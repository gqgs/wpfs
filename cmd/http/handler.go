package main

import (
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
)

func randomImageHandler(rootFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("new request")
		matches, _ := fs.Glob(rootFS, "*")
		randomFile := matches[rand.Intn(len(matches))]
		slog.Info("returning new request", "file", randomFile)
		http.ServeFileFS(w, r, rootFS, randomFile)
	}
}

func handler(opts options) error {
	root, err := os.OpenRoot(opts.root)
	if err != nil {
		return fmt.Errorf("failed to open root: %w", err)
	}

	http.HandleFunc("GET /random", randomImageHandler(root.FS()))

	log.Printf("Server is running. Visit http://localhost:%d/random", opts.port)

	return http.ListenAndServe(":"+fmt.Sprint(opts.port), nil)
}
