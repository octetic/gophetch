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

// InlineStrategy represents the different strategies for inlining images.
type InlineStrategy int

const (
	// InlineAll indicates that all images should be inlined.
	InlineAll InlineStrategy = iota

	// InlineNone indicates that no images should be inlined, and all images should be uploaded to
	// cloud storage using the upload function.
	InlineNone

	// InlineHybrid indicates a hybrid approach to inlining images where images are inlined if they are smaller
	// than the maxContentSize, maxWidth, and maxHeight options, and uploaded to cloud storage otherwise.
	InlineHybrid
)

// SrcsetStrategy represents the different strategies for handling srcset attributes.
type SrcsetStrategy int

const (
	// SrcsetSmallestImage selects the smallest image in the srcset.
	SrcsetSmallestImage SrcsetStrategy = iota

	// SrcsetLargestImage selects the largest image in the srcset.
	SrcsetLargestImage

	// SrcsetPreferredDescriptors selects an image based on the preferred descriptors.
	// Currently only looks for 2x, 1.5x, and 1x, in that order.
	SrcsetPreferredDescriptors

	// SrcsetAllImages includes all images in the srcset.
	SrcsetAllImages
)

// ImageFetcher is an interface for fetching images
type ImageFetcher interface {
	NewImageFromURL(url string) (*image.Image, error)
}

// RealImageFetcher uses the actual implementation
type RealImageFetcher struct{}

// NewImageFromURL fetches an image from the given URL.
func (r *RealImageFetcher) NewImageFromURL(url string) (*image.Image, error) {
	return image.NewImageFromURL(url)
}

// UploadFunc is the function signature to use for uploading images to cloud storage.
type UploadFunc func(*image.Image) (string, error)

// ShouldInlineFunc is the function signature use for determining whether an image should be inlined.
type ShouldInlineFunc func(*image.Image) bool

// ImageInliner is responsible for fetching and replacing images in HTML documents.
type ImageInliner struct {
	ShouldInline   ShouldInlineFunc
	fetcher        ImageFetcher
	uploadFunc     UploadFunc
	inlineStrategy InlineStrategy
	srcsetStrategy SrcsetStrategy
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
	// InlineStrategy is the storage strategy to use. Default is InlineAll.
	InlineStrategy InlineStrategy
	// SrcsetStrategy is the strategy to use for handling srcset attributes. Default is SrcsetSmallestImage.
	SrcsetStrategy SrcsetStrategy
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
		fetcher: func() ImageFetcher {
			if opts.Fetcher == nil {
				return &RealImageFetcher{}
			}
			return opts.Fetcher
		}(),
		uploadFunc: func() UploadFunc {
			if opts.UploadFunc == nil {
				return func(img *image.Image) (string, error) {
					return "", fmt.Errorf("no upload function provided")
				}
			}
			return opts.UploadFunc
		}(),
		inlineStrategy: opts.InlineStrategy,
		srcsetStrategy: opts.SrcsetStrategy,
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
	imgSelector, err := cascadia.Compile("img, video, picture > source")
	if err != nil {
		return "", fmt.Errorf("failed to compile selector: %v", err)
	}
	imgVideoNodes := imgSelector.MatchAll(doc)

	for _, node := range imgVideoNodes {
		wg.Add(1)
		go func(node *html.Node) {
			defer wg.Done()
			for i := range node.Attr {
				attr := &node.Attr[i]
				if attr.Key == "src" || attr.Key == "srcset" || attr.Key == "poster" {
					// If the image is already inlined, skip it
					if strings.HasPrefix(attr.Val, "data:") {
						continue
					}

					// Determine storage strategy based on file size, type, etc.
					switch inliner.inlineStrategy {
					case InlineAll:
						// StrategyInline as base64
						attr.Val = inliner.fetchAndInline(attr)
					case InlineNone:
						// Upload to cloud storage and replace URL
						attr.Val = inliner.uploadAndReplaceAttr(attr)
					case InlineHybrid:
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
	if attr.Key == "srcset" {
		urls, descriptors = inliner.selectSrcsetURL(urls, descriptors)
	}

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
	if attr.Key == "srcset" {
		urls, descriptors = inliner.selectSrcsetURL(urls, descriptors)
	}

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
	if attr.Key == "srcset" {
		urls, descriptors = inliner.selectSrcsetURL(urls, descriptors)
	}

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

// selectSrcsetURL selects URLs based on the chosen SrcsetStrategy.
// It returns a slice of selected URLs and their corresponding descriptors.
func (inliner *ImageInliner) selectSrcsetURL(urls []string, descriptors []string) ([]string, []string) {
	var selectedURLs []string
	var selectedDescriptors []string

	preferredDescriptors := []string{"2x", "1.5x", "1x"} // Add more if needed

	switch inliner.srcsetStrategy {
	case SrcsetSmallestImage:
		// Assuming the first descriptor usually points to the smallest image. This is not always true, but it's a
		// reasonable assumption.
		smallestIndex := 0
		selectedURLs = append(selectedURLs, urls[smallestIndex])
		selectedDescriptors = append(selectedDescriptors, descriptors[smallestIndex])

	case SrcsetLargestImage:
		// Assuming the last descriptor usually points to the largest image. This is not always true, but it's a
		// reasonable assumption.
		largestIndex := len(urls) - 1
		selectedURLs = append(selectedURLs, urls[largestIndex])
		selectedDescriptors = append(selectedDescriptors, descriptors[largestIndex])

	case SrcsetPreferredDescriptors:
		found := false
		for _, preferred := range preferredDescriptors {
			for i, descriptor := range descriptors {
				if descriptor == preferred {
					selectedURLs = append(selectedURLs, urls[i])
					selectedDescriptors = append(selectedDescriptors, descriptor)
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			// Default to the last if no preferred descriptor is found
			lastIndex := len(urls) - 1
			selectedURLs = append(selectedURLs, urls[lastIndex])
			selectedDescriptors = append(selectedDescriptors, descriptors[lastIndex])
		}

	case SrcsetAllImages:
		// Use all URLs and descriptors (default behavior)
		selectedURLs = urls
		selectedDescriptors = descriptors
	}

	return selectedURLs, selectedDescriptors
}
