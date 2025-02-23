package main

import (
	"log/slog"
	"os"
)

//go:generate go tool argsgen

type options struct {
	port   int    `arg:"server port,required"`
	apiKey string `arg:"auth key to access api,required"`
}

func main() {
	o := options{
		apiKey: os.Getenv("WPFS_WALLHAVEN_API_KEY"),
		port:   9999,
	}
	o.MustParse()

	if err := handler(o); err != nil {
		slog.Error(err.Error())
	}
}
