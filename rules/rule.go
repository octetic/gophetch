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
	Extract(node *html.Node, targetURL *url.URL) ([]string, error)
}

// ExtractFunc is the function signature for all extractors
// It takes in a node, a target url, and a list of selectors
type ExtractFunc func(node *html.Node, targetURL *url.URL, selectors []string) ([]string, bool)

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
func (br *BaseRule) Extract(node *html.Node, targetURL *url.URL) ([]string, error) {
	for _, strategy := range br.Strategies {
		if value, found := strategy.Extractor(node, targetURL, strategy.Selectors); found {
			return value, nil
		}
	}

	return []string{}, ErrValueNotFound
}
