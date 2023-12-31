package rules

import (
	"net/url"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

import (
	"encoding/json"
	// import other required packages
)

type SelectorInfo struct {
	Attr     string
	InMeta   bool
	Selector string
}

// ExtractJSONLD extracts the given JSON-LD attribute from the given document.
func ExtractJSONLD(node *html.Node, _ *url.URL, selectors []string) ExtractResult {
	jsonLdNodes := cascadia.QueryAll(node, cascadia.MustCompile(`script[type="application/ld+json"]`))
	for _, jsonLdNode := range jsonLdNodes {
		if jsonLdNode.FirstChild != nil {
			var obj map[string]any
			err := json.Unmarshal([]byte(jsonLdNode.FirstChild.Data), &obj)
			if err != nil {
				continue
			}

			for _, selector := range selectors {
				keys := strings.Split(selector, ".")
				val := obj
				for _, key := range keys {
					if nextVal, ok := val[key].(map[string]any); ok {
						val = nextVal
					} else if finalVal, ok := val[key].(string); ok {
						return NewStringResult(finalVal, SelectorInfo{
							Attr:     key,
							InMeta:   false,
							Selector: selector,
						}, true)
					} else {
						break
					}
				}
			}
		}
	}
	return NewNoResult()
}

// ExtractCSS extracts the given CSS selector from the given document.
func ExtractCSS(node *html.Node, _ *url.URL, selectors []string) ExtractResult {
	for _, selector := range selectors {
		cssNode := cascadia.Query(node, cascadia.MustCompile(selector))
		if cssNode != nil && cssNode.FirstChild != nil {
			return NewStringResult(
				strings.TrimSpace(cssNode.FirstChild.Data),
				SelectorInfo{
					Attr:     "text",
					InMeta:   false,
					Selector: selector,
				},
				true,
			)
		}
	}
	return NewNoResult()
}

// ExtractAttr extracts a selector from the given document using the given attribute.
func ExtractAttr(attribute string) ExtractFunc {
	return func(node *html.Node, _ *url.URL, selectors []string) ExtractResult {
		for _, selector := range selectors {
			cssNode := cascadia.Query(node, cascadia.MustCompile(selector))
			if cssNode != nil {
				for _, attr := range cssNode.Attr {
					if attr.Key == attribute {
						return NewStringResult(
							attr.Val,
							SelectorInfo{
								Attr:     attribute,
								InMeta:   false,
								Selector: selector,
							},
							true,
						)
					}
				}
			}
		}
		return NewNoResult()
	}
}

// ExtractMeta extracts the given meta tag from the given document.
func ExtractMeta(node *html.Node, targetURL *url.URL, selectors []string) ExtractResult {
	fn := ExtractAttr("content")
	result := fn(node, targetURL, selectors)

	// Check if the result is found, is not nil, and the stripped value is not empty
	if result.Found() && result.Value() != nil && strings.TrimSpace(result.Value().(string)) != "" {
		sr := result.(*StringResult)
		return NewStringResult(
			sr.value,
			SelectorInfo{
				Attr:     sr.SelectorInfo().Attr,
				InMeta:   true,
				Selector: sr.SelectorInfo().Selector,
			},
			result.Found(),
		)
	}
	return NewNoResult()
}
