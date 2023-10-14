package rules

// LangRule is the rule for extracting the language information from a page.
type LangRule struct {
	BaseRule
}

func NewLangRule() *LangRule {
	return &LangRule{
		BaseRule: BaseRule{
			Strategies: langStrategies,
		},
	}
}

var langStrategies = []ExtractionStrategy{
	{
		Selectors: []string{
			"meta[property='og:locale']",
			"meta[itemprop='inLanguage']",
		},
		Extractor: ExtractMeta,
	},
	{
		Selectors: []string{
			"html",
		},
		Extractor: ExtractAttr("lang"),
	},
}
