package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type ApiResponse struct {
	Data []struct {
		Path string `json:"path"`
	} `json:"data"`
}

func randomImageHandler(opts options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := make(url.Values, 0)
		query.Add("apikey", opts.apiKey)
		query.Add("sort", "random")
		query.Add("resolutions", opts.resolution)
		query.Add("ratios", opts.ratio)
		query.Add("categories", opts.category)
		query.Add("purity", opts.purity)
		query.Add("seed", fmt.Sprint(time.Now().Unix()))

		url := "https://wallhaven.cc/api/v1/search?" + query.Encode()
		slog.Info("new request", "url", url)

		searchResp, err := http.Get(url)
		if err != nil {
			slog.Error("failed to get random image", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer searchResp.Body.Close()

		var apiResponse ApiResponse
		if err := json.NewDecoder(searchResp.Body).Decode(&apiResponse); err != nil {
			slog.Error("failed to decode random image", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(apiResponse.Data) == 0 {
			slog.Error("no data found")
			http.Error(w, "no data found", http.StatusInternalServerError)
			return
		}

		path := apiResponse.Data[rand.Intn(len(apiResponse.Data))].Path
		slog.Info("random image", "path", path)

		imageResp, err := http.Get(path)
		if err != nil {
			slog.Error("failed to get random image", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer imageResp.Body.Close()

		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		if _, err = io.Copy(w, imageResp.Body); err != nil {
			slog.Error("failed to copy random image", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handler(opts options) error {

	http.HandleFunc("GET /random", randomImageHandler(opts))

	log.Printf("Server is running. Visit http://localhost:%d/random", opts.port)

	return http.ListenAndServe(":"+fmt.Sprint(opts.port), nil)
}
