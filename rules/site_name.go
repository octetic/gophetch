package rules

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// SiteNameRule is the rule for extracting the site name information from a page.
type SiteNameRule struct {
	BaseRule
}

func NewSiteNameRule() *SiteNameRule {
	return &SiteNameRule{
		BaseRule: BaseRule{
			Strategies: siteNameStrategies,
		},
	}
}

var siteNameStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			"meta[property='og:site_name']",
			"meta[name='og:site_name']",
			"meta[property='twitter:site_name']",
			"meta[name='twitter:site_name']",
			"meta[itemprop='name']",
			"meta[name='application-name']",
		},
		Extractor: ExtractAttr("content"),
	},
}

func (r *SiteNameRule) Extract(node *html.Node, targetURL *url.URL) (ExtractResult, error) {
	result, err := r.BaseRule.Extract(node, targetURL)
	if err == nil && result.Found() {
		return result, nil
	}

	// If no site name was found, use the domain name without the TLD and 'www'
	domain := targetURL.Hostname()
	domain = strings.TrimPrefix(domain, "www.")
	return NewStringResult(
		domain,
		SelectorInfo{
			Attr:     "content",
			InMeta:   false,
			Selector: "domain",
		},
		true,
	), nil
}
