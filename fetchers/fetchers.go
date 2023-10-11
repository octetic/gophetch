package fetchers

import (
	"io"
	"net/http"

	"github.com/pixiesys/gophetch/metadata"
)

type HTMLFetcher interface {
	// Name returns the name of the fetcher.
	Name() string
	// FetchHTML fetches the HTML from the given URL.
	FetchHTML(url string) (*http.Response, io.ReadCloser, error)
	// HasMetadata returns true if the fetcher has metadata.
	HasMetadata() bool
	// Metadata returns the metadata for the fetcher.
	Metadata() metadata.Metadata
}
