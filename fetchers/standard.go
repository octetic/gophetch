package fetchers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pixiesys/gophetch/metadata"
)

type StandardHTTPFetcher struct{}

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
