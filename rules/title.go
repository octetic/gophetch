package rules

// TitleRule is the rule for extracting the title information from a page.
type TitleRule struct {
	BaseRule
}

func NewTitleRule() *TitleRule {
	return &TitleRule{
		BaseRule: BaseRule{
			Strategies: titleStrategies,
		},
	}
}

var titleStrategies = []ExtractionStrategy{
	{
		Selectors: []string{"meta[property='og:title']", "meta[name='twitter:title']", "meta[property='twitter:title']"},
		Extractor: ExtractAttr("content"),
	},
	{
		Selectors: []string{"title"},
		Extractor: ExtractCSS,
	},
	{
		Selectors: []string{"headline"},
		Extractor: ExtractJSONLD,
	},
	{
		Selectors: []string{".post-title", ".entry-title", "h1[class*='title' i] a", "h1[class*='title' i]"},
		Extractor: ExtractCSS,
	},
}
