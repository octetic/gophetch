package rules

import (
	"net/url"

	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
)

// ReadableRule is the rule for extracting the readable content
type ReadableRule struct {
	BaseRule
}

// NewReadableRule creates a new ReadableRule
func NewReadableRule() *ReadableRule {
	return &ReadableRule{
		BaseRule: BaseRule{
			Strategies: readableStrategies,
		},
	}
}

var readableStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			`#readability-page-1`,
		},
		Extractor: extractReadable,
	},
}

// extractReadable extracts the readable content from the node
// It uses readability's article extraction routines
// Returns the excerpt, html content, and text content in that order
func extractReadable(node *html.Node, targetURL *url.URL, _ []string) ExtractResult {
	// Using readability's article extraction routines as they are more reliable than ours
	readabilityArticle, err := readability.FromDocument(node, targetURL)
	if err != nil {
		return ExtractResult{}
	}

	excerpt := readabilityArticle.Excerpt
	if len(excerpt) > 255 {
		excerpt = excerpt[:255] + "..."
	}

	htmlContent := readabilityArticle.Content
	textContent := readabilityArticle.TextContent
	image := readabilityArticle.Image
	lang := readabilityArticle.Language
	title := readabilityArticle.Title
	byline := readabilityArticle.Byline
	siteName := readabilityArticle.SiteName

	value := []string{excerpt, htmlContent, textContent, image, lang, title, byline, siteName}
	return ExtractResult{
		Value: value,
		Selector: SelectorInfo{
			Attr:     "readable",
			InMeta:   false,
			Selector: "readable",
		},
		Found: true,
	}
}
