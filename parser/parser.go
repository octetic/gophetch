package parser

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Headers map[string][]string

type Parser struct {
	reader   io.Reader
	response *http.Response
	headers  Headers
	node     *html.Node
	url      *url.URL
}

func New() *Parser {
	return &Parser{}
}

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

func (p *Parser) Node() *html.Node {
	return p.node
}

func (p *Parser) Headers() Headers {
	return p.headers
}

func (p *Parser) URL() *url.URL {
	return p.url
}

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
