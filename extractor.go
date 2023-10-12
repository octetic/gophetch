package gophetch

import (
	"fmt"
	"net/url"
	"strings"

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

	doc, err := e.renderHTML(node)
	if err != nil {
		return metadata.Metadata{}, err
	}
	meta.HTML = doc

	for key, rule := range e.Rules {
		result, err := e.ExtractRule(node, targetURL, rule)
		if err != nil {
			e.handleError(err)
			continue
		} else if !result.Found() {
			continue
		}

		result.ApplyMetadata(key, targetURL, &meta)
		if err != nil {
			e.handleError(err)
		}
	}

	return meta, nil
}

func (e *Extractor) ExtractRuleByKey(node *html.Node, targetURL *url.URL, key string) (rules.ExtractResult, error) {
	rule, ok := e.Rules[key]
	if !ok {
		return rules.NewNoResult(), fmt.Errorf("rule %s not found", key)
	}
	return e.ExtractRule(node, targetURL, rule)
}

func (e *Extractor) ExtractRule(node *html.Node, targetURL *url.URL, rule rules.Rule) (rules.ExtractResult, error) {
	result, err := rule.Extract(node, targetURL)
	if err != nil {
		return rules.NewNoResult(), err
	}
	if !result.Found() {
		return rules.NewNoResult(), nil // or some sentinel error if needed
	}
	return result, nil
}

// ApplySiteSpecificRules applies the custom rules for the given site.
func (e *Extractor) ApplySiteSpecificRules(site sites.Site) {
	for key, customRule := range site.Rules() {
		// Replace the default rule with the custom one for this key
		e.Rules[key] = customRule
	}
}

func (e *Extractor) handleError(err error) {
	e.Errors = append(e.Errors, err)
}

func (e *Extractor) renderHTML(node *html.Node) (string, error) {
	var sb strings.Builder
	err := html.Render(&sb, node)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}
