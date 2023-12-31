package media_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/octetic/gophetch/media"
)

func TestIsValidImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/invalid" {
			http.NotFound(w, r)
			return
		}

		reader := bytes.NewReader(imgData["mark.png"])
		_, err := io.Copy(w, reader)
		assert.NoError(t, err)
	}))

	defer server.Close()

	tests := []struct {
		name     string
		imgURL   string
		expected bool
	}{
		{
			name:     "valid image",
			imgURL:   server.URL,
			expected: true,
		},
		{
			name:     "invalid image",
			imgURL:   "https://www.example.com",
			expected: false,
		},
		{
			name:     "Place Kitten",
			imgURL:   "https://placekitten.com/g/2048/2048",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := media.IsValidImage(tt.imgURL)
			assert.Equal(t, tt.expected, valid)
		})
	}
}

func TestIsValidFavicon(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/invalid" {
			reader := bytes.NewReader(imgData["mark.tif"])
			_, err := io.Copy(w, reader)
			assert.NoError(t, err)
			return
		}

		reader := bytes.NewReader(imgData["mark.ico"])
		_, err := io.Copy(w, reader)
		assert.NoError(t, err)
	}))

	defer server.Close()

	tests := []struct {
		name     string
		imgURL   string
		expected bool
	}{
		{
			name:     "valid favicon",
			imgURL:   server.URL,
			expected: true,
		},
		{
			name:     "invalid image",
			imgURL:   server.URL + "/invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := media.IsValidFavicon(tt.imgURL)
			assert.Equal(t, tt.expected, valid)
		})
	}
}
