package sites

import "github.com/octetic/gophetch/rules"

// YouTube is a site overrider for YouTube data extraction.
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
