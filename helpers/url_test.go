package helpers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/octetic/gophetch/helpers"
)

func TestCleanURL(t *testing.T) {
	// Table-driven tests
	testCases := []struct {
		name     string
		inputURL string
		want     string
	}{
		{
			name:     "URL with multiple tracking params",
			inputURL: "https://example.com?fbclid=123&gclid=456&normal_param=value",
			want:     "https://example.com?normal_param=value",
		},
		{
			name:     "URL with tracking params and path",
			inputURL: "https://example.com/foobar/?fbclid=123&gclid=456&normal_param=value",
			want:     "https://example.com/foobar/?normal_param=value",
		},
		{
			name:     "URL with no tracking params",
			inputURL: "https://example.com?normal_param=value",
			want:     "https://example.com?normal_param=value",
		},
		{
			name:     "URL with no params and path",
			inputURL: "https://example.com/foobar/",
			want:     "https://example.com/foobar/",
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := helpers.CleanURL(tc.inputURL)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIsURLValid(t *testing.T) {
	// Table-driven tests
	testCases := []struct {
		name     string
		inputURL string
		want     bool
	}{
		{
			name:     "Valid http URL",
			inputURL: "http://example.com",
			want:     true,
		},
		{
			name:     "Valid https URL",
			inputURL: "https://example.com/foobar/?fbclid=123&gclid=456&normal_param=value",
			want:     true,
		},
		{
			name:     "Invalid URL with no scheme",
			inputURL: "example.com",
			want:     false,
		},
		{
			name:     "Invalid URL with only a scheme",
			inputURL: "http://",
			want:     false,
		},
		{
			name:     "Invalid URL with scheme and host",
			inputURL: "http://examplecom",
			want:     false,
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := helpers.IsURLValid(tc.inputURL)
			assert.Equal(t, tc.want, got)
		})
	}
}
