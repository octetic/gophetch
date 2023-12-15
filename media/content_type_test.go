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

// assuming these functions are defined in the same package for simplicity
// if defined in another package, make sure to import that package

func TestContentTypeFunctions(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		extension      string
		extensionAlt   string
		isValidImage   bool
		isValidFavicon bool
	}{
		{
			name:           "JPG Test",
			contentType:    "image/jpeg",
			extension:      ".jpg",
			isValidImage:   true,
			isValidFavicon: true,
		},
		{
			name:           "JPEG Test",
			contentType:    "image/jpeg",
			extension:      ".jpg",
			extensionAlt:   ".jpeg",
			isValidImage:   true,
			isValidFavicon: true,
		},
		{
			name:           "PNG Test",
			contentType:    "image/png",
			extension:      ".png",
			isValidImage:   true,
			isValidFavicon: true,
		},
		{
			name:           "Invalid Content Type",
			contentType:    "invalid/content-type",
			extension:      ".invalid",
			isValidImage:   false,
			isValidFavicon: false,
		},
		{
			name:           "SVG Test",
			contentType:    "image/svg+xml",
			extension:      ".svg",
			isValidImage:   true,
			isValidFavicon: true,
		},
		{
			name:           "ICO Test",
			contentType:    "image/x-icon",
			extension:      ".ico",
			isValidImage:   true,
			isValidFavicon: true,
		},
		{
			name:           "TIFF Test",
			contentType:    "image/tiff",
			extension:      ".tif",
			isValidImage:   true,
			isValidFavicon: false,
		},
		{
			name:           "BMP Test",
			contentType:    "image/bmp",
			extension:      ".bmp",
			isValidImage:   true,
			isValidFavicon: true,
		},
		// ... add more cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			// Test IsValidImageContentType
			a.Equal(tt.isValidImage, media.IsValidImageContentType(tt.contentType), "IsValidImageContentType mismatch")

			// Test IsValidFaviconContentType
			a.Equal(tt.isValidFavicon, media.IsValidFaviconContentType(tt.contentType), "IsValidFaviconContentType mismatch")

			// Test ContentTypeByExtension
			extToTest := tt.extension
			if tt.extensionAlt != "" {
				extToTest = tt.extensionAlt
			}
			contentType, err := media.ContentTypeByExtension(extToTest)
			if tt.isValidImage {
				a.NoError(err, "ContentTypeByExtension should not return an error for valid extension")
				a.Equal(tt.contentType, contentType, "ContentTypeByExtension mismatch")
			} else {
				a.Error(err, "ContentTypeByExtension should return an error for invalid extension")
			}

			// Test ExtensionByContentType
			extension, err := media.ExtensionByContentType(tt.contentType)
			if tt.isValidImage {
				a.NoError(err, "ExtensionByContentType should not return an error for valid content-type")
				a.Equal(tt.extension, extension, "ExtensionByContentType mismatch")
			} else {
				a.Error(err, "ExtensionByContentType should return an error for invalid content-type")
			}
		})
	}
}

func TestContentTypeFromURL(t *testing.T) {
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
		name        string
		imgURL      string
		expected    string
		expectedErr string
	}{
		{
			name:     "valid image: tests the fallback when a HEAD request fails",
			imgURL:   server.URL,
			expected: "image/png",
		},
		{
			name:        "invalid image",
			imgURL:      server.URL + "/invalid",
			expectedErr: "could not extract metadata within the first 512 KB",
		},
		{
			name:     "Place Kitten",
			imgURL:   "https://placekitten.com/g/2048/2048",
			expected: "image/jpeg",
		},
		{
			name:     "Wikipedia Image",
			imgURL:   "http://upload.wikimedia.org/wikipedia/commons/9/9a/SKA_dishes_big.jpg",
			expected: "image/jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentType, err := media.ContentTypeFromURL(tt.imgURL)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, contentType)
			}
		})
	}
}
