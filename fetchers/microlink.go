package fetchers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/minsoft-io/gophetch/metadata"
)

const endpoint = "https://pro.microlink.io/"

type MicrolinkFetchedJSON struct {
	Status  string            `json:"status"`
	Data    metadata.Metadata `json:"data"`
	Message string            `json:"message"`
}

type MicrolinkDataQueryRule struct {
	Selector string `json:"selector"`
	Type     string `json:"type"`
	Attr     string `json:"attr"`
}

func (d MicrolinkDataQueryRule) AsMap(prefix string) map[string]string {
	return map[string]string{
		prefix + ".selector": d.Selector,
		prefix + ".type":     d.Type,
		prefix + ".attr":     d.Attr,
	}
}

// MicrolinkFetcher is the struct that encapsulates the microlink.io fetcher. It is responsible for fetching the
// HTML from the given URL. This fetcher returns metadata.
type MicrolinkFetcher struct {
	AdBlock   bool
	APIKey    string
	Prerender bool
	metadata  metadata.Metadata
}

func (m *MicrolinkFetcher) Name() string {
	return "microlink"
}

func (m *MicrolinkFetcher) FetchHTML(targetURL string) (*http.Response, io.ReadCloser, error) {
	fmt.Println("Fetching HTML from Microlink")

	// Create the URL with the token in the query parameters
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, nil, err
	}
	q := m.buildQueryParams(targetURL, map[string]MicrolinkDataQueryRule{
		"html": {
			Selector: "html",
			Type:     "string",
			Attr:     "html",
		},
	})

	u.RawQuery = q.Encode()
	urlString := u.String()

	// Create a new HTTP request
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("x-api-key", m.APIKey)

	// Create a new HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var fetchedJSON MicrolinkFetchedJSON
	err = json.Unmarshal(body, &fetchedJSON)
	if err != nil {
		return resp, nil, err
	}

	m.metadata = fetchedJSON.Data

	// Find the html content
	htmlContent := fetchedJSON.Data.HTML
	if htmlContent == "" {
		return resp, nil, fmt.Errorf("unable to find HTML content in response")
	}

	// Create a new response body
	respBody := io.NopCloser(strings.NewReader(htmlContent))

	return resp, respBody, nil
}

func (m *MicrolinkFetcher) HasMetadata() bool {
	return true
}

func (m *MicrolinkFetcher) Metadata() metadata.Metadata {
	return m.metadata
}

func (m *MicrolinkFetcher) buildQueryParams(urlPath string, rules map[string]MicrolinkDataQueryRule) url.Values {
	q := url.Values{}
	q.Add("url", urlPath)
	q.Add("adblock", fmt.Sprintf("%t", m.AdBlock))
	q.Add("prerender", fmt.Sprintf("%t", m.Prerender))
	for ruleName, rule := range rules {
		prefix := "data." + ruleName
		for key, value := range rule.AsMap(prefix) {
			q.Add(key, value)
		}
	}
	return q
}
