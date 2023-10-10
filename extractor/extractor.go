package extractor

import (
	"fmt"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"

	"github.com/pixiesys/gophetch/metadata"
	"github.com/pixiesys/gophetch/rules"
	"github.com/pixiesys/gophetch/sites"
)

type Extractor struct {
	Rules  map[string]rules.Rule
	Errors []error
}

func New() *Extractor {
	return &Extractor{
		Rules: map[string]rules.Rule{
			"author":      rules.NewAuthorRule(),
			"date":        rules.NewDateRule(),
			"description": rules.NewDescriptionRule(),
			"title":       rules.NewTitleRule(),
		},
	}
}

func (e *Extractor) ExtractMetadata(node *html.Node) (metadata.Metadata, error) {
	var meta metadata.Metadata

	if node == nil {
		return metadata.Metadata{}, fmt.Errorf("node is nil")
	}

	// Get the HTML as a string
	var sb strings.Builder
	err := html.Render(&sb, node)
	if err != nil {
		return metadata.Metadata{}, err
	}
	meta.HTML = sb.String()

	for key, rule := range e.Rules {
		value, err := rule.Extract(node)
		if err != nil {
			e.Errors = append(e.Errors, err)
			continue
		}

		switch key {
		case "author":
			meta.Author = Normalize(value)
		case "date":
			meta.Date = Normalize(value)
		case "description":
			meta.Description = Normalize(value)
		case "title":
			meta.Title = Normalize(value)
		default:
			meta.Dynamic[key] = value
		}
	}
	return meta, nil
}

func (e *Extractor) ApplySiteSpecificRules(site sites.Site) {
	for key, customRule := range site.Rules() {
		// Replace the default rule with the custom one for this key
		e.Rules[key] = customRule
	}
}

// Normalize cleans up the extracted string, removing HTML tags,
// decoding HTML entities, and trimming whitespace.
func Normalize(input string) string {
	// Strip HTML tags
	p := bluemonday.StripTagsPolicy()
	clean := p.Sanitize(input)

	// Decode HTML entities
	decoded := html.UnescapeString(clean)

	// Trim whitespace
	normalized := strings.TrimSpace(decoded)

	return normalized
}
