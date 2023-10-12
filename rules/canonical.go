package rules

import (
	"errors"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// CanonicalRule is the rule for extracting the canonical URL of a page.
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

func (cr *CanonicalRule) Extract(node *html.Node, targetURL *url.URL) (ExtractResult, error) {
	result, err := cr.BaseRule.Extract(node, targetURL)
	if err != nil {
		if !errors.Is(err, ErrValueNotFound) {
			return &NoResult{}, err
		}
	}

	if result.Found() {

		sr, ok := result.(*StringResult)
		if !ok {
			return NewNoResult(), errors.New("invalid result type")
		}

		inMeta := sr.SelectorInfo().Selector == "content"

		// If the value is not a full URL, then prepend the scheme and host
		if !strings.HasPrefix(sr.value, "http") {
			sr.value = targetURL.Scheme + "://" + targetURL.Host + sr.value
		}

		return NewStringResult(
			sr.value,
			SelectorInfo{
				Attr:     sr.SelectorInfo().Attr,
				InMeta:   inMeta,
				Selector: sr.SelectorInfo().Selector,
			},
			result.Found(),
		), nil
	}

	// If the values are empty, then return the current url
	return NewStringResult(
		targetURL.String(),
		SelectorInfo{
			Attr:     "href",
			InMeta:   false,
			Selector: "content",
		},
		true,
	), nil
}
