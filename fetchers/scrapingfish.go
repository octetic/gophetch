package fetchers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/minsoft-io/gophetch/metadata"
)

// ScrapingfishFetcher is the struct that encapsulates the scrapingfish.com fetcher. It is responsible for fetching the
// HTML from the given URL. This fetcher does not return metadata.
type ScrapingfishFetcher struct {
	APIKey string
}

func (s *ScrapingfishFetcher) Name() string {
	return "scrapingfish"
}

func (s *ScrapingfishFetcher) FetchHTML(targetURL string) (*http.Response, io.ReadCloser, error) {
	fmt.Println("Fetching HTML from Scrapingfish")
	const endpoint = "https://scraping.narf.ai/api/v1/"

	// Create the URL with the token in the query parameters
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	q.Add("api_key", s.APIKey)
	q.Add("url", targetURL)
	q.Add("render_js", "true")
	u.RawQuery = q.Encode()
	urlString := u.String()

	// Create a new HTTP request
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, nil, err
	}

	// Create a new HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return resp, resp.Body, nil
}

func (s *ScrapingfishFetcher) HasMetadata() bool {
	return false
}

func (s *ScrapingfishFetcher) Metadata() metadata.Metadata {
	return metadata.Metadata{}
}
