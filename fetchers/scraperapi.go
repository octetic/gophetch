package fetchers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/octetic/gophetch/metadata"
)

// ScraperapiFetcher is the struct that encapsulates the scraperapi.com fetcher. It is responsible for fetching the
// HTML from the given URL. This fetcher does not return metadata.
type ScraperapiFetcher struct {
	APIKey string
}

func (s *ScraperapiFetcher) Name() string {
	return "scraperapi"
}

func (s *ScraperapiFetcher) FetchHTML(targetURL string) (*http.Response, io.ReadCloser, error) {
	fmt.Println("Fetching HTML from Scraperapi")
	const endpoint = "https://api.scraperapi.com"

	// Create the URL with the token in the query parameters
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, nil, err
	}

	q := u.Query()
	q.Add("api_key", s.APIKey)
	q.Add("url", targetURL)
	q.Add("render", "true")
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

func (s *ScraperapiFetcher) HasMetadata() bool {
	return false
}

func (s *ScraperapiFetcher) Metadata() metadata.Metadata {
	return metadata.Metadata{}
}
