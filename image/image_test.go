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

func TestImageFromBytes(t *testing.T) {
	tests := []struct {
		data     []byte
		expected image.Metadata
	}{
		{
			data: imgData["test_image.bmp"],
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.bmp"])),
				ContentType: "image/bmp",
			},
		},
		{
			data: imgData["test_image.gif"],
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.gif"])),
				ContentType: "image/gif",
			},
		},
		{
			data: imgData["test_image.jpeg"],
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.jpeg"])),
				ContentType: "image/jpeg",
			},
		},
		{
			data: imgData["test_image.png"],
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.png"])),
				ContentType: "image/png",
			},
		},
		{
			data: imgData["test_image.tiff"],
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.tiff"])),
				ContentType: "application/octet-stream",
			},
		},
		{
			data: imgData["test_image.webp"],
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.webp"])),
				ContentType: "image/webp",
			},
		},
		{
			data: imgData["test_image.ico"],
			expected: image.Metadata{
				Width:       48,
				Height:      48,
				ContentSize: int64(len(imgData["test_image.ico"])),
				ContentType: "image/x-icon",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			img, err := image.ImageFromBytes(tt.data)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, img.Metadata)
		})
	}
}

func TestImageFromURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/invalid" {
			http.NotFound(w, r)
			return
		}

		if r.URL.Path == "/favicon.ico" {
			reader := bytes.NewReader(imgData["test_image.ico"])
			_, err := io.Copy(w, reader)
			assert.NoError(t, err)
			return
		}

		if r.URL.Path == "/image.webb" {
			reader := bytes.NewReader(imgData["test_image.webp"])
			_, err := io.Copy(w, reader)
			assert.NoError(t, err)
			return
		}

		if r.URL.Path == "/image.jpg" {
			reader := bytes.NewReader(imgData["test_image.jpeg"])
			_, err := io.Copy(w, reader)
			assert.NoError(t, err)
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
			name:   "valid png image",
			imgURL: server.URL,
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.png"])),
				ContentType: "image/png",
			},
		},
		{
			name:   "valid jpeg image",
			imgURL: server.URL + "/image.jpg",
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.jpeg"])),
				ContentType: "image/jpeg",
			},
		},
		{
			name:   "valid webp image",
			imgURL: server.URL + "/image.webb",
			expected: image.Metadata{
				Width:       50,
				Height:      50,
				ContentSize: int64(len(imgData["test_image.webp"])),
				ContentType: "image/webp",
			},
		},
		{
			name:   "valid favicon",
			imgURL: server.URL + "/favicon.ico",
			expected: image.Metadata{
				Width:       48,
				Height:      48,
				ContentSize: int64(len(imgData["test_image.ico"])),
				ContentType: "image/x-icon",
			},
		},
		{
			name:        "invalid image",
			imgURL:      server.URL + "/invalid",
			expectedErr: "image: unknown format",
		},
		{
			name:   "Wikipedia Image",
			imgURL: "http://upload.wikimedia.org/wikipedia/commons/9/9a/SKA_dishes_big.jpg",
			expected: image.Metadata{
				Width:       5000,
				Height:      2813,
				ContentSize: 10001439,
				ContentType: "image/jpeg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := image.ImageFromURL(tt.imgURL)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, img.Metadata)

				// assert img is not nil and is of type image.Image
				assert.NotNil(t, img)
			}
		})
	}
}
