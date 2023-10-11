package sites

import "github.com/pixiesys/gophetch/rules"

type YouTube struct{}

func (yt YouTube) DomainKey() string {
	return "youtube.com"
}

func (yt YouTube) Rules() map[string]rules.Rule {
	return map[string]rules.Rule{
		"author": rules.NewAuthorRule(
			rules.ExtractionStrategy{
				Selectors: []string{
					`[class*="user-info"]`,
				},
				Extractor: rules.ExtractCSS,
			},
			rules.ExtractionStrategy{
				Selectors: []string{
					`[itemprop="author"] [itemprop="name"]`,
				},
				Extractor: rules.ExtractAttr("content"),
			},
		),
		"image": rules.NewDateRule(rules.ExtractionStrategy{
			Selectors: []string{},
			Extractor: rules.ExtractAttr("content"),
		}),
	}
}
