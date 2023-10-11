package rules

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
