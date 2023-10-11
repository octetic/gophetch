package rules

import (
	"bytes"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"strings"

	ico "github.com/biessek/golang-ico"
)

func IsValidImage(url string) bool {
	valid, err := ValidateImage(url)
	if err != nil {
		return false
	}
	return valid
}

func ValidateImage(url string) (bool, error) {
	if url == "" {
		return false, nil
	}

	// Make a GET request to the image URL.
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Read in the image as bytes
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Check the content type is one of the valid favicon types.
	contentType := http.DetectContentType(data)
	switch {
	case strings.HasPrefix(contentType, "image/png"):
		_, err = png.Decode(bytes.NewReader(data))
	case strings.HasPrefix(contentType, "image/jpeg"):
		_, err = jpeg.Decode(bytes.NewReader(data))
	case strings.HasPrefix(contentType, "image/gif"):
		_, err = gif.Decode(bytes.NewReader(data))
	//case strings.HasPrefix(contentType, "image/webp"):
	//	_, err = webp.Decode(bytes.NewReader(data))
	case strings.HasPrefix(contentType, "image/x-icon"), strings.HasPrefix(contentType, "image/vnd.microsoft.icon"):
		_, err = ico.Decode(bytes.NewReader(data))
	default:
		return false, ErrInvalidImageFormat
	}

	// We'll get here if the content type was one of the valid types. If there was an error decoding the image, it was invalid.
	if err != nil {
		return false, err
	}

	// If we got here, the image is valid.
	return true, nil
}

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
