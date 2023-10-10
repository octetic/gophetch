package sites

import "github.com/pixiesys/gophetch/rules"

type YouTube struct{}

func (yt YouTube) DomainKey() string {
	return "youtube.com"
}

func (yt YouTube) Rules() map[string]rules.Rule {
	return map[string]rules.Rule{
		"author": rules.NewAuthorRule(rules.ExtractionStrategy{
			Selectors: []string{`meta[property="og:site_name"]`},
			Extractor: rules.ExtractMeta,
		}),
		"date": rules.NewDateRule(rules.ExtractionStrategy{
			Selectors: []string{`meta[property="og:updated_time"]`},
			Extractor: rules.ExtractMeta,
		}),
	}
}
