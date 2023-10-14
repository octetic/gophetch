package image_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/minsoft-io/gophetch/image"
)

func TestFetchMetadataFromURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/invalid" {
			http.NotFound(w, r)
			return
		}

		reader := bytes.NewReader(imgData["test_image.png"])
		_, err := io.Copy(w, reader)
		assert.NoError(t, err)
	}))
	defer server.Close()

	tests := []struct {
		name        string
		imgURL      string
		expected    image.Metadata
		expectedErr string
	}{
		{
			name:   "valid image",
			imgURL: server.URL,
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.png"])),
				ContentType: "image/png",
			},
		},
		{
			name:        "invalid image",
			imgURL:      server.URL + "/invalid",
			expectedErr: "could not extract metadata within the first 512 KB",
		},
		{
			name:   "Place Kitten",
			imgURL: "https://placekitten.com/g/2048/2048",
			expected: image.Metadata{
				Width:       2048,
				Height:      2048,
				ContentSize: 0,
				ContentType: "image/jpeg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := image.FetchMetadataFromHeader(tt.imgURL, 512*1024)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, metadata.Width)
				assert.NotZero(t, metadata.Height)
				// ... other assertions ...
			}
		})
	}
}
