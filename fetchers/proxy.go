package fetchers

import (
	"io"
	"net/http"
	"net/url"
)

// ProxyHTTPFetcher is a fetcher that uses a proxy to fetch HTML from a URL.
type ProxyHTTPFetcher struct {
	ProxyURL string
}

func (p *ProxyHTTPFetcher) FetchHTML(targetURL string) (io.ReadCloser, error) {
	proxyURL, err := url.Parse(p.ProxyURL)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
