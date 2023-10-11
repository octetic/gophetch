package rules

import (
	"errors"
	"net/url"

	"golang.org/x/net/html"
)

var ErrValueNotFound = errors.New("no value found")

// Rule is the interface for all rules
type Rule interface {
	// Extract extracts the value from the node
	Extract(node *html.Node, targetURL *url.URL) (ExtractResult, error)
}

// ExtractResult is the result of an extraction
type ExtractResult struct {
	Value    []string
	Selector SelectorInfo
	Found    bool
}

// ExtractFunc is the function signature for all extractors
// It accepts the node to extract from, the target URL, and the selectors to use
// It returns the value as an array of strings, a string indicating where it was found, and a boolean indicating if the value was found
type ExtractFunc func(node *html.Node, targetURL *url.URL, selectors []string) ExtractResult

// ExtractionStrategy is the strategy for extracting a value
type ExtractionStrategy struct {
	Selectors []string
	Extractor ExtractFunc
}

// BaseRule is the base rule for all rules
type BaseRule struct {
	Strategies []ExtractionStrategy
}

// Extract extracts the value from the node
// It iterates through all the strategies and returns the first value found
func (br *BaseRule) Extract(node *html.Node, targetURL *url.URL) (ExtractResult, error) {
	for _, strategy := range br.Strategies {
		result := strategy.Extractor(node, targetURL, strategy.Selectors)
		if result.Found {
			return result, nil
		}
	}
	return ExtractResult{}, ErrValueNotFound
}
