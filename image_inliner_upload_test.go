package gophetch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/minsoft-io/gophetch"
	"github.com/minsoft-io/gophetch/image"
)

func TestUploadStrategy(t *testing.T) {
	SetupFiles()

	tests := []struct {
		name         string
		inputHTML    string
		expectedHTML string
	}{
		{
			name:         "image gets uploaded",
			inputHTML:    `<img src="mark.png">`,
			expectedHTML: `<html><head></head><body><img src="new_image_url.png"/></body></html>`,
		},
		{
			name:         "image gets uploaded with srcset",
			inputHTML:    `<img srcset="https://example.com/mark.png x1, https://example.com/mark.png x2">`,
			expectedHTML: `<html><head></head><body><img srcset="new_image_url.png x1, new_image_url.png x2"/></body></html>`,
		},
		{
			name:         "image gets uploaded with srcset and sizes",
			inputHTML:    `<img srcset="https://example.com/mark.png x1, https://example.com/mark.png x2" sizes="100vw">`,
			expectedHTML: `<html><head></head><body><img srcset="new_image_url.png x1, new_image_url.png x2" sizes="100vw"/></body></html>`,
		},
	}

	mockFetcher := new(MockImageFetcher)
	mockUploadFunc := func(img *image.Image) (string, error) {
		return "new_image_url.png", nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inliner := gophetch.NewImageInliner(gophetch.ImageInlinerOptions{
				Fetcher:        mockFetcher,
				UploadFunc:     mockUploadFunc,
				InlineStrategy: gophetch.InlineNone,
			})

			actualHTML, err := inliner.InlineImages(tt.inputHTML)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedHTML, actualHTML)
		})
	}
}
