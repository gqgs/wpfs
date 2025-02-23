package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

type ApiResponse struct {
	Data []struct {
		Path string `json:"path"`
	} `json:"data"`
}

func randomImageHandler(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("new request")
		searchResp, err := http.Get("https://wallhaven.cc/api/v1/search?apikey=" + apiKey + "&sort=random&resolutions=3840x2160&ratios=16x9&categories=010&purity=100&seed=" + fmt.Sprint(time.Now().Unix()))
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

	http.HandleFunc("GET /random", randomImageHandler(opts.apiKey))

	log.Printf("Server is running. Visit http://localhost:%d/random", opts.port)

	return http.ListenAndServe(":"+fmt.Sprint(opts.port), nil)
}
