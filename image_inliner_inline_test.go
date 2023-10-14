package gophetch_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/minsoft-io/gophetch"
)

func TestInlinerInlineStrategy(t *testing.T) {
	SetupFiles()

	tests := []struct {
		name          string
		inputHTML     string
		expectedFiles []string
		expectError   bool
	}{
		{
			name:          "empty HTML",
			inputHTML:     "",
			expectedFiles: []string{},
			expectError:   false,
		},
		// Add more cases here
		{
			name:          "single src",
			inputHTML:     `<html><head><body><img src="https://example.com/mark.jpg"></body></html>`,
			expectedFiles: []string{"mark.jpg"},
			expectError:   false,
		},
		{
			name:      "srcset multiple sources",
			inputHTML: `<html><head><body><img srcset="https://example.com/mark.jpg 1x, https://example.com/mark.png 2x"></body></html>`,
			expectedFiles: []string{
				"mark.jpg",
				"mark.png",
			},
			expectError: false,
		},
		{
			name: "multiple images",
			inputHTML: `<html><head><body>
			<img src="https://example.com/mark.jpg">
			<img src="https://example.com/mark.png">
			</body></html>`,
			expectedFiles: []string{
				"mark.jpg",
				"mark.png",
			},
		},
		{
			name: "multiple images with srcset",
			inputHTML: `<html><head><body>
			<img src="https://example.com/mark.jpg">
			<img srcset="https://example.com/mark.png 1x,https://example.com/mark.webp 2x">
			</body></html>`,
			expectedFiles: []string{
				"mark.jpg",
				"mark.png",
			},
		},
		{
			name: "multiple images with srcset and sizes",
			inputHTML: `<html><head><body>
			<img src="https://example.com/mark.jpg">
			<img srcset="https://example.com/mark.png 1x,https://example.com/mark.webp 2x" sizes="100vw">
			</body></html>`,
			expectedFiles: []string{
				"mark.jpg",
				"mark.png",
			},
		},
		{
			name: "multiple images with srcset and sizes and width",
			inputHTML: `<html><head><body>
			<img src="https://example.com/mark.jpg">
			<img srcset="https://example.com/mark.png 1x,https://example.com/mark.webp 2x" sizes="100vw" width="100">
			</body></html>`,
			expectedFiles: []string{
				"mark.jpg",
				"mark.png",
			},
		},
		{
			name: "multiple images with srcset and sizes and width and height",
			inputHTML: `<html><head><body>
			<img src="https://example.com/mark.jpg">
			<img srcset="https://example.com/mark.png 1x,https://example.com/mark.webp 2x" sizes="100vw" width="100" height="100">
			</body></html>`,
			expectedFiles: []string{
				"mark.jpg",
				"mark.png",
			},
		},
		{
			name: "multiple images nested in divs",
			inputHTML: `<html><head><body>
		    <main>
			<div><img src="https://example.com/mark.jpg"></div>
			<div><img src="https://example.com/mark.png"></div>	
			<div><img srcset="https://example.com/mark.png 1x,https://example.com/mark.webp 2x"></div>
			</main>
			</body></html>`,
			expectedFiles: []string{
				"mark.jpg",
				"mark.png",
				"mark.webp",
			},
		},
		{
			name: "video with poster",
			inputHTML: `<html><head><body>
 		    <main>
			<video poster="https://example.com/mark.jpg">
				<source src="https://example.com/mark.mp4" type="video/mp4">
			</video>
			</main>
			</body></html>`,
			expectedFiles: []string{
				"mark.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFetcher := &MockImageFetcher{}
			inliner := gophetch.NewImageInliner(gophetch.ImageInlinerOptions{
				Fetcher:  mockFetcher,
				Strategy: gophetch.StrategyInline,
			})
			outputHTML, err := inliner.InlineImages(tt.inputHTML)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Here you might want to use a more sophisticated HTML comparison,
			// but for simplicity, we're using string comparison.
			//assert.Equal(t, tt.expectedHTML, strings.TrimSpace(outputHTML))

			for _, filename := range tt.expectedFiles {
				expectedBase64, ok := imageBase64Map[filename]
				assert.True(t, ok, "Expected base64 for %s in testdata not found", filename)
				assert.Contains(t, outputHTML, expectedBase64, "Base64 data string for %s not found. Output was: \n %s", filename, outputHTML)

				if strings.Contains(tt.inputHTML, "sizes=\"100vw\"") {
					assert.Contains(t, outputHTML, "sizes=\"100vw\"", "sizes attribute not found in output HTML")
				}

				if strings.Contains(tt.inputHTML, "width=\"100\"") {
					assert.Contains(t, outputHTML, "width=\"100\"", "width attribute not found in output HTML")
				}

				if strings.Contains(tt.inputHTML, "height=\"100\"") {
					assert.Contains(t, outputHTML, "height=\"100\"", "height attribute not found in output HTML")
				}
			}
		})
	}
}
