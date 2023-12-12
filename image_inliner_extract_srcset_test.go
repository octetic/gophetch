package gophetch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/octetic/gophetch"
)

func TestExtractSrcset(t *testing.T) {
	tests := []struct {
		name          string
		srcset        string
		expectedURLs  []string
		expectedDescs []string
	}{
		{
			name:          "single source",
			srcset:        `https://example.com/image1.png 1x`,
			expectedURLs:  []string{"https://example.com/image1.png"},
			expectedDescs: []string{"1x"},
		},
		{
			name:   "multiple sources",
			srcset: `https://example.com/image1.png 1x, https://example.com/image2.png 2x`,
			expectedURLs: []string{
				"https://example.com/image1.png",
				"https://example.com/image2.png",
			},
			expectedDescs: []string{"1x", "2x"},
		},
		{
			name:   "single source with width",
			srcset: `https://example.com/image1.png 484w`,
			expectedURLs: []string{
				"https://example.com/image1.png",
			},
			expectedDescs: []string{"484w"},
		},
		{
			name:   "multiple sources with width",
			srcset: `https://example.com/image1.png 484w, https://example.com/image2.png 968w`,
			expectedURLs: []string{
				"https://example.com/image1.png",
				"https://example.com/image2.png",
			},
			expectedDescs: []string{"484w", "968w"},
		},
		{
			name:   "multiple sources with width and density",
			srcset: `https://example.com/image1.png 484w 1x, https://example.com/image2.png 968w 2x, https://example.com/image3.png 1936w 4x`,
			expectedURLs: []string{
				"https://example.com/image1.png",
				"https://example.com/image2.png",
				"https://example.com/image3.png",
			},
			expectedDescs: []string{"484w 1x", "968w 2x", "1936w 4x"},
		},
		{
			name:   "multiple sources with query params",
			srcset: `https://example.com/image1.png?foo=bar 484w 1x, https://example.com/image2.png?foo=bar 968w 2x, https://example.com/image3.png?foo=bar 1936w 4x`,
			expectedURLs: []string{
				"https://example.com/image1.png?foo=bar",
				"https://example.com/image2.png?foo=bar",
				"https://example.com/image3.png?foo=bar",
			},
			expectedDescs: []string{"484w 1x", "968w 2x", "1936w 4x"},
		},
		{
			name:   "multiples sources with decimal density",
			srcset: `https://example.com/image1.png 1.5x, https://example.com/image2.png 2.5x, https://example.com/image3.png 3.5x`,
			expectedURLs: []string{
				"https://example.com/image1.png",
				"https://example.com/image2.png",
				"https://example.com/image3.png",
			},
			expectedDescs: []string{"1.5x", "2.5x", "3.5x"},
		},
		{
			name:   "multiples sources with decimal density and width",
			srcset: `https://example.com/image1.png 484w 1.5x, https://example.com/image2.png 968w 2.5x, https://example.com/image3.png 1936w 3.5x`,
			expectedURLs: []string{
				"https://example.com/image1.png",
				"https://example.com/image2.png",
				"https://example.com/image3.png",
			},
			expectedDescs: []string{"484w 1.5x", "968w 2.5x", "1936w 3.5x"},
		},
		{
			name:          "multiple sources with no descriptors (non-standard)",
			srcset:        `https://example.com/image1.png, https://example.com/image2.png, https://example.com/image3.png`,
			expectedURLs:  nil,
			expectedDescs: nil,
		},
		{
			name:   "substack style srcset",
			srcset: `https://substackcdn.com/image/fetch/w_424,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fexample.png 424w, https://substackcdn.com/image/fetch/w_848,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fexample.png 848w`,
			expectedURLs: []string{
				"https://substackcdn.com/image/fetch/w_424,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fexample.png",
				"https://substackcdn.com/image/fetch/w_848,c_limit,f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fexample.png",
			},
			expectedDescs: []string{"424w", "848w"},
		},
		{
			name:          "empty srcset",
			srcset:        ``,
			expectedURLs:  nil,
			expectedDescs: nil,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls, descs := gophetch.ExtractSrcset(tt.srcset)
			assert.Equal(t, tt.expectedURLs, urls)
			assert.Equal(t, tt.expectedDescs, descs)
		})
	}
}
