package fetchers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pixiesys/gophetch/metadata"
)

// StandardHTTPFetcher is the struct that encapsulates the standard HTTP fetcher using the standard library.
// It does not support metadata.
type StandardHTTPFetcher struct{}

func (s *StandardHTTPFetcher) Name() string {
	return "standard"
}

func (s *StandardHTTPFetcher) FetchHTML(url string) (*http.Response, io.ReadCloser, error) {
	fmt.Println("Fetching HTML from Standard HTTP")

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}

	return resp, resp.Body, nil
}

func (s *StandardHTTPFetcher) HasMetadata() bool {
	return false
}

func (s *StandardHTTPFetcher) Metadata() metadata.Metadata {
	return metadata.Metadata{}
}
