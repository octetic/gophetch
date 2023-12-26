package gophetch_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/octetic/gophetch"
)

func TestPrefixStrategy(t *testing.T) {
	SetupFiles()

	relativeURL, err := url.Parse("https://example.com/foo/bar/baz")
	assert.NoError(t, err)

	prefixProxy := "https://example.art/proxy"

	tests := []struct {
		name         string
		inputHTML    string
		expectedHTML string
		expectError  bool
	}{
		{
			name:         "empty HTML",
			inputHTML:    "",
			expectedHTML: "<html><head></head><body></body></html>",
			expectError:  false,
		},
		// Add more cases here
		{
			name:         "single src",
			inputHTML:    `<html><head></head><body><img src="https://example.com/mark.jpg"></body></html>`,
			expectedHTML: `<html><head></head><body><img src="` + prefixProxy + `?url=https://example.com/mark.jpg"/></body></html>`,
			expectError:  false,
		},
		{
			name:         "srcset multiple sources",
			inputHTML:    `<html><head></head><body><img srcset="https://example.com/mark.jpg 1x, https://example.com/mark.png 2x"></body></html>`,
			expectedHTML: `<html><head></head><body><img srcset="` + prefixProxy + `?url=https://example.com/mark.jpg 1x, ` + prefixProxy + `?url=https://example.com/mark.png 2x"/></body></html>`,
			expectError:  false,
		},
		{
			name: "multiple images",
			inputHTML: `<html><head></head><body>
			<img src="https://example.com/mark.jpg">
			<img src="https://example.com/mark.png">
			</body></html>`,
			expectedHTML: `<html><head></head><body>
			<img src="` + prefixProxy + `?url=https://example.com/mark.jpg"/>
			<img src="` + prefixProxy + `?url=https://example.com/mark.png"/>
			</body></html>`,
		},
		{
			name: "multiple images with srcset",
			inputHTML: `<html><head></head><body>
			<img src="https://example.com/mark.jpg">
			<img srcset="https://example.com/mark.png 1x,https://example.com/mark.webp 2x">
			</body></html>`,
			expectedHTML: `<html><head></head><body>
			<img src="` + prefixProxy + `?url=https://example.com/mark.jpg"/>
			<img srcset="` + prefixProxy + `?url=https://example.com/mark.png 1x, ` + prefixProxy + `?url=https://example.com/mark.webp 2x"/>
			</body></html>`,
		},
		{
			name: "multiple images with srcset and sizes",
			inputHTML: `<html><head></head><body>
			<img src="https://example.com/mark.jpg">
			<img srcset="https://example.com/mark.png 1x,https://example.com/mark.webp 2x" sizes="(max-width: 600px) 200px, 50vw">
			</body></html>`,
			expectedHTML: `<html><head></head><body>
			<img src="` + prefixProxy + `?url=https://example.com/mark.jpg"/>
			<img srcset="` + prefixProxy + `?url=https://example.com/mark.png 1x, ` + prefixProxy + `?url=https://example.com/mark.webp 2x" sizes="(max-width: 600px) 200px, 50vw"/>
			</body></html>`,
		},
		{
			name: "multiple images with picture > source",
			inputHTML: `<html><head></head><body>
			<picture>
				<source srcset="https://example.com/mark.png 1x,https://example.com/mark.webp 2x" sizes="(max-width: 600px) 200px, 50vw">	
				<img src="https://example.com/mark.jpg">
			</picture>
			</body></html>`,
			expectedHTML: `<html><head></head><body>
			<picture>
				<source srcset="` + prefixProxy + `?url=https://example.com/mark.png 1x, ` + prefixProxy + `?url=https://example.com/mark.webp 2x" sizes="(max-width: 600px) 200px, 50vw"/>	
				<img src="` + prefixProxy + `?url=https://example.com/mark.jpg"/>
			</picture>
			</body></html>`,
		},
		{
			name: "multiple images with relative URLs",
			inputHTML: `<html><head></head><body>
			<img src="/mark.jpg">
			<img srcset="/mark.png 1x,/mark.webp 2x">
			</body></html>`,
			expectedHTML: `<html><head></head><body>
			<img src="` + prefixProxy + `?url=https://example.com/mark.jpg"/>
			<img srcset="` + prefixProxy + `?url=https://example.com/mark.png 1x, ` + prefixProxy + `?url=https://example.com/mark.webp 2x"/>
			</body></html>`,
		},
		{
			name: "data:image/png;base64 image",
			inputHTML: `<html><head></head><body>
			<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wr4H/wAAAABJRU5ErkJggg==">
			</body></html>`,
			expectedHTML: `<html><head></head><body>
			<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wr4H/wAAAABJRU5ErkJggg=="/>
			</body></html>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFetcher := new(MockImageFetcher)
			inliner := gophetch.NewImageInliner(gophetch.ImageInlinerOptions{
				Fetcher:        mockFetcher,
				InlineStrategy: gophetch.InlineMediaProxy,
				SrcsetStrategy: gophetch.SrcsetAllImages,
				MediaProxyURL:  prefixProxy,
				RelativeURL:    relativeURL,
			})

			actualHTML, err := inliner.InlineImages(tt.inputHTML)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedHTML, actualHTML)
		})
	}
}
