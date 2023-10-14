package rules

// PublisherRule is the rule for extracting the publisher information from a page.
type PublisherRule struct {
	BaseRule
}

func NewPublisherRule() *PublisherRule {
	return &PublisherRule{
		BaseRule: BaseRule{
			Strategies: publisherStrategies,
		},
	}
}

var publisherStrategies = []ExtractionStrategy{
	{
		Selectors: []string{"publisher.name", "brand.name"},
		Extractor: ExtractJSONLD,
	},
	{
		Selectors: []string{
			"meta[property='og:site_name']",
			"meta[name*='application-name']",
			"meta[name*='app-title']",
			"meta[property*='app_name']",
			"meta[name='publisher']",
			"meta[name='twitter:app:name:iphone']",
			"meta[property='twitter:app:name:iphone']",
			"meta[name='twitter:app:name:ipad']",
			"meta[property='twitter:app:name:ipad']",
			"meta[name='twitter:app:name:googleplay']",
			"meta[property='twitter:app:name:googleplay']",
		},
		Extractor: ExtractMeta,
	},
	{
		Selectors: []string{
			"#logo",
			".logo",
			"a[class*='brand']",
			"[class*='brand']",
		},
		Extractor: ExtractCSS,
	},
	{
		Selectors: []string{
			"[class*='logo'] a img[alt]",
			"[class*='logo'] img[alt]",
		},
		Extractor: ExtractAttr("alt"),
	},
}
