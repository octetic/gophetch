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

// Split the feedStrategies into a separate strategies, so we can extract all available feeds
var feedStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			"link[type='application/rss+xml']",
		},
		Extractor: ExtractAttr("href"),
	},
	{
		Selectors: []string{
			"link[type='application/feed+json']",
		},
		Extractor: ExtractAttr("href"),
	},
	{
		Selectors: []string{
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
		if result.Found() {
			mvr := result.(*StringResult)
			feeds = append(feeds, mvr.value)
		}
	}

	if len(feeds) == 0 {
		return NewNoResult(), ErrValueNotFound
	}

	return NewMultiStringResult(
		feeds,
		SelectorInfo{
			Attr:     "href",
			InMeta:   false,
			Selector: "feed",
		},
		len(feeds) > 0,
	), nil
}
