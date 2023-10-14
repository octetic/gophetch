package gophetch

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"

	"github.com/minsoft-io/gophetch/image"
)

type InlineStrategy string

const (
	// StrategyInline stores images as base64 inline
	StrategyInline InlineStrategy = "inline"
	// StrategyHybrid stores images as base64 inline if they are less than maxContentSize or
	// smaller than maxWidth x maxHeight, otherwise it uploads them to cloud storage
	StrategyHybrid InlineStrategy = "hybrid"
	// StrategyUpload stores images in cloud storage
	StrategyUpload InlineStrategy = "upload"
)

// ImageFetcher is an interface for fetching images
type ImageFetcher interface {
	NewImageFromURL(url string) (*image.Image, error)
}

// RealImageFetcher uses the actual implementation
type RealImageFetcher struct{}

func (r *RealImageFetcher) NewImageFromURL(url string) (*image.Image, error) {
	return image.NewImageFromURL(url)
}

type UploadFunc func(*image.Image) (string, error)

type ShouldInlineFunc func(*image.Image) bool

// ImageInliner is responsible for fetching and replacing images in HTML documents.
type ImageInliner struct {
	ShouldInline   ShouldInlineFunc
	fetcher        ImageFetcher
	uploadFunc     UploadFunc
	strategy       InlineStrategy
	maxContentSize int64
	maxWidth       int
	maxHeight      int
}

// ImageInlinerOptions are options for creating a new ImageInliner.
type ImageInlinerOptions struct {
	// ShouldInlineFunc is the function to use for determining whether an image should be inlined. Default is to inline
	// if image size is less than 100KB or if dimensions are smaller than 800x600 (based on the maxContentSize, maxWidth,
	// and maxHeight options).
	ShouldInlineFunc ShouldInlineFunc
	// Fetcher is the ImageFetcher to use for fetching images.
	Fetcher ImageFetcher
	// UploadFunc is the function to use for uploading images to cloud storage.
	UploadFunc UploadFunc
	// Strategy is the storage strategy to use.
	Strategy InlineStrategy
	// MaxContentSize is the maximum size in bytes for images to be processed in a hybrid strategy. Default is 100KB.
	MaxContentSize int64
	// MaxWidth is the maximum width in pixels for images to be processed in a hybrid strategy. Default is 800.
	MaxWidth int
	// MaxHeight is the maximum height in pixels for images to be processed in a hybrid strategy. Default is 600.
	MaxHeight int
}

// NewImageInliner creates a new ImageInliner with the given fetcher, upload function, and storage strategy.
func NewImageInliner(opts ImageInlinerOptions) *ImageInliner {
	return &ImageInliner{
		ShouldInline: func() ShouldInlineFunc {
			if opts.ShouldInlineFunc == nil {
				return func(img *image.Image) bool {
					// Inline if image size is less than 100KB
					if img.ContentSize < 100*1024 {
						return true
					}
					// Inline if dimensions are smaller than 800x600
					if img.Width < 800 && img.Height < 600 {
						return true
					}
					return false
				}
			}
			return opts.ShouldInlineFunc
		}(),
		fetcher:    opts.Fetcher,
		uploadFunc: opts.UploadFunc,
		strategy:   opts.Strategy,
		maxContentSize: func() int64 {
			if opts.MaxContentSize == 0 {
				return 100 * 1024
			}
			return opts.MaxContentSize
		}(),
		maxWidth: func() int {
			if opts.MaxWidth == 0 {
				return 800
			}
			return opts.MaxWidth
		}(),
		maxHeight: func() int {
			if opts.MaxHeight == 0 {
				return 600
			}
			return opts.MaxHeight
		}(),
	}
}

// InlineImages replaces image URLs with either base64 inline versions or cloud URLs based on the set strategy.
func (inliner *ImageInliner) InlineImages(readableHTML string) (string, error) {
	var wg sync.WaitGroup

	doc, err := html.Parse(strings.NewReader(readableHTML))
	if err != nil {
		return "", err
	}

	// Find all image nodes
	imgSelector := cascadia.MustCompile("img, video")
	imgVideoNodes := imgSelector.MatchAll(doc)

	for _, node := range imgVideoNodes {
		wg.Add(1)
		go func(node *html.Node) {
			defer wg.Done()
			for i := range node.Attr {
				attr := &node.Attr[i]
				if attr.Key == "src" || attr.Key == "srcset" || attr.Key == "poster" {
					// Determine storage strategy based on file size, type, etc.
					switch inliner.strategy {
					case StrategyInline:
						// StrategyInline as base64
						attr.Val = inliner.fetchAndInline(attr)
					case StrategyUpload:
						// Upload to cloud storage and replace URL
						attr.Val = inliner.uploadAndReplaceAttr(attr)
					case StrategyHybrid:
						// Hybrid strategy
						attr.Val = inliner.processHybrid(attr)
					}
				}
			}
		}(node)
	}

	// Wait for all go routines to finish
	wg.Wait()

	// Convert the modified doc back to HTML string
	var b strings.Builder
	if err := html.Render(&b, doc); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (inliner *ImageInliner) fetchAndInline(attr *html.Attribute) string {
	urls, descriptors := inliner.parseSrcAndSrcset(attr)
	var newURLs []string

	for i, url := range urls {
		img, err := inliner.fetcher.NewImageFromURL(url)
		if err != nil {
			log.Printf("Failed to download image: %v", err)
			continue
		}
		imgBase64 := base64.StdEncoding.EncodeToString(img.Bytes)
		newURL := fmt.Sprintf("data:%s;base64,%s", img.ContentType, imgBase64)
		if attr.Key == "srcset" && descriptors[i] != "" {
			newURL += " " + descriptors[i]
		}
		newURLs = append(newURLs, newURL)
	}

	return strings.Join(newURLs, ", ")
}

func (inliner *ImageInliner) processHybrid(attr *html.Attribute) string {
	urls, descriptors := inliner.parseSrcAndSrcset(attr)
	var newURLs []string

	for i, url := range urls {
		img, err := inliner.fetcher.NewImageFromURL(url)
		if err != nil {
			newURLs = append(newURLs, url)
			continue
		}

		var newURL string
		if inliner.ShouldInline(img) {
			newURL = fmt.Sprintf("data:%s;base64,%s", img.ContentType, base64.StdEncoding.EncodeToString(img.Bytes))
		} else {
			newURL, err = inliner.uploadFunc(img)
			if err != nil {
				// If upload fails, use the original URL
				newURL = url
			}
		}

		if attr.Key == "srcset" && descriptors[i] != "" {
			newURL += " " + descriptors[i]
		}
		newURLs = append(newURLs, newURL)
	}

	return strings.Join(newURLs, ", ")
}

func (inliner *ImageInliner) uploadAndReplaceAttr(attr *html.Attribute) string {
	urls, descriptors := inliner.parseSrcAndSrcset(attr)
	var newURLs []string

	for i, url := range urls {
		img, err := inliner.fetcher.NewImageFromURL(url)
		if err != nil {
			log.Printf("Failed to download image: %v", err)
			continue
		}

		newURL, err := inliner.uploadFunc(img)
		if err != nil {
			// If upload fails, use the original URL
			newURL = url
		}

		if attr.Key == "srcset" && descriptors[i] != "" {
			newURL += " " + descriptors[i]
		}
		newURLs = append(newURLs, newURL)
	}

	return strings.Join(newURLs, ", ")
}

func (inliner *ImageInliner) parseSrcAndSrcset(attr *html.Attribute) ([]string, []string) {
	var urls []string
	var descriptors []string

	if attr.Key == "src" {
		urls = []string{attr.Val}
	} else {
		entries := strings.Split(attr.Val, ",")
		for _, entry := range entries {
			parts := strings.Fields(strings.TrimSpace(entry))
			if len(parts) > 0 {
				urls = append(urls, parts[0])
				if len(parts) > 1 {
					descriptors = append(descriptors, parts[1])
				} else {
					descriptors = append(descriptors, "")
				}
			}
		}
	}

	return urls, descriptors
}
