package downloader

import (
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestNewDownloader(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		expectedTTL   time.Duration
		expectedError error
	}{
		{
			name: "valid REDIS_ADDR and TTL",
			envVars: map[string]string{
				"REDIS_ADDR": "localhost:6379",
				"TTL":        "30s",
			},
			expectedTTL: 30 * time.Second,
		},
		{
			name: "invalid TTL",
			envVars: map[string]string{
				"REDIS_ADDR": "localhost:6379",
				"TTL":        " invalid",
			},
			expectedTTL: 10 * time.Second,
		},
		{
			name: "missing TTL",
			envVars: map[string]string{
				"REDIS_ADDR": "localhost:6379",
			},
			expectedTTL: 10 * time.Second,
		},
		{
			name: "missing REDIS_ADDR",
			envVars: map[string]string{
				"TTL": "30s",
			},
			expectedError: redis.Nil,
			expectedTTL:   30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Reset environment variables
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			// Create a new downloader
			downloader := NewDownloader()

			// Check TTL
			if downloader.ttl != tt.expectedTTL {
				t.Errorf("expected TTL %v, got %v", tt.expectedTTL, downloader.ttl)
			}

			// Check error
			if tt.expectedError != nil {
				if downloader.redclient == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				}
			} else {
				if downloader.redclient == nil {
					t.Errorf("expected non-nil redis client, got nil")
				}
			}
		})
	}
}

func TestGetVideoID(t *testing.T) {
	tests := []struct {
		name    string
		vidUrl  string
		wantID  string
		wantErr bool
	}{
		{
			name:    "Valid YouTube URL with query parameter",
			vidUrl:  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			wantID:  "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "Valid YouTube URL with path",
			vidUrl:  "https://www.youtube.com/dQw4w9WgXcQ",
			wantID:  "dQw4w9WgXcQ",
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			vidUrl:  "invalid url",
			wantID:  "",
			wantErr: true,
		},
		{
			name:    "URL without video ID",
			vidUrl:  "https://www.youtube.com",
			wantID:  "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := getVideoID(tt.vidUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVideoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("getVideoID() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}
