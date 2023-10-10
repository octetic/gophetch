package fetchers

import (
	"io"
	"net/http"

	"github.com/pixiesys/gophetch/metadata"
)

type HTMLFetcher interface {
	FetchHTML(url string) (*http.Response, io.ReadCloser, error)
	HasMetadata() bool
	Metadata() metadata.Metadata
}
