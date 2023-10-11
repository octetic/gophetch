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
			return ExtractResult{}, err
		}
	}

	inMeta := result.Selector.Selector == "content"

	// If the value is not a full URL, then prepend the scheme and host
	for i, value := range result.Value {
		if !strings.HasPrefix(value, "http") {
			result.Value[i] = targetURL.Scheme + "://" + targetURL.Host + value
		}
	}

	result = ExtractResult{
		Value: result.Value,
		Selector: SelectorInfo{
			Attr:     result.Selector.Attr,
			InMeta:   inMeta,
			Selector: result.Selector.Selector,
		},
		Found: true,
	}

	if len(result.Value) > 0 {
		return result, nil
	}

	// If the values are empty, then return the current url
	result.Value = []string{targetURL.String()}
	return result, nil
}
