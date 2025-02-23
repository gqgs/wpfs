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
		slog.Info("new request")
		url, err := url.Parse("https://wallhaven.cc/api/v1/search")
		if err != nil {
			slog.Error("failed to parse url", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		url.Query().Add("apikey", opts.apiKey)
		url.Query().Add("sort", "random")
		url.Query().Add("resolutions", opts.resolution)
		url.Query().Add("ratios", opts.ratio)
		url.Query().Add("categories", opts.category)
		url.Query().Add("purity", opts.purity)
		url.Query().Add("seed", fmt.Sprint(time.Now().Unix()))

		searchResp, err := http.Get(url.String())
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
