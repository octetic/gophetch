package helpers

import (
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

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

// FixRelativePath converts a relative path to an absolute path for the given URL.
func FixRelativePath(url *url.URL, path string) string {
	path = strings.TrimSpace(path)

	if strings.HasPrefix(path, "http") {
		return path
	}

	var buf strings.Builder

	if strings.HasPrefix(path, "//") {
		if url.Scheme != "" {
			buf.WriteString(url.Scheme + ":")
		}
		buf.WriteString(path)
		return buf.String()
	}

	if url.Scheme != "" {
		buf.WriteString(url.Scheme + ":")
	}

	if url.Host != "" || url.Path != "" || url.User != nil {
		buf.WriteString("//")

		if ui := url.User; ui != nil {
			buf.WriteString(ui.String() + "@")
		}

		buf.WriteString(url.Host)
	}

	if url.Host != "" && path != "" && path[0] != '/' {
		buf.WriteByte('/')
	}

	buf.WriteString(path)

	return buf.String()
}
