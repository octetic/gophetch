package rules

import (
	"fmt"
	"net/url"

	"golang.org/x/net/html"
)

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
	if err == nil && len(result.Value) > 0 {
		return result, nil
	}

	// If no favicon was found, try to extract it from the /favicon.ico file.
	faviconUrl := fmt.Sprintf("%s://%s/favicon.ico", targetURL.Scheme, targetURL.Host)
	if IsValidImage(faviconUrl) {
		return ExtractResult{
			Value: []string{faviconUrl},
			Selector: SelectorInfo{
				Attr:     "href",
				InMeta:   false,
				Selector: "favicon.ico",
			},
			Found: true,
		}, nil
	}

	return ExtractResult{}, ErrValueNotFound
}
