package rules

import (
	"errors"
	"net/url"

	"golang.org/x/net/html"

	"github.com/octetic/gophetch/helpers"
	"github.com/octetic/gophetch/metadata"
)

var ErrValueNotFound = errors.New("no value found")

// A Rule is a rule for extracting a value from a node. It encapsulates multiple strategies for extracting a value.
// Each strategy is tried in order of priority until a value is found, or all strategies have been tried.
type Rule interface {
	// Extract extracts the value from the node
	Extract(node *html.Node, targetURL *url.URL) (ExtractResult, error)
}

// ExtractFunc is the function signature for all extractors that can be used in a strategy.
// It accepts the node to extract from, the target URL, and the selectors to use
// It returns the value as an array of strings, a string indicating where it was found, and a boolean indicating if the value was found
type ExtractFunc func(node *html.Node, targetURL *url.URL, selectors []string) ExtractResult

// ExtractionStrategy is the strategy for extracting a value
type ExtractionStrategy struct {
	Selectors []string
	Extractor ExtractFunc
}

// BaseRule is the base rule for all rules
type BaseRule struct {
	Strategies []ExtractionStrategy
}

// Extract extracts the value from the node
// It iterates through all the strategies and returns the first value found
func (br *BaseRule) Extract(node *html.Node, targetURL *url.URL) (ExtractResult, error) {
	for _, strategy := range br.Strategies {
		result := strategy.Extractor(node, targetURL, strategy.Selectors)
		if result.Found() {
			return result, nil
		}
	}
	return NewNoResult(), ErrValueNotFound
}

// ----------------------------------------

// ExtractResult is the result of an extraction.
type ExtractResult interface {
	ApplyMetadata(key string, u *url.URL, m *metadata.Metadata)
	Found() bool
	SelectorInfo() SelectorInfo
	Value() any
}

type BaseResult struct {
	found        bool
	selectorInfo SelectorInfo
}

func (r *BaseResult) Found() bool {
	return r.found
}

func (r *BaseResult) SelectorInfo() SelectorInfo {
	return r.selectorInfo
}

type NoResult struct {
	*BaseResult
}

func NewNoResult() *NoResult {
	return &NoResult{
		BaseResult: &BaseResult{
			found: false,
		},
	}
}

func (r *BaseResult) ApplyMetadata(_ string, _ *url.URL, _ *metadata.Metadata) {
	// no-op
}

func (r *NoResult) Value() any {
	return nil
}

type StringResult struct {
	*BaseResult
	value string
}

func NewStringResult(value string, selectorInfo SelectorInfo, found bool) *StringResult {
	return &StringResult{
		BaseResult: &BaseResult{
			found:        found,
			selectorInfo: selectorInfo,
		},
		value: value,
	}
}

func (r *StringResult) ApplyMetadata(key string, u *url.URL, m *metadata.Metadata) {
	switch key {
	case "author":
		m.Author = helpers.Normalize(r.value)
	case "canonical":
		canonicalURL := helpers.FixRelativePath(u, r.value)
		m.CanonicalURL = canonicalURL
		m.URL = canonicalURL
	case "date":
		m.Date = helpers.Normalize(r.value)
	case "description":
		m.Description = helpers.Normalize(r.value)
	case "favicon":
		m.FaviconURL = helpers.FixRelativePath(u, r.value)
	case "lang":
		m.Lang = helpers.Normalize(r.value)
	case "lead_image":
		m.LeadImageURL = helpers.FixRelativePath(u, r.value)
		m.LeadImageInMeta = r.selectorInfo.InMeta
	case "publisher":
		m.Publisher = helpers.Normalize(r.value)
	case "site_name":
		m.SiteName = helpers.Normalize(r.value)
	case "title":
		m.Title = helpers.Normalize(r.value)
	default:
		m.Dynamic[key] = r.value
	}
}

func (r *StringResult) Value() any {
	return r.value
}

type MultiStringResult struct {
	*BaseResult
	value []string
}

func NewMultiStringResult(value []string, selectorInfo SelectorInfo, found bool) *MultiStringResult {
	return &MultiStringResult{
		BaseResult: &BaseResult{
			found:        found,
			selectorInfo: selectorInfo,
		},
		value: value,
	}
}

func (r *MultiStringResult) ApplyMetadata(key string, _ *url.URL, m *metadata.Metadata) {
	switch key {
	case "feed":
		m.FeedURLs = r.value
	default:
		m.Dynamic[key] = r.value
	}
}

func (r *MultiStringResult) Found() bool {
	return len(r.value) > 0 && r.BaseResult.Found()
}

func (r *MultiStringResult) Value() any {
	return r.value
}

type ReadableValue struct {
	Excerpt    string
	HTML       string
	Text       string
	Image      string
	Lang       string
	Length     int
	Title      string
	Byline     string
	SiteName   string
	IsReadable bool
}

type ReadableResult struct {
	*BaseResult
	value ReadableValue
}

func NewReadableResult(value ReadableValue, selectorInfo SelectorInfo, found bool) *ReadableResult {
	return &ReadableResult{
		BaseResult: &BaseResult{
			found:        found,
			selectorInfo: selectorInfo,
		},
		value: value,
	}
}

func (r *ReadableResult) ApplyMetadata(_ string, _ *url.URL, m *metadata.Metadata) {
	m.ReadableExcerpt = r.value.Excerpt
	m.ReadableHTML = r.value.HTML
	m.ReadableText = r.value.Text
	m.ReadableImage = r.value.Image
	m.ReadableLang = r.value.Lang
	m.ReadableLength = r.value.Length
	m.ReadableTitle = r.value.Title
	m.ReadableByline = r.value.Byline
	m.ReadableSiteName = r.value.SiteName
	m.IsReadable = r.value.IsReadable
}

func (r *ReadableResult) Value() any {
	return r.value
}
