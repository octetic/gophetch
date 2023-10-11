package gophetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/pixiesys/gophetch/extractor"
	"github.com/pixiesys/gophetch/fetchers"
	"github.com/pixiesys/gophetch/metadata"
	"github.com/pixiesys/gophetch/parser"
	"github.com/pixiesys/gophetch/sites"
)

type Gophetch struct {
	Parser       *parser.Parser
	Extractor    *extractor.Extractor
	Fetchers     []fetchers.HTMLFetcher
	SiteRegistry map[string]sites.Site
}

type FetchedData struct {
	HTMLNode    *html.Node
	Headers     map[string][]string
	IsHTML      bool
	Metadata    metadata.Metadata
	MimeType    string
	Response    *http.Response
	StatusCode  int
	FetcherName string
}

func New(fetchers ...fetchers.HTMLFetcher) *Gophetch {
	g := &Gophetch{
		Parser:       parser.New(),
		Extractor:    extractor.New(),
		Fetchers:     fetchers,
		SiteRegistry: make(map[string]sites.Site),
	}
	g.RegisterSite(sites.YouTube{})
	return g
}

func (g *Gophetch) FetchAndExtractFromReader(r io.Reader, targetURL string) (FetchedData, error) {
	err := g.Parser.Parse(r, nil, targetURL)
	if err != nil {
		return FetchedData{
			HTMLNode:    g.Parser.Node(),
			Headers:     g.Parser.Headers(),
			IsHTML:      g.Parser.IsHTML(),
			MimeType:    g.Parser.MimeType(),
			Response:    nil,
			StatusCode:  0,
			FetcherName: "",
		}, err
	}

	fetchedData := FetchedData{
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

// FetchAndExtractFromURL fetches and extracts data from a URL using the available fetchers
func (g *Gophetch) FetchAndExtractFromURL(targetURL string) (FetchedData, error) {
	var err error
	var body io.ReadCloser
	var resp *http.Response
	var fetcherName string

	// If no fetchers are provided, use the standard HTTP fetcher
	if len(g.Fetchers) == 0 {
		g.Fetchers = append(g.Fetchers, &fetchers.StandardHTTPFetcher{})
	}

	var data metadata.Metadata
	hasData := false
	for _, fetcher := range g.Fetchers {
		resp, body, err = fetcher.FetchHTML(targetURL)
		if err == nil {
			data = fetcher.Metadata()
			hasData = fetcher.HasMetadata()
			fetcherName = fetcher.Name()
			break
		}
	}

	if err != nil {
		return FetchedData{}, err
	} else if resp == nil || body == nil {
		return FetchedData{}, fmt.Errorf("unable to fetch HTML from %s", targetURL)
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(body)

	err = g.Parser.Parse(body, resp, targetURL)
	if err != nil {
		return FetchedData{}, err
	}

	fetchedData := FetchedData{
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
	if hasData {
		fetchedData.Metadata = data
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

	fetchedData.Metadata = data
	return fetchedData, nil
}

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
