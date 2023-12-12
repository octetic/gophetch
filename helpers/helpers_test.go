package helpers_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/octetic/gophetch/helpers"
)

func TestFixRelativePath(t *testing.T) {
	testCases := []struct {
		name     string
		url      *url.URL
		path     string
		expected string
	}{
		{
			name:     "absolute URL",
			url:      &url.URL{Scheme: "http", Host: "example.com"},
			path:     "http://example.com",
			expected: "http://example.com",
		},
		{
			name:     "relative URL with leading //",
			url:      &url.URL{Scheme: "http", Host: "example.com"},
			path:     "//example.com",
			expected: "http://example.com",
		},
		{
			name:     "relative URL with leading /",
			url:      &url.URL{Scheme: "http", Host: "example.com"},
			path:     "/page",
			expected: "http://example.com/page",
		},
		{
			name:     "relative URL without leading /",
			url:      &url.URL{Scheme: "http", Host: "example.com"},
			path:     "page",
			expected: "http://example.com/page",
		},
		{
			name:     "URL with user info",
			url:      &url.URL{Scheme: "http", Host: "example.com", User: url.UserPassword("user", "pass")},
			path:     "/page",
			expected: "http://user:pass@example.com/page",
		},
		{
			name:     "path with spaces",
			url:      &url.URL{Scheme: "http", Host: "example.com"},
			path:     " /page ",
			expected: "http://example.com/page",
		},
		{
			name:     "empty path",
			url:      &url.URL{Scheme: "http", Host: "example.com"},
			path:     "",
			expected: "http://example.com",
		},
		{
			name:     "data URL",
			url:      &url.URL{Scheme: "http", Host: "example.com"},
			path:     "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wr+3HwAAAABJRU5ErkJggg==",
			expected: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wr+3HwAAAABJRU5ErkJggg==",
		},
		{
			name:     "empty URL",
			url:      &url.URL{},
			path:     "/page",
			expected: "/page",
		}, // Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := helpers.FixRelativePath(tc.url, tc.path)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
