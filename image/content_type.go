package image

import (
	"fmt"
	"io"
	"net/http"
)

var contentTypeToExt = map[string][]string{
	"image/avif":               {".avif"},
	"image/bmp":                {".bmp"},
	"image/gif":                {".gif"},
	"image/heic":               {".heic"},
	"image/heif":               {".heif"},
	"image/jpeg":               {".jpg", ".jpeg"},
	"image/png":                {".png"},
	"image/svg+xml":            {".svg"},
	"image/tiff":               {".tif", ".tiff"},
	"image/vnd.microsoft.icon": {".ico"},
	"image/webp":               {".webp"},
	"image/x-icon":             {".ico"},
	"image/x-windows-bmp":      {".bmp"},
}

var extToContentType = map[string]string{}

var invalidFaviconContentTypes = map[string]struct{}{
	"image/avif": {},
	"image/heic": {},
	"image/heif": {},
	"image/tiff": {},
}

func init() {
	// Build the reverse mapping from extension to content-type
	for contentType, extension := range contentTypeToExt {
		for _, ext := range extension {
			extToContentType[ext] = contentType
		}
	}
}

func IsValidImageContentType(contentType string) bool {
	_, ok := contentTypeToExt[contentType]
	return ok
}

func IsValidFaviconContentType(contentType string) bool {
	if _, ok := invalidFaviconContentTypes[contentType]; ok {
		return false
	}
	return IsValidImageContentType(contentType)
}

func ContentTypeByExtension(extension string) (string, error) {
	switch extension {
	case ".ico":
		return "image/x-icon", nil
	case ".bmp":
		return "image/bmp", nil
	default:
		contentType, ok := extToContentType[extension]
		if !ok {
			return "", fmt.Errorf("unknown extension: %s", extension)
		}
		return contentType, nil
	}
}

func ExtensionByContentType(contentType string) (string, error) {
	extension, ok := contentTypeToExt[contentType]
	if !ok {
		return "", fmt.Errorf("unknown content-type: %s", contentType)
	}
	return extension[0], nil
}

// ContentTypeFromURL returns the Content-Type header from the given URL. However, it does not download the image.
// It will attempt to get the Content-Type header by making a HEAD request to the URL. If the Content-Type header
// cannot be extracted using a HEAD request, it will attempt to get it by making a GET request
// to the URL and reviewing the first 512 KB. An error is returned if the Content-Type header still cannot be extracted.
func ContentTypeFromURL(url string) (string, error) {
	resp, err := http.Head(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Check for a successful or redirected response
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		meta, err := FetchMetadataFromHeader(url, defaultMaxBytes)
		if err != nil {
			return "", err
		}
		return meta.ContentType, nil
	}

	contentType := resp.Header.Get("Content-Type")
	return contentType, nil
}
