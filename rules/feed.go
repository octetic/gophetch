package rules

import (
	"net/url"

	"golang.org/x/net/html"
)

// FeedRule is the rule for extracting the feed URL of a page. It will respond with an array of feed URLs it found.
type FeedRule struct {
	BaseRule
}

func NewFeedRule() *FeedRule {
	return &FeedRule{
		BaseRule: BaseRule{
			Strategies: feedStrategies,
		},
	}
}

var feedStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			"link[type='application/rss+xml']",
			"link[type='application/feed+json']",
			"link[type='application/atom+xml']",
		},
		Extractor: ExtractAttr("href"),
	},
}

// Extract extracts the value from the node
func (fr *FeedRule) Extract(node *html.Node, targetURL *url.URL) (ExtractResult, error) {
	var feeds []string

	for _, strategy := range fr.Strategies {
		result := strategy.Extractor(node, targetURL, strategy.Selectors)
		if result.Found {
			// For each value found, fix the URL if it's relative and add it to the list of feeds.
			for i, v := range result.Value {
				result.Value[i] = FixRelativePath(targetURL, v)
			}
			feeds = append(feeds, result.Value...)
		}
	}

	return ExtractResult{
		Value: feeds,
		Selector: SelectorInfo{
			Attr:     "href",
			InMeta:   false,
			Selector: "feed",
		},
		Found: len(feeds) > 0,
	}, nil
}
