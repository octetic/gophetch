package gophetch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/minsoft-io/gophetch"
	"github.com/minsoft-io/gophetch/image"
)

func TestHybridStrategy(t *testing.T) {
	SetupFiles()

	tests := []struct {
		name           string
		inputHTML      string
		expectedHTML   string
		expectedFiles  []string
		shouldInline   bool
		mockFetchImage *image.Image
		mockFetchError error
	}{
		{
			name: "images get inlined",
			inputHTML: `<html><head><body>
			<img src="https://example.com/mark.jpg">
			<img srcset="https://example.com/mark.png 1x,https://example.com/mark.webp 2x" sizes="100vw" width="100">
			</body></html>`,
			shouldInline: true,
			expectedHTML: "",
			expectedFiles: []string{
				"mark.png",
				"mark.webp",
			},
		},
		{
			name:         "image gets uploaded",
			inputHTML:    `<img src="mark.png">`,
			expectedHTML: `<html><head></head><body><img src="new_image_url.png"/></body></html>`,
			shouldInline: false,
			mockFetchImage: &image.Image{
				Bytes:    []byte("image data"),
				Metadata: image.Metadata{ContentType: "image/png"},
			},
			mockFetchError: nil,
		},
		// Add more test cases if needed
	}

	mockUploadFunc := func(img *image.Image) (string, error) {
		return "new_image_url.png", nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the behavior of ImageFetcher
			mockFetcher := new(MockImageFetcherHybrid)
			mockFetcher.On("NewImageFromURL", mock.Anything).Return(tt.mockFetchImage, tt.mockFetchError)

			inliner := gophetch.NewImageInliner(gophetch.ImageInlinerOptions{
				Fetcher:    mockFetcher,
				UploadFunc: mockUploadFunc,
				Strategy:   gophetch.StrategyHybrid,
			})
			inliner.ShouldInline = func(img *image.Image) bool {
				return tt.shouldInline
			}

			outputHTML, err := inliner.InlineImages(tt.inputHTML)
			assert.NoError(t, err)

			if tt.shouldInline {
				for _, filename := range tt.expectedFiles {
					expectedBase64, ok := imageBase64Map[filename]
					assert.True(t, ok, "Expected base64 for %s in testdata not found", filename)
					assert.Contains(t, outputHTML, expectedBase64, "Base64 data string for %s not found. Output was: \n %s", filename, outputHTML)
				}
			} else {
				assert.Equal(t, tt.expectedHTML, outputHTML)
			}
		})
	}
}
