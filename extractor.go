package gophetch

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"

	"github.com/pixiesys/gophetch/metadata"
	"github.com/pixiesys/gophetch/rules"
	"github.com/pixiesys/gophetch/sites"
)

// Extractor is the struct that encapsulates the rules used to extract metadata from HTML.
type Extractor struct {
	Rules  map[string]rules.Rule
	Errors []error
}

// NewExtractor creates a new Extractor struct with the default rules.
func NewExtractor() *Extractor {
	return &Extractor{
		Rules: map[string]rules.Rule{
			"author":      rules.NewAuthorRule(),
			"canonical":   rules.NewCanonicalRule(),
			"date":        rules.NewDateRule(),
			"description": rules.NewDescriptionRule(),
			"favicon":     rules.NewFaviconRule(),
			"feed":        rules.NewFeedRule(),
			"lang":        rules.NewLangRule(),
			"lead_image":  rules.NewLeadImageRule(),
			"publisher":   rules.NewPublisherRule(),
			"readable":    rules.NewReadableRule(),
			"site_name":   rules.NewSiteNameRule(),
			"title":       rules.NewTitleRule(),
		},
	}
}

// ExtractMetadata extracts metadata from the given HTML node. The url parameter is used to fix relative paths.
func (e *Extractor) ExtractMetadata(node *html.Node, targetURL *url.URL) (metadata.Metadata, error) {
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
		//fmt.Printf("Extracting %s\n", key)
		result, err := rule.Extract(node, targetURL)
		if err != nil {
			e.Errors = append(e.Errors, err)
			continue
		} else if len(result.Value) == 0 || !result.Found {
			continue
		}

		value := result.Value

		switch key {
		case "author":
			meta.Author = Normalize(value[0])
		case "canonical":
			canonicalURL := rules.FixRelativePath(targetURL, value[0])
			meta.CanonicalURL = canonicalURL
			meta.URL = canonicalURL
		case "date":
			meta.Date = Normalize(value[0])
		case "description":
			meta.Description = Normalize(value[0])
		case "favicon":
			meta.FaviconURL = rules.FixRelativePath(targetURL, value[0])
		case "feed":
			meta.FeedURLs = value
		case "lang":
			meta.Lang = Normalize(value[0])
		case "lead_image":
			meta.LeadImageURL = rules.FixRelativePath(targetURL, value[0])
			meta.LeadImageInMeta = result.Selector.InMeta
		case "publisher":
			meta.Publisher = Normalize(value[0])
		case "readable":
			if len(value) != 9 {
				return metadata.Metadata{}, fmt.Errorf("readable rule returned %d values, expected 9", len(value))
			}
			meta.ReadableExcerpt = value[0]
			meta.ReadableHTML = value[1]
			meta.ReadableText = value[2]
			meta.ReadableImage = value[3]
			meta.ReadableLang = value[4]
			meta.ReadableLength = utf8.RuneCountInString(value[2])
			meta.ReadableTitle = value[5]
			meta.ReadableByline = value[6]
			meta.ReadableSiteName = value[7]
			meta.IsReadable = value[8] == "true"

		case "site_name":
			meta.SiteName = Normalize(value[0])
		case "title":
			meta.Title = Normalize(value[0])
		default:
			meta.Dynamic[key] = value
		}
	}
	return meta, nil
}

// ApplySiteSpecificRules applies the custom rules for the given site.
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
