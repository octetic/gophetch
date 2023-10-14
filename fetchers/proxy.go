package fetchers

import (
	"io"
	"net/http"
	"net/url"

	"github.com/minsoft-io/gophetch/metadata"
)

// ProxyHTTPFetcher is a fetcher that uses a proxy to fetch HTML from a URL. It does not return metadata.
type ProxyHTTPFetcher struct {
	ProxyURL string
}

// Name returns the name of the fetcher.
func (p *ProxyHTTPFetcher) Name() string {
	return "proxy"
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

func (p *ProxyHTTPFetcher) HasMetadata() bool {
	return false
}

func (p *ProxyHTTPFetcher) Metadata() metadata.Metadata {
	return metadata.Metadata{}
}
