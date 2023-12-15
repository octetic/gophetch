package media

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
	// pdf content type
	"application/pdf": {".pdf"},
	// doc content type
	"application/msword": {".doc"},
	// docx content type
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {".docx"},
	// odt content type
	"application/vnd.oasis.opendocument.text": {".odt"},
	// xls content type
	"application/vnd.ms-excel": {".xls"},
	// xlsx content type
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": {".xlsx"},
	// ods content type
	"application/vnd.oasis.opendocument.spreadsheet": {".ods"},
	// ppt content type
	"application/vnd.ms-powerpoint": {".ppt"},
	// pptx content type
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": {".pptx"},
	// odp content type
	"application/vnd.oasis.opendocument.presentation": {".odp"},
	// zip content type
	"application/zip": {".zip"},
	// rar content type
	"application/x-rar-compressed": {".rar"},
	// tar content type
	"application/x-tar": {".tar"},
	// 7z content type
	"application/x-7z-compressed": {".7z"},
	// mp3 content type
	"audio/mpeg": {".mp3"},
	// mp4 content type
	"video/mp4": {".mp4"},
	// webm content type
	"video/webm": {".webm"},
	// ogg content type
	"video/ogg": {".ogg"},
	// ogg content type
	"audio/ogg": {".ogg"},
	// wav content type
	"audio/wav": {".wav"},
	// woff content type
	"font/woff": {".woff"},
	// woff2 content type
	"font/woff2": {".woff2"},
	// ttf content type
	"font/ttf": {".ttf"},
	// eot content type
	"application/vnd.ms-fontobject": {".eot"},
	// otf content type
	"font/otf": {".otf"},
	// txt content type
	"text/plain": {".txt"},
	// csv content type
	"text/csv": {".csv"},
	// html content type
	"text/html": {".html"},
	// xml content type
	"text/xml": {".xml"},
	// json content type
	"application/json": {".json"},
	// js content type
	"application/javascript": {".js"},
	// css content type
	"text/css": {".css"},
}

var extToContentType = map[string]string{}

var invalidFaviconContentTypes = map[string]struct{}{
	"image/avif": {},
	"image/heic": {},
	"image/heif": {},
	"image/tiff": {},
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         {},
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": {},
	"font/woff":                     {},
	"font/woff2":                    {},
	"font/ttf":                      {},
	"application/vnd.ms-fontobject": {},
	"font/otf":                      {},
	"text/plain":                    {},
	"text/csv":                      {},
	"text/xml":                      {},
	"application/json":              {},
	"application/javascript":        {},
	"text/css":                      {},
	"application/pdf":               {},
	"application/msword":            {},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {},
	"application/vnd.ms-excel":      {},
	"application/vnd.ms-powerpoint": {},
	"application/zip":               {},
	"application/x-rar-compressed":  {},
	"application/x-tar":             {},
	"application/x-7z-compressed":   {},
	"audio/mpeg":                    {},
	"video/mp4":                     {},
	"video/webm":                    {},
	"video/ogg":                     {},
	"audio/ogg":                     {},
	"audio/wav":                     {},
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
