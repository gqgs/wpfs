package main

import (
	"log/slog"
	"os"
)

//go:generate go tool argsgen

type options struct {
	root string `arg:"root of folder to be server,required"`
	port int    `arg:"server port,required"`
}

func main() {
	o := options{
		root: os.Getenv("WPFS_ROOT"),
		port: 9999,
	}
	o.MustParse()

	if err := handler(o); err != nil {
		slog.Error(err.Error())
	}
}
