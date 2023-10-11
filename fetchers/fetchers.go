package fetchers

import (
	"io"
	"net/http"

	"github.com/pixiesys/gophetch/metadata"
)

// HTMLFetcher is the interface that encapsulates the HTML fetcher. An HTML fetcher is responsible for fetching the
// HTML from the given URL. It also has the option to return metadata.
type HTMLFetcher interface {
	// Name returns the name of the fetcher.
	Name() string
	// FetchHTML fetches the HTML from the given URL.
	FetchHTML(url string) (*http.Response, io.ReadCloser, error)
	// HasMetadata returns true if the fetcher can return metadata.
	HasMetadata() bool
	// Metadata returns the metadata for the fetcher, if available.
	Metadata() metadata.Metadata
}
