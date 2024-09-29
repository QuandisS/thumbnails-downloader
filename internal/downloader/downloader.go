package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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

func DowloadThumbnail(vidUrl string) ([]byte, error) {
	vidId, err := getVideoID(vidUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get video id: %w", err)
	}
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", vidId)

	resp, err := http.Get(thumbnailURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download thumbnail: %w", err)
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imgData, nil
}
