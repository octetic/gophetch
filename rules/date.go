package rules

type DateRule struct {
	BaseRule
}

func NewDateRule(strategies ...ExtractionStrategy) *DateRule {
	innerStrategies := dateStrategies
	if len(strategies) > 0 {
		innerStrategies = strategies
	}
	return &DateRule{
		BaseRule: BaseRule{
			Strategies: innerStrategies,
		},
	}
}

var dateStrategies = []ExtractionStrategy{
	// JSON LD
	{
		Selectors: []string{"datePublished", "dateCreated", "dateModified"},
		Extractor: ExtractJSONLD,
	},

	// Meta selectors
	{
		Selectors: []string{
			"meta[property='article:published_time']",
			"meta[property*='published_time']",
			"[itemprop*='datepublished']",
			"meta[property='og:published_time']",
			"meta[name='article:published_time']",
			"meta[name='og:published_time']",
			"meta[property*='modified_time']",
			"[itemprop*='datemodified']",
			"[itemprop*='date']",
		},
		Extractor: ExtractMeta,
	},

	// Time selectors
	{
		Selectors: []string{
			"time[itemprop*='date']",
			"time[datetime]",
		},
		Extractor: ExtractTime,
	},

	// Common CSS selectors
	{
		Selectors: []string{
			".post-date",
			".entry-date",
			".article-date",
			"[id*='date']",
			"[class*='date']",
			"[class*='time']",
		},
		Extractor: ExtractCSS,
	},
}
