package fetchers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/octetic/gophetch/metadata"
)

const browserlessEndpoint = "https://chrome.browserless.io"

type BrowserlessFetcher struct {
	APIToken    string
	GotoOptions BrowserlessGoToOptions
}

type BrowserlessMargin struct {
	Top    string `json:"top,omitempty"`
	Right  string `json:"right,omitempty"`
	Bottom string `json:"bottom,omitempty"`
	Left   string `json:"left,omitempty"`
}

type BrowserlessPdfOptions struct {
	DisplayHeaderFooter bool              `json:"displayHeaderFooter,omitempty"`
	FooterTemplate      string            `json:"footerTemplate,omitempty"`
	Format              string            `json:"format,omitempty"`
	Height              int               `json:"height,omitempty"`
	Landscape           bool              `json:"landscape,omitempty"`
	HeaderTemplate      string            `json:"headerTemplate,omitempty"`
	Margin              BrowserlessMargin `json:"margin,omitempty"`
	PageRanges          string            `json:"pageRanges,omitempty"`
	PreferCSSPageSize   bool              `json:"preferCSSPageSize,omitempty"`
	PrintBackground     bool              `json:"printBackground,omitempty"`
	Scale               float32           `json:"scale,omitempty"`
	Width               int               `json:"width,omitempty"`
	OmitBackground      bool              `json:"omitBackground,omitempty"`
	Timeout             int               `json:"timeout,omitempty"`
}

type BrowserlessGoToOptions struct {
	WaitUntil string `json:"waitUntil,omitempty"`
	Timeout   int    `json:"timeout,omitempty"`
}

type BrowserlessPdfRequest struct {
	URL  string                 `json:"url,omitempty"`
	Opts BrowserlessPdfOptions  `json:"options,omitempty"`
	GoTo BrowserlessGoToOptions `json:"gotoOptions,omitempty"`
}

type BrowserlessClip struct {
	Height int `json:"height,omitempty"`
	Width  int `json:"width,omitempty"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type BrowserlessScreenshotOptions struct {
	Clip           BrowserlessClip `json:"clip,omitempty"`
	FullPage       bool            `json:"fullPage,omitempty"`
	OmitBackground bool            `json:"omitBackground,omitempty"`
	Quality        int             `json:"quality,omitempty"`
	Type           string          `json:"type,omitempty"`
	Encoding       string          `json:"encoding,omitempty"`
}

type ScreenshotRequest struct {
	URL  string                       `json:"url,omitempty"`
	Opts BrowserlessScreenshotOptions `json:"options,omitempty"`
}

type ContentRequest struct {
	URL  string                 `json:"url,omitempty"`
	GoTo BrowserlessGoToOptions `json:"gotoOptions,omitempty"`
}

//func (c *Client) GetScreenshotBytes(targetURL string, options ScreenshotOptions) ([]byte, error) {
//	request := &ScreenshotRequest{
//		URL:  targetURL,
//		Opts: options,
//	}
//	return c.makeBrowserlessRequest("/screenshot", request)
//}
//
//func (c *Client) GetPdfBytes(targetURL string, options PdfOptions, gotoOptions GoToOptions) ([]byte, error) {
//	request := &PdfRequest{
//		URL:  targetURL,
//		Opts: options,
//		GoTo: gotoOptions,
//	}
//	return c.makeBrowserlessRequest("/pdf", request)
//}

func (b *BrowserlessFetcher) Name() string {
	return "browserless"
}

func (b *BrowserlessFetcher) FetchHTML(targetURL string) (*http.Response, io.ReadCloser, error) {
	request := &ContentRequest{
		URL:  targetURL,
		GoTo: b.GotoOptions,
	}
	resp, body, err := b.makeBrowserlessRequest("/content", request)
	if err != nil {
		return resp, nil, err
	}
	return resp, body, nil
}

func (b *BrowserlessFetcher) makeBrowserlessRequest(path string, request interface{}) (*http.Response, io.ReadCloser, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, nil, err
	}

	// Create the URL with the token in the query parameters
	apiPath := browserlessEndpoint + path
	u, err := url.Parse(apiPath)
	if err != nil {
		return nil, nil, err
	}
	q := u.Query()
	q.Add("token", b.APIToken)
	u.RawQuery = q.Encode()

	// Create a new HTTP request
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, err
	}

	// Add headers
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Content-Type", "application/json")

	// Create a new HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, nil, err
	}

	return resp, resp.Body, nil
}

func (s *BrowserlessFetcher) HasMetadata() bool {
	return false
}

func (s *BrowserlessFetcher) Metadata() metadata.Metadata {
	return metadata.Metadata{}
}
