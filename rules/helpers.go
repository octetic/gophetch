package rules

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

var validImageMIMETypes = map[string]struct{}{
	"image/jpeg":               {},
	"image/png":                {},
	"image/gif":                {},
	"image/bmp":                {},
	"image/x-windows-bmp":      {},
	"image/webp":               {},
	"image/svg+xml":            {},
	"image/tiff":               {},
	"image/vnd.microsoft.icon": {},
	"image/x-icon":             {},
	"image/avif":               {},
	"image/heif":               {},
	"image/heic":               {},
}

var validFaviconMIMETypes = map[string]struct{}{
	"image/x-icon":             {},
	"image/vnd.microsoft.icon": {},
	"image/jpeg":               {},
	"image/png":                {},
	"image/gif":                {},
	"image/bmp":                {},
	"image/x-windows-bmp":      {},
	"image/webp":               {},
	"image/svg+xml":            {},
}

func isValidImageMIMEType(mimeType string) bool {
	_, ok := validImageMIMETypes[mimeType]
	return ok
}

func isValidFaviconMIMEType(mimeType string) bool {
	_, ok := validFaviconMIMETypes[mimeType]
	return ok
}

func IsValidImage(url string) bool {
	valid, err := ValidateImage(url, false)
	if err != nil {
		return false
	}
	return valid
}

// IsValidFavicon checks if the given URL is a valid favicon.
func IsValidFavicon(url string) bool {
	valid, err := ValidateImage(url, true)
	if err != nil {
		return false
	}
	return valid
}

// ValidateImage checks if the given URL is a valid image.
func ValidateImage(url string, faviconOnly bool) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Check for a successful or redirected response
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return false, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")

	if faviconOnly {
		if !isValidFaviconMIMEType(contentType) {
			return false, fmt.Errorf("invalid MIME type: %s", contentType)
		}
	} else {
		if !isValidImageMIMEType(contentType) {
			return false, fmt.Errorf("invalid MIME type: %s", contentType)
		}
	}

	return true, nil
}

// FixRelativePath converts a relative path to an absolute path for the given URL.
func FixRelativePath(u *url.URL, path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	} else if strings.HasPrefix(path, "//") {
		return u.Scheme + ":" + path
	} else {
		if strings.HasPrefix(path, "/") {
			return u.Scheme + "://" + u.Host + path
		}
		return u.Scheme + "://" + u.Host + "/" + path
	}
}

// Normalize cleans up the extracted string, removing HTML tags,
// decoding HTML entities, and trimming whitespace.
func Normalize(input string) string {
	// Strip HTML tags
	p := bluemonday.StripTagsPolicy()
	clean := p.Sanitize(input)

	// Decode HTML entities
	decoded := html.UnescapeString(clean)

	// Trim whitespace
	normalized := strings.TrimSpace(decoded)

	return normalized
}
