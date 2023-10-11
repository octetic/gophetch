package rules

type AuthorRule struct {
	BaseRule
}

func NewAuthorRule(strategies ...ExtractionStrategy) *AuthorRule {
	innerStrategies := authorStrategies
	if len(strategies) > 0 {
		innerStrategies = strategies
	}
	return &AuthorRule{
		BaseRule: BaseRule{
			Strategies: innerStrategies,
		},
	}
}

var authorStrategies = []ExtractionStrategy{
	{
		Selectors: []string{"author.name", "brand.name", "creator.name"},
		Extractor: ExtractJSONLD,
	},
	{
		Selectors: []string{
			"meta[name='author']",
			"meta[property='article:author']",
			"meta[property='dc:creator']",
			`meta[property="schema:author"]`,
			`meta[property="dc:creator"]`,
			`meta[itemprop="author"]`,
		},
		Extractor: ExtractAttr("content"),
	},
	{
		Selectors: []string{
			// RDFa Selectors
			`span[property="schema:author"]`,
			`div[typeof="schema:Person"] span[property="schema:name"]`,
			`span[property="dc:creator"]`,
			`div[typeof="dc:Person"] span[property="dc:name"]`,

			// Microdata Selectors
			`span[itemprop="author"]`,
			`div[itemtype="http://schema.org/Person"] span[itemprop="name"]`,

			// Common class or ID-based selectors
			`span[class="author"]`,
			`a[rel="author"]`,
			`span[id="author"]`,
		},
		Extractor: ExtractCSS,
	},
}
