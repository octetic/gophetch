package rules

import (
	"fmt"
	"net/url"

	"golang.org/x/net/html"
)

// FaviconRule is the rule for extracting the favicon URL of a page.
type FaviconRule struct {
	BaseRule
}

func NewFaviconRule() *FaviconRule {
	return &FaviconRule{
		BaseRule: BaseRule{
			Strategies: faviconStrategies,
		},
	}
}

var faviconStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			"link[rel='icon']",
			"link[rel='shortcut icon']",
			"link[rel='apple-touch-icon']",
			"link[rel='apple-touch-icon-precomposed']",
			"link[rel~='mask-icon']",
		},
		Extractor: ExtractAttr("href"),
	},
}

func (r *FaviconRule) Extract(node *html.Node, targetURL *url.URL) (ExtractResult, error) {
	result, err := r.BaseRule.Extract(node, targetURL)
	if err == nil && result.Found() {
		return result, nil
	}

	// If no favicon was found, try to extract it from the /favicon.ico file.
	faviconURL := fmt.Sprintf("%s://%s/favicon.ico", targetURL.Scheme, targetURL.Host)
	if IsValidImage(faviconURL) {
		return NewStringResult(
			faviconURL,
			SelectorInfo{
				Attr:     "href",
				InMeta:   false,
				Selector: "favicon.ico",
			},
			true,
		), nil
	}

	return NewNoResult(), ErrValueNotFound
}
