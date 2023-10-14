package rules

// DescriptionRule is the rule for extracting the description information from a page.
type DescriptionRule struct {
	BaseRule
}

func NewDescriptionRule() *DescriptionRule {
	return &DescriptionRule{
		BaseRule: BaseRule{
			Strategies: descriptionStrategies,
		},
	}
}

var descriptionStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			"meta[property='og:description']",
			"meta[name='twitter:description']",
			"meta[property='twitter:description']",
			"meta[name='description']",
			"meta[itemprop='description']",
		},
		Extractor: ExtractMeta,
	},

	// JSON LD
	{
		Selectors: []string{"description", "articleBody"},
		Extractor: ExtractJSONLD,
	},

	// Common CSS selectors
	{
		Selectors: []string{
			".post-description",
			".entry-description",
			".article-description",
			".post-content p",
			".entry-content p",
			".article-content p",
			".post-content",
			".entry-content",
			".article-content",
			".post-body",
			".entry-body",
			".article-body",
			".post",
			".entry",
		},
		Extractor: ExtractCSS,
	},
}
