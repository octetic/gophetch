package image_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/minsoft-io/gophetch/image"
)

func TestParseCacheHeader(t *testing.T) {
	testCases := []struct {
		name     string
		header   http.Header
		expected image.Cache
	}{
		{
			name: "Cache-Control no-cache",
			header: http.Header{
				"Cache-Control": []string{"no-cache"},
			},
			expected: image.Cache{
				Available:      true,
				MaxAge:         0,
				NoCache:        true,
				NoStore:        false,
				MustRevalidate: false,
			},
		},
		{
			name: "Cache-Control no-store",
			header: http.Header{
				"Cache-Control": []string{"no-store"},
			},
			expected: image.Cache{
				Available:      true,
				MaxAge:         0,
				NoCache:        false,
				NoStore:        true,
				MustRevalidate: false,
			},
		},
		{
			name: "Cache-Control must-revalidate",
			header: http.Header{
				"Cache-Control": []string{"must-revalidate"},
			},
			expected: image.Cache{
				Available:      true,
				MaxAge:         0,
				NoCache:        false,
				NoStore:        false,
				MustRevalidate: true,
			},
		},
		{
			name: "Cache-Control no-cache, no-store",
			header: http.Header{
				"Cache-Control": []string{"no-cache,no-store"},
			},
			expected: image.Cache{
				Available:      true,
				MaxAge:         0,
				NoCache:        true,
				NoStore:        true,
				MustRevalidate: false,
			},
		},
		{
			name: "Cache-Control no-cache, must-revalidate",
			header: http.Header{
				"Cache-Control": []string{"no-cache, must-revalidate"},
			},
			expected: image.Cache{
				Available:      true,
				MaxAge:         0,
				NoCache:        true,
				NoStore:        false,
				MustRevalidate: true,
			},
		},
		{
			name:   "No Cache-Control or Expires",
			header: http.Header{},
			expected: image.Cache{
				Available:      false,
				MaxAge:         -1,
				Expires:        time.Time{},
				NoCache:        false,
				NoStore:        false,
				MustRevalidate: false,
			},
		},
		{
			name: "No Cache-Control, Expires in the past",
			header: http.Header{
				"Expires": []string{"Thu, 01 Jan 1970 00:00:00 GMT"},
			},
			expected: image.Cache{
				Available:      true,
				MaxAge:         -1,
				Expires:        parseTime("Thu, 01 Jan 1970 00:00:00 GMT"),
				NoCache:        false,
				NoStore:        false,
				MustRevalidate: false,
			},
		},
		{
			name: "No Cache-Control, Expires in the future",
			header: http.Header{
				"Expires": []string{"Thu, 01 Jan 2030 00:00:00 GMT"},
			},
			expected: image.Cache{
				Available:      true,
				MaxAge:         -1,
				Expires:        parseTime("Thu, 01 Jan 2030 00:00:00 GMT"),
				NoCache:        false,
				NoStore:        false,
				MustRevalidate: false,
			},
		},
		{
			name: "Max-Age directive",
			header: http.Header{
				"Cache-Control": []string{"max-age=60"},
			},
			expected: image.Cache{
				Available:      true,
				MaxAge:         60,
				NoCache:        false,
				NoStore:        false,
				MustRevalidate: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := image.ParseCacheHeader(tc.header)
			assert.Equal(t, tc.expected.Available, actual.Available, "Available mismatch")
			assert.Equal(t, tc.expected.MaxAge, actual.MaxAge, "MaxAge mismatch")
			assert.Equal(t, tc.expected.NoCache, actual.NoCache, "NoCache mismatch")
			assert.Equal(t, tc.expected.NoStore, actual.NoStore, "NoStore mismatch")
			assert.Equal(t, tc.expected.MustRevalidate, actual.MustRevalidate, "MustRevalidate mismatch")

			// TODO: add more specific requirements on how to compare `Expires`
			if !tc.expected.Expires.IsZero() {
				assert.True(t, tc.expected.Expires.Sub(actual.Expires) < time.Second)
			}
		})
	}
}

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC1123, s)
	if err != nil {
		panic(err)
	}
	return t
}
