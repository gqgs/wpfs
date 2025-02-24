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

	"golang.org/x/sync/singleflight"
)

type ApiResponse struct {
	Data []struct {
		Path string `json:"path"`
	} `json:"data"`
}

type Image struct {
	Body     []byte
	MimeType string
}

func randomImageHandler(opts options) http.HandlerFunc {
	g := new(singleflight.Group)
	return func(w http.ResponseWriter, r *http.Request) {
		seed := time.Now().Truncate(time.Minute).Unix()
		query := make(url.Values, 0)
		query.Add("apikey", opts.apiKey)
		query.Add("sort", "random")
		query.Add("resolutions", opts.resolution)
		query.Add("ratios", opts.ratio)
		query.Add("categories", opts.category)
		query.Add("purity", opts.purity)
		query.Add("seed", fmt.Sprint(seed))

		url := "https://wallhaven.cc/api/v1/search?" + query.Encode()
		slog.Info("new request", "url", url)

		result, err, _ := g.Do(fmt.Sprint(seed), func() (any, error) {
			searchResp, err := http.Get(url)
			if err != nil {
				return nil, fmt.Errorf("failed to get random image: %w", err)
			}
			defer searchResp.Body.Close()

			var apiResponse ApiResponse
			if err := json.NewDecoder(searchResp.Body).Decode(&apiResponse); err != nil {
				return nil, fmt.Errorf("failed to decode random image: %w", err)
			}

			if len(apiResponse.Data) == 0 {
				return nil, fmt.Errorf("no data found")
			}

			path := apiResponse.Data[rand.Intn(len(apiResponse.Data))].Path
			slog.Info("random image", "path", path)

			imageResp, err := http.Get(path)
			if err != nil {
				return nil, fmt.Errorf("failed to get random image: %w", err)
			}
			defer imageResp.Body.Close()

			body, err := io.ReadAll(imageResp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed read body data: %w", err)
			}
			return &Image{
				Body:     body,
				MimeType: imageResp.Header.Get("Content-Type"),
			}, nil
		})

		if err != nil {
			slog.Error("failed to get data from url", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", result.(*Image).MimeType)
		w.WriteHeader(http.StatusOK)
		w.Write(result.(*Image).Body)
	}
}

func handler(opts options) error {

	http.HandleFunc("GET /random", randomImageHandler(opts))

	log.Printf("Server is running. Visit http://localhost:%d/random", opts.port)

	return http.ListenAndServe(":"+fmt.Sprint(opts.port), nil)
}
