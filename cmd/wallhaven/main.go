package main

import (
	"log/slog"
	"os"
)

//go:generate go tool argsgen

type options struct {
	port       int    `arg:"server port,required"`
	apiKey     string `arg:"auth key to access api,required"`
	resolution string `arg:"resolution,required"`
	ratio      string `arg:"ratio,required"`
	category   string `arg:"category,required"`
	purity     string `arg:"purity,required"`
}

func main() {
	o := options{
		apiKey:     os.Getenv("WPFS_WALLHAVEN_API_KEY"),
		port:       9999,
		resolution: "3840x2160",
		ratio:      "16x9",
		category:   "010",
		purity:     "100",
	}
	o.MustParse()

	if err := handler(o); err != nil {
		slog.Error(err.Error())
	}
}
