package parser

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Headers map[string][]string

type Parser struct {
	reader   io.Reader
	response *http.Response
	headers  Headers
	node     *html.Node
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(rd io.Reader, re *http.Response) error {
	p.reader = rd
	p.response = re
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
