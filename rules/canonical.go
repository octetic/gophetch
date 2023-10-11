package rules

import (
	"errors"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type CanonicalRule struct {
	BaseRule
}

func NewCanonicalRule() *CanonicalRule {
	return &CanonicalRule{
		BaseRule: BaseRule{
			Strategies: canonicalStrategies,
		},
	}
}

var canonicalStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			"meta[property='og:url']",
			"meta[name='twitter:url']",
			"meta[property='twitter:url']",
		},
		Extractor: ExtractAttr("content"),
	},
	{
		Selectors: []string{
			"link[rel='canonical']",
			"link[rel='alternate'][hreflang='x-default']",
		},
		Extractor: ExtractAttr("href"),
	},
}

func (cr *CanonicalRule) Extract(node *html.Node, targetURL *url.URL) ([]string, error) {
	values, err := cr.BaseRule.Extract(node, targetURL)
	if err != nil {
		if !errors.Is(err, ErrValueNotFound) {
			return []string{}, err
		}
	}

	// If the value is not a full URL, then prepend the scheme and host
	for i, value := range values {
		if !strings.HasPrefix(value, "http") {
			values[i] = targetURL.Scheme + "://" + targetURL.Host + value
		}
	}

	if len(values) > 0 {
		return values, nil
	}

	// If the values are empty, then return the current url
	return []string{targetURL.String()}, nil
}
