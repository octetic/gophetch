package rules

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
		Extractor: ExtractAttr("content"),
	},
	{
		Selectors: []string{
			"html",
		},
		Extractor: ExtractAttr("lang"),
	},
}
