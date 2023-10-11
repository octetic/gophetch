package rules

import (
	"net/url"

	"golang.org/x/net/html"
)

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
func (fr *FeedRule) Extract(node *html.Node, targetURL *url.URL) ([]string, error) {
	var feeds []string

	for _, strategy := range fr.Strategies {
		if value, found := strategy.Extractor(node, targetURL, strategy.Selectors); found {
			// For each value found, fix the URL if it's relative and add it to the list of feeds.
			for i, v := range value {
				value[i] = FixRelativePath(targetURL, v)
			}
			feeds = append(feeds, value...)
		}
	}

	return feeds, ErrValueNotFound
}
