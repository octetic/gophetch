package rules

import (
	"bytes"
	"errors"
	"image/gif"
	_ "image/gif" // This is required to initialize the GIF decoder
	"image/jpeg"
	_ "image/jpeg" // This is required to initialize the JPEG decoder
	"image/png"
	_ "image/png" // This is required to initialize the PNG decoder
	"io"
	"net/http"
	"strings"

	ico "github.com/biessek/golang-ico"
)

var ErrInvalidImageFormat = errors.New("invalid image format")

type LeadImageRule struct {
	BaseRule
}

func NewLeadImageRule() *LeadImageRule {
	return &LeadImageRule{
		BaseRule: BaseRule{
			Strategies: leadImageStrategies,
		},
	}
}

var leadImageStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			"meta[property='og:image:secure_url']",
			"meta[property='og:image:url']",
			"meta[property='og:image']",
			"meta[name='og:image']",
			"meta[name='twitter:image:src']",
			"meta[property='twitter:image:src']",
			"meta[name='twitter:image']",
			"meta[property='twitter:image']",
			"meta[itemprop='image']",
		},
		Extractor: ExtractAttr("content"),
	},
	{
		Selectors: []string{
			"img[src]:not([width='1']):not([height='1'])",
			"img[srcset]:not([width='1']):not([height='1'])",
			"img[data-src]:not([width='1']):not([height='1'])",
			"img[data-srcset]:not([width='1']):not([height='1'])",
			"img[data-lazy-src]:not([width='1']):not([height='1'])",
			"img[data-lazy-srcset]:not([width='1']):not([height='1'])",
			"img[data-lazyload]:not([width='1']):not([height='1'])",
		},
		Extractor: ExtractCSS,
	},
}

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
