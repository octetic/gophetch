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
		name     string
		data     []byte
		expected image.Metadata
	}{
		{
			name: "valid bmp image",
			data: imgData["mark.bmp"],
			expected: image.Metadata{
				Width:       100,
				Height:      100,
				ContentSize: int64(len(imgData["mark.bmp"])),
				ContentType: "image/bmp",
			},
		},
		{
			name: "valid gif image",
			data: imgData["mark.gif"],
			expected: image.Metadata{
				Width:       100,
				Height:      100,
				ContentSize: int64(len(imgData["mark.gif"])),
				ContentType: "image/gif",
			},
		},
		{
			name: "valid jpg image",
			data: imgData["mark.jpg"],
			expected: image.Metadata{
				Width:       100,
				Height:      100,
				ContentSize: int64(len(imgData["mark.jpg"])),
				ContentType: "image/jpeg",
			},
		},
		{
			name: "valid png image",
			data: imgData["mark.png"],
			expected: image.Metadata{
				Width:       100,
				Height:      100,
				ContentSize: int64(len(imgData["mark.png"])),
				ContentType: "image/png",
			},
		},
		{
			name: "valid tiff image",
			data: imgData["mark.tif"],
			expected: image.Metadata{
				Width:       100,
				Height:      100,
				ContentSize: int64(len(imgData["mark.tif"])),
				ContentType: "application/octet-stream",
			},
		},
		{
			name: "valid webp image",
			data: imgData["mark.webp"],
			expected: image.Metadata{
				Width:       100,
				Height:      100,
				ContentSize: int64(len(imgData["mark.webp"])),
				ContentType: "image/webp",
			},
		},
		{
			name: "valid favicon",
			data: imgData["mark.ico"],
			expected: image.Metadata{
				Width:       48,
				Height:      48,
				ContentSize: int64(len(imgData["mark.ico"])),
				ContentType: "image/x-icon",
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			img, err := image.NewImageFromBytes(tt.data)
			assert.NoError(t, err, tt.name+" error mismatch")
			assert.Equal(t, tt.expected, img.Metadata, tt.name+" metadata mismatch")
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
			reader := bytes.NewReader(imgData["mark.ico"])
			_, err := io.Copy(w, reader)
			assert.NoError(t, err)
			return
		}

		if r.URL.Path == "/image.webb" {
			reader := bytes.NewReader(imgData["mark.webp"])
			_, err := io.Copy(w, reader)
			assert.NoError(t, err)
			return
		}

		if r.URL.Path == "/image.jpg" {
			reader := bytes.NewReader(imgData["mark.jpg"])
			_, err := io.Copy(w, reader)
			assert.NoError(t, err)
			return
		}

		reader := bytes.NewReader(imgData["mark.png"])
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
				Width:       100,
				Height:      100,
				ContentSize: int64(len(imgData["mark.png"])),
				ContentType: "image/png",
			},
		},
		{
			name:   "valid jpeg image",
			imgURL: server.URL + "/image.jpg",
			expected: image.Metadata{
				Width:       100,
				Height:      100,
				ContentSize: int64(len(imgData["mark.jpg"])),
				ContentType: "image/jpeg",
			},
		},
		{
			name:   "valid webp image",
			imgURL: server.URL + "/image.webb",
			expected: image.Metadata{
				Width:       100,
				Height:      100,
				ContentSize: int64(len(imgData["mark.webp"])),
				ContentType: "image/webp",
			},
		},
		{
			name:   "valid favicon",
			imgURL: server.URL + "/favicon.ico",
			expected: image.Metadata{
				Width:       48,
				Height:      48,
				ContentSize: int64(len(imgData["mark.ico"])),
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
			img, err := image.NewImageFromURL(tt.imgURL, 0)
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
				assert.True(t, strings.HasSuffix(filename, tc.suffix), tc.name+"The filename should have the correct suffix")
				// Remove suffix for further tests
				filename = strings.TrimSuffix(filename, tc.suffix)
			} else {
				assert.False(t, strings.HasSuffix(filename, "."), tc.name+"The filename should not have a period as the last character")
			}

			// Test the length of the filename based on whether URL is empty
			if tc.image.URL == "" {
				assert.Equal(t, 32, len(filename), tc.name+"Random hash should be 32 characters long") // 16 bytes encoded as hex
			} else {
				assert.Equal(t, 40, len(filename), tc.name+"SHA-1 hash should be 40 characters long")
			}

			// Test that filename only contains valid characters (0-9, a-f for hash)
			assert.Regexp(t, "^[a-f0-9]+$", filename, tc.name+"The filename should only contain valid characters")
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

func TestNewImageFromDataURI(t *testing.T) {
	tests := []struct {
		name                string
		inputFile           string
		input64             string
		expectedContentType string
		expectedWidth       int
		expectedHeight      int
		expectErr           bool
	}{
		{
			name:                "ValidDataURI PNG",
			inputFile:           "mark.png",
			expectedContentType: "image/png",
			expectedWidth:       100,
			expectedHeight:      100,
			expectErr:           false,
		},
		{
			name:                "ValidDataURI JPG",
			inputFile:           "mark.jpg",
			expectedContentType: "image/jpeg",
			expectedWidth:       100,
			expectedHeight:      100,
			expectErr:           false,
		},
		{
			name:                "ValidDataURI WEBP",
			inputFile:           "mark.webp",
			expectedContentType: "image/webp",
			expectedWidth:       100,
			expectedHeight:      100,
			expectErr:           false,
		},
		{
			name:                "ValidDataURI ICO",
			inputFile:           "mark.ico",
			expectedContentType: "image/x-icon",
			expectedWidth:       48,
			expectedHeight:      48,
			expectErr:           false,
		},
		{
			name:      "InvalidDataURI",
			input64:   "invalidData",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base64Str := tt.input64
			if tt.inputFile != "" {
				base64Str, _ = ReadAndEncodeImage(tt.inputFile)
			}
			img, err := image.NewImageFromDataURI(base64Str)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, img)
				if tt.expectedContentType != "" {
					assert.Equal(t, tt.expectedContentType, img.ContentType)
				}
				if tt.expectedWidth != 0 {
					assert.Equal(t, tt.expectedWidth, img.Width)
				}
				if tt.expectedHeight != 0 {
					assert.Equal(t, tt.expectedHeight, img.Height)
				}
			}
		})
	}
}

func TestDataURIToBytes(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		input64   string
		expectErr bool
	}{
		{
			name:      "ValidDataURI PNG",
			inputFile: "mark.png",
			expectErr: false,
		},
		{
			name:      "ValidDataURI JPG",
			inputFile: "mark.jpg",
			expectErr: false,
		},
		{
			name:      "ValidDataURI WEBP",
			inputFile: "mark.webp",
			expectErr: false,
		},
		{
			name:      "InvalidDataURI",
			inputFile: "",
			input64:   "invalidData",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			base64Str := tt.input64
			if tt.inputFile != "" {
				base64Str, err = ReadAndEncodeImage(tt.inputFile)
				if err != nil {
					t.Fatalf("Error reading file: %v", err)
				}
			}
			uriBytes, err := image.DataURIToBytes(base64Str)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, uriBytes)
				// Add more assertions based on what DataURIToBytes is supposed to do
			}
		})
	}
}
