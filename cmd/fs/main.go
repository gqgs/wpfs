package main

import (
	"log/slog"
	"os"
)

//go:generate go tool argsgen

type options struct {
	mountpoint string `arg:"mountpoint,required"`
	fileServer string `arg:"file server endpoint,required"`
}

func main() {
	o := options{
		mountpoint: os.Getenv("WPFS_MOUNTPOINT"),
		fileServer: os.Getenv("WPFS_FILE_SERVER"),
	}
	o.MustParse()

	if err := handler(o); err != nil {
		slog.Error(err.Error())
	}
}
