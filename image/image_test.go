package image_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
			img, err := image.NewImageFromBytes(tt.data)
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
		// Uncomment this test to test against a large, remote image
		//{
		//	name:   "Wikipedia Image",
		//	imgURL: "http://upload.wikimedia.org/wikipedia/commons/9/9a/SKA_dishes_big.jpg",
		//	expected: image.Metadata{
		//		Width:       5000,
		//		Height:      2813,
		//		ContentSize: 10001439,
		//		ContentType: "image/jpeg",
		//	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := image.NewImageFromURL(tt.imgURL)
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

func TestGenerateUniqueFilename(t *testing.T) {
	testCases := []struct {
		name      string
		image     image.Image
		nonEmpty  bool // Whether the returned filename should be non-empty
		hasSuffix bool // Whether the filename should have the provided extension as suffix
		suffix    string
	}{
		{
			name:      "With URL and extension",
			image:     image.Image{URL: "https://example.com/image.jpg", Extension: "jpg"},
			nonEmpty:  true,
			hasSuffix: true,
			suffix:    ".jpg",
		},
		{
			name:      "With URL but no extension",
			image:     image.Image{URL: "https://example.com/image", Extension: ""},
			nonEmpty:  true,
			hasSuffix: false,
			suffix:    "",
		},
		{
			name:      "Without URL but with extension",
			image:     image.Image{URL: "", Extension: "jpg"},
			nonEmpty:  true,
			hasSuffix: true,
			suffix:    ".jpg",
		},
		{
			name:      "Without URL and extension",
			image:     image.Image{URL: "", Extension: ""},
			nonEmpty:  true,
			hasSuffix: false,
			suffix:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filename := tc.image.GenerateUniqueFilename()

			assert.NotEmpty(t, filename, "The filename should not be empty")

			if tc.hasSuffix {
				assert.True(t, strings.HasSuffix(filename, tc.suffix), "The filename should have the correct suffix")
				// Remove suffix for further tests
				filename = strings.TrimSuffix(filename, tc.suffix)
			} else {
				assert.False(t, strings.HasSuffix(filename, "."), "The filename should not have a period as the last character")
			}

			// Test the length of the filename based on whether URL is empty
			if tc.image.URL == "" {
				assert.Equal(t, 32, len(filename), "Random hash should be 32 characters long") // 16 bytes encoded as hex
			} else {
				assert.Equal(t, 40, len(filename), "SHA-1 hash should be 40 characters long")
			}

			// Test that filename only contains valid characters (0-9, a-f for hash)
			assert.Regexp(t, "^[a-f0-9]+$", filename, "The filename should only contain valid characters")
		})
	}

	// Test for randomness (uniqueness) by generating multiple filenames for an empty URL
	t.Run("Randomness", func(t *testing.T) {
		filenames := make(map[string]bool)
		for i := 0; i < 1000; i++ {
			img := &image.Image{URL: "", Extension: ""}
			filename := img.GenerateUniqueFilename()
			_, exists := filenames[filename]
			assert.False(t, exists, "Filename should be unique")
			filenames[filename] = true
		}
	})
}
