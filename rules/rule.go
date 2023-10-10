package rules

import (
	"errors"

	"golang.org/x/net/html"
)

var ErrValueNotFound = errors.New("no value found")

// Rule is the interface for all rules
type Rule interface {
	// Extract extracts the value from the node
	Extract(node *html.Node) (string, error)
}

// ExtractFunc is the function signature for all extractors
type ExtractFunc func(*html.Node, []string) (string, bool)

// ExtractionStrategy is the strategy for extracting a value
type ExtractionStrategy struct {
	Selectors []string
	Extractor ExtractFunc
}

// BaseRule is the base rule for all rules
type BaseRule struct {
	Strategies []ExtractionStrategy
}

func (br *BaseRule) Extract(node *html.Node) (string, error) {
	for _, strategy := range br.Strategies {
		if value, found := strategy.Extractor(node, strategy.Selectors); found {
			return value, nil
		}
	}

	return "", ErrValueNotFound
}
