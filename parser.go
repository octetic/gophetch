package gophetch

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Headers is a map of HTTP headers
type Headers map[string][]string

// Parser is the struct that encapsulates the HTML parser.
type Parser struct {
	reader   io.Reader
	response *http.Response
	headers  Headers
	node     *html.Node
	url      *url.URL
}

// NewParser creates a new Parser struct.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses the HTML content from the provided io.Reader, and encapsulates the parsed HTML into a html.Node struct.
// It will also parse the HTTP headers from the provided http.Response struct. The targetURL parameter is used to fix
// relative paths.
func (p *Parser) Parse(reader io.Reader, resp *http.Response, targetURL string) error {
	u, err := url.Parse(targetURL)
	if err != nil {
		return err
	}

	p.url = u
	p.reader = reader
	p.response = resp
	p.headers = p.parseHeaders()

	doc, err := html.Parse(p.reader)
	if err != nil {
		return err
	}
	p.node = doc
	return nil
}

func (p *Parser) parseHeaders() map[string][]string {
	if p.response == nil {
		return nil
	}

	headers := make(map[string][]string)

	for k, v := range p.response.Header {
		headers[k] = v
	}

	p.headers = headers
	return headers
}

// Node returns the parsed HTML as a html.Node struct.
func (p *Parser) Node() *html.Node {
	return p.node
}

// Headers returns the HTTP headers as a map.
func (p *Parser) Headers() Headers {
	return p.headers
}

// URL returns the target URL as a url.URL struct.
func (p *Parser) URL() *url.URL {
	return p.url
}

// IsHTML returns true if the response is HTML, false otherwise.
func (p *Parser) IsHTML() bool {
	if p.headers == nil {
		return false
	}
	cp, ok := p.headers["Content-Type"]
	return ok && strings.Contains(cp[0], "text/html")
}

// MimeType returns the MIME type of the response
func (p *Parser) MimeType() string {
	if p.headers == nil {
		return ""
	}
	cp, ok := p.headers["Content-Type"]
	if !ok {
		return ""
	}
	return cp[0]
}
