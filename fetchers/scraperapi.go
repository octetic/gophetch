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
	APIKey          string
	RenderJS        *bool   // default: true
	FollowRedirects *bool   // default: true
	Autoparse       *bool   // default: false
	Retry404        *bool   // default: false
	WaitForSelector *string // default: ""
	CountryCode     *string // default: nil
	DeviceType      *string // default: nil
	SessionNumber   *int    // default: nil
	BinaryTarget    *bool   // default: false
	UseOwnHeaders   *bool   // default: false
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
	if s.RenderJS == nil {
		q.Add("render", "true")
	} else {
		q.Add("render", fmt.Sprintf("%t", *s.RenderJS))
	}
	if s.FollowRedirects != nil {
		q.Add("follow_redirect", fmt.Sprintf("%t", *s.FollowRedirects))
	}
	if s.Autoparse != nil {
		q.Add("autoparse", fmt.Sprintf("%t", *s.Autoparse))
	}
	if s.Retry404 != nil {
		q.Add("retry_404", fmt.Sprintf("%t", *s.Retry404))
	}
	if s.WaitForSelector != nil {
		q.Add("wait_for_selector", *s.WaitForSelector)
	}
	if s.CountryCode != nil {
		q.Add("country_code", *s.CountryCode)
	}
	if s.DeviceType != nil {
		q.Add("device_type", *s.DeviceType)
	}
	if s.SessionNumber != nil {
		q.Add("session_number", fmt.Sprintf("%d", *s.SessionNumber))
	}
	if s.BinaryTarget != nil {
		q.Add("binary_content", fmt.Sprintf("%t", *s.BinaryTarget))
	}
	if s.UseOwnHeaders != nil {
		q.Add("use_own_headers", fmt.Sprintf("%t", *s.UseOwnHeaders))
	}

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
