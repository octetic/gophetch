package rules

import (
	"errors"
	_ "image/gif"  // This is required to initialize the GIF decoder
	_ "image/jpeg" // This is required to initialize the JPEG decoder
	_ "image/png"  // This is required to initialize the PNG decoder
)

var ErrInvalidImageFormat = errors.New("invalid image format")

// LeadImageRule is the rule for extracting the lead image from a page.
type LeadImageRule struct {
	BaseRule
}

func NewLeadImageRule() *LeadImageRule {
	return &LeadImageRule{
		BaseRule: BaseRule{
			Strategies: leadImageStrategies,
		},
	}
}

var leadImageStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			"meta[property='og:image:secure_url']",
			"meta[property='og:image:url']",
			"meta[property='og:image']",
			"meta[name='og:image']",
			"meta[name='twitter:image:src']",
			"meta[property='twitter:image:src']",
			"meta[name='twitter:image']",
			"meta[property='twitter:image']",
			"meta[itemprop='image']",
		},
		Extractor: ExtractAttr("content"),
	},
	{
		Selectors: []string{
			"img[src]:not([width='1']):not([height='1'])",
			"img[srcset]:not([width='1']):not([height='1'])",
			"img[data-src]:not([width='1']):not([height='1'])",
			"img[data-srcset]:not([width='1']):not([height='1'])",
			"img[data-lazy-src]:not([width='1']):not([height='1'])",
			"img[data-lazy-srcset]:not([width='1']):not([height='1'])",
			"img[data-lazyload]:not([width='1']):not([height='1'])",
		},
		Extractor: ExtractCSS,
	},
}
