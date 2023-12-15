// Package gophetch is a library for fetching and extracting metadata from HTML pages.
package gophetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/octetic/gophetch/fetchers"
	"github.com/octetic/gophetch/helpers"
	"github.com/octetic/gophetch/metadata"
	"github.com/octetic/gophetch/sites"
)

// Gophetch is the main struct that encapsulates the parser, extractor, and fetchers.
type Gophetch struct {
	Parser       *Parser
	Extractor    *Extractor
	Fetchers     []fetchers.HTMLFetcher
	SiteRegistry map[string]sites.Site
}

// Result is the struct that encapsulates the extracted metadata, along with the response data.
type Result struct {
	HTMLNode    *html.Node
	Headers     map[string][]string
	IsHTML      bool
	Metadata    metadata.Metadata
	MimeType    string
	Response    *http.Response
	StatusCode  int
	FetcherName string
}

// New creates a new Gophetch struct with the provided fetchers.
func New(fetchers ...fetchers.HTMLFetcher) *Gophetch {
	g := &Gophetch{
		Parser:       NewParser(),
		Extractor:    NewExtractor(),
		Fetchers:     fetchers,
		SiteRegistry: make(map[string]sites.Site),
	}
	g.RegisterSite(sites.YouTube{})
	return g
}

// ReadAndParse accepts two parameters: an io.Reader containing the HTML to be parsed, and a
// target URL string. It reads the HTML content from the provided io.Reader, parses it to extract metadata, and
// encapsulates the extracted metadata, along with the response data, into a Result struct which is then returned.
// This method is useful when the HTML content is already available and does not need to be fetched from the internet.
func (g *Gophetch) ReadAndParse(r io.Reader, targetURL string) (Result, error) {
	err := g.Parser.Parse(r, nil, targetURL)
	if err != nil {
		return Result{
			HTMLNode:    g.Parser.Node(),
			Headers:     g.Parser.Headers(),
			IsHTML:      g.Parser.IsHTML(),
			MimeType:    g.Parser.MimeType(),
			Response:    nil,
			StatusCode:  0,
			FetcherName: "",
		}, err
	}

	fetchedData := Result{
		HTMLNode:    g.Parser.Node(),
		Headers:     g.Parser.Headers(),
		IsHTML:      g.Parser.IsHTML(),
		MimeType:    g.Parser.MimeType(),
		Response:    nil,
		StatusCode:  0,
		FetcherName: "",
	}

	data, err := g.Extractor.ExtractMetadata(g.Parser.Node(), g.Parser.URL())
	if err != nil {
		return fetchedData, err
	}

	fetchedData.Metadata = data
	return fetchedData, nil
}

// FetchAndParse accepts a target URL string as its parameter. It initiates an HTTP request to fetch
// the HTML content from the specified URL, parses the fetched HTML to extract metadata, and encapsulates the
// extracted metadata, along with the response data, into a Result struct which is then returned. This method is
// useful when the HTML content needs to be fetched from the internet before parsing.
func (g *Gophetch) FetchAndParse(targetURL string) (Result, error) {
	var err error
	var body io.ReadCloser
	var resp *http.Response
	var fetcherName string

	// If no fetchers are provided, use the standard HTTP fetcher
	if len(g.Fetchers) == 0 {
		g.Fetchers = append(g.Fetchers, &fetchers.StandardHTTPFetcher{})
	}

	var data metadata.Metadata
	hasMetadata := false
	for _, fetcher := range g.Fetchers {
		resp, body, err = fetcher.FetchHTML(targetURL)
		if err == nil {
			data = fetcher.Metadata()
			hasMetadata = fetcher.HasMetadata()
			fetcherName = fetcher.Name()
			break
		}
	}

	if err != nil {
		return Result{}, err
	} else if resp == nil || body == nil {
		return Result{}, fmt.Errorf("unable to fetch HTML from %s", targetURL)
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)

	err = g.Parser.Parse(body, resp, targetURL)
	if err != nil {
		return Result{}, err
	}

	fetchedData := Result{
		HTMLNode:    g.Parser.Node(),
		Headers:     g.Parser.Headers(),
		IsHTML:      g.Parser.IsHTML(),
		Metadata:    metadata.Metadata{},
		MimeType:    g.Parser.MimeType(),
		Response:    resp,
		StatusCode:  resp.StatusCode,
		FetcherName: fetcherName,
	}

	// If the fetcher provided metadata, use that instead
	if hasMetadata {
		fetchedData.Metadata = data
		result, err := g.Extractor.ExtractRuleByKey(g.Parser.Node(), g.Parser.URL(), "readable")
		if err == nil {
			result.ApplyMetadata("readable", g.Parser.URL(), &fetchedData.Metadata)
		}
		result2, err := g.Extractor.ExtractRuleByKey(g.Parser.Node(), g.Parser.URL(), "lead_image")
		if err == nil {
			result2.ApplyMetadata("lead_image", g.Parser.URL(), &fetchedData.Metadata)
		}
		return fetchedData, nil
	}

	domain, err := ExtractDomain(targetURL)
	if err != nil {
		return fetchedData, err
	}

	if site, found := g.findSite(domain); found {
		g.Extractor.ApplySiteSpecificRules(site)
	}

	data, err = g.Extractor.ExtractMetadata(g.Parser.Node(), g.Parser.URL())

	if err != nil {
		return fetchedData, err
	}

	data.CleanURL = helpers.CleanURL(g.Parser.URL())
	fetchedData.Metadata = data
	return fetchedData, nil
}

// RegisterSite registers a site with the Gophetch instance. This allows the Gophetch instance to apply
// site-specific rules when extracting metadata from the HTML content.
func (g *Gophetch) RegisterSite(site sites.Site) {
	g.SiteRegistry[site.DomainKey()] = site
}

func (g *Gophetch) findSite(domain string) (sites.Site, bool) {
	site, found := g.SiteRegistry[domain]
	return site, found
}

// ExtractDomain extracts the domain from a given URL string
func ExtractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// This will extract "www.example.com" from "https://www.example.com/path"
	domain := u.Hostname()

	// Optionally, remove 'www.' prefix?
	domain = strings.TrimPrefix(domain, "www.")

	return domain, nil
}
