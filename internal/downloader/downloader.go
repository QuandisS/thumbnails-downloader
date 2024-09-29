package downloader

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Downloader struct {
	redclient *redis.Client
	ttl       time.Duration
}

func NewDownloader() *Downloader {
	rclient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	ttl, ok := os.LookupEnv("TTL")
	if !ok {
		ttl = "10"
	}

	t, err := time.ParseDuration(ttl)
	if err != nil {
		t = 10 * time.Second
	}

	log.Println("Cache TTL: ", t)
	return &Downloader{redclient: rclient, ttl: t}
}
func getVideoID(vidUrl string) (string, error) {
	u, err := url.Parse(vidUrl)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	if id := u.Query().Get("v"); id != "" {
		return id, nil
	}

	parts := strings.Split(u.Path, "/")
	if len(parts) > 1 {
		return parts[1], nil
	}

	return "", fmt.Errorf("cant parse video id from url: %s", vidUrl)
}

func (d *Downloader) DowloadThumbnail(vidUrl string) ([]byte, error) {
	vidId, err := getVideoID(vidUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get video id: %w", err)
	}

	imgData, err := d.redclient.Get(context.Background(), vidId).Bytes()
	if err == nil {
		log.Println("Using cached thumbnail for:", vidId)
		return imgData, nil
	}

	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", vidId)

	resp, err := http.Get(thumbnailURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download thumbnail: %w", err)
	}
	defer resp.Body.Close()

	imgData, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	log.Println("Caching thumbnail:", vidId)
	err = d.redclient.Set(context.Background(), vidId, imgData, d.ttl).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to cache thumbnail: %w", err)
	}
	return imgData, nil
}
