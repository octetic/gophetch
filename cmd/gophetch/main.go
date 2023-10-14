package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3"

	"github.com/minsoft-io/gophetch"
	"github.com/minsoft-io/gophetch/fetchers"
	"github.com/minsoft-io/gophetch/metadata"
)

type config struct {
	URL                string
	MicrolinkAPIKey    string
	ScrapingFishAPIKey string
	PrintHeaders       bool
	PrintMetadata      bool
	PrintHTML          bool
}

func main() {
	cfg := config{}

	fs := flag.NewFlagSet("rgserver", flag.ContinueOnError)
	fs.StringVar(&cfg.MicrolinkAPIKey, "microlink-api-key", "", "Microlink API key")
	fs.StringVar(&cfg.ScrapingFishAPIKey, "scrapingfish-api-key", "", "ScrapingFish API key")
	fs.StringVar(&cfg.URL, "url", "", "URL to fetch")
	fs.BoolVar(&cfg.PrintHeaders, "headers", false, "Print headers")
	fs.BoolVar(&cfg.PrintMetadata, "metadata", false, "Print metadata")
	fs.BoolVar(&cfg.PrintHTML, "html", false, "Print HTML")

	showVersion := fs.Bool("v", false, "display version and exit")

	// Parse flags first, then environment variables
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("GOPHETCH"))
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}

		panic(err)
	}

	if *showVersion {
		fmt.Printf("\nGoPhetch\n")
		fmt.Printf("Version: %s\n", "0.0.2")
		os.Exit(0)
	}

	if cfg.URL == "" {
		fmt.Println("URL is required")
		os.Exit(1)
	}

	standardFetcher := &fetchers.StandardHTTPFetcher{}
	microLinkFetcher := &fetchers.MicrolinkFetcher{
		AdBlock:   true,
		APIKey:    cfg.MicrolinkAPIKey,
		Prerender: true,
	}
	//proxyFetcher := &fetchers.ProxyHTTPFetcher{ProxyURL: "http://your-proxy-url"}
	scrapingfishFetcher := &fetchers.ScrapingfishFetcher{
		APIKey: cfg.ScrapingFishAPIKey,
	}

	allow := false
	var htmlFetchers []fetchers.HTMLFetcher

	if allow {
		if cfg.ScrapingFishAPIKey != "" {
			htmlFetchers = append(htmlFetchers, scrapingfishFetcher)
		}

		if cfg.MicrolinkAPIKey != "" {
			htmlFetchers = append(htmlFetchers, microLinkFetcher)
		}
	}

	// Add the standard fetcher last
	htmlFetchers = append(htmlFetchers, standardFetcher)

	g := gophetch.New(htmlFetchers...)
	data, err := g.FetchAndParse(cfg.URL)
	if err != nil {
		panic(err)
	}

	printStatusAndMime(data.StatusCode, cfg.URL)

	if cfg.PrintHeaders {
		printHeaders(data.Headers)
	}

	if cfg.PrintMetadata {
		printMetadata(data.Metadata)
	}

	if cfg.PrintHTML {
		printHTML(data.Metadata.HTML)
	}
}

func printStatusAndMime(statusCode int, url string) {
	fmt.Println("Parsed URL: ", url)
	fmt.Printf("Status: %d\n", statusCode)
	fmt.Printf("MIME type: %s\n", "text/html")
}

func printHeaders(headers map[string][]string) {
	fmt.Println("HEADERS: ")
	for key, value := range headers {
		fmt.Printf("%s: %s\n", key, value)
	}
}

func printMetadata(metadata metadata.Metadata) {
	fmt.Println("METADATA: ")
	fmt.Printf("Audio: %v\n", metadata.Audio)
	fmt.Printf("Author: %s\n", metadata.Author)
	fmt.Printf("CanonicalURL: %s\n", metadata.CanonicalURL)
	fmt.Printf("Date: %s\n", metadata.Date)
	fmt.Printf("Description: %s\n", metadata.Description)
	//fmt.Printf("HTML: %s\n", metadata.HTML)
	fmt.Printf("FaviconURL: %s\n", metadata.FaviconURL)
	fmt.Printf("Feed: %v\n", metadata.FeedURLs)
	fmt.Printf("LeadImageURL: %s\n", metadata.LeadImageURL)
	fmt.Printf("Image: %v\n", metadata.Image)
	fmt.Printf("Lang: %s\n", metadata.Lang)
	fmt.Printf("Logo: %v\n", metadata.Logo)
	fmt.Printf("Meta: %v\n", metadata.Meta)
	fmt.Printf("Publisher: %s\n", metadata.Publisher)
	fmt.Printf("Title: %s\n", metadata.Title)
	fmt.Printf("URL: %s\n", metadata.URL)
	fmt.Printf("Video: %v\n", metadata.Video)
	fmt.Printf("ReadableExcerpt: %s\n", metadata.ReadableExcerpt)
	fmt.Printf("ReadableHTML: %s\n", metadata.ReadableHTML)
	fmt.Printf("ReadableText: %s\n", metadata.ReadableText)
	fmt.Printf("ReadableImage: %s\n", metadata.ReadableImage)
	fmt.Printf("ReadableLang: %s\n", metadata.ReadableLang)
	fmt.Printf("ReadableTitle: %s\n", metadata.ReadableTitle)
	fmt.Printf("ReadableByline: %s\n", metadata.ReadableByline)
	fmt.Printf("ReadableSiteName: %s\n", metadata.ReadableSiteName)
	fmt.Printf("SiteName: %s\n", metadata.SiteName)

	for key, value := range metadata.Dynamic {
		fmt.Printf("%s: %s\n", key, value)
	}
}

func printHTML(html string) {
	fmt.Println("HTML: ")
	fmt.Printf("%s\n", html)
}
