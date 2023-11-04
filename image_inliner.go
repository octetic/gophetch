package gophetch

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
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
	// than the maxInlinedSize, maxWidth, and maxHeight options, and uploaded to cloud storage otherwise.
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
	NewImageFromURL(url string, maxSize int) (*image.Image, error)
}

// RealImageFetcher uses the actual implementation
type RealImageFetcher struct{}

// NewImageFromURL fetches an image from the given URL.
func (r *RealImageFetcher) NewImageFromURL(url string, maxSize int) (*image.Image, error) {
	return image.NewImageFromURL(url, maxSize)
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
	maxInlinedSize int
	maxWidth       int
	maxHeight      int
}

// ImageInlinerOptions are options for creating a new ImageInliner.
type ImageInlinerOptions struct {
	// ShouldInlineFunc is the function to use for determining whether an image should be inlined. Default is to inline
	// if image size is less than 100KB or if dimensions are smaller than 800x600 (based on the maxInlinedSize, maxWidth,
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
	// MaxContentSize is the maximum size in bytes for images to be processed and uploaded. Default is 10MB.
	MaxContentSize int64
	// MaxInlinedSize is the maximum size in bytes for images to be processed in a hybrid strategy. Default is 100KB.
	MaxInlinedSize int
	// MaxWidth is the maximum width in pixels for images to be processed in a hybrid strategy. Default is 800.
	MaxWidth int
	// MaxHeight is the maximum height in pixels for images to be processed in a hybrid strategy. Default is 600.
	MaxHeight int
}

// NewImageInliner creates a new ImageInliner with the given fetcher, upload function, and storage strategy.
func NewImageInliner(opts ImageInlinerOptions) *ImageInliner {
	inlinedSize := 100 * 1024
	if opts.MaxInlinedSize > 0 {
		inlinedSize = opts.MaxInlinedSize
	}
	maxWidth := 800
	if opts.MaxWidth > 0 {
		maxWidth = opts.MaxWidth
	}
	maxHeight := 600
	if opts.MaxHeight > 0 {
		maxHeight = opts.MaxHeight
	}
	return &ImageInliner{
		ShouldInline: func() ShouldInlineFunc {
			if opts.ShouldInlineFunc == nil {
				return func(img *image.Image) bool {
					// Inline if image size is less than 100KB
					if img.ContentSize < int64(inlinedSize) {
						return true
					}
					// Inline if dimensions are smaller than 800x600
					if img.Width < maxWidth && img.Height < maxHeight {
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
				return 10 * 1024 * 1024
			}
			return opts.MaxContentSize
		}(),
		maxInlinedSize: inlinedSize,
		maxWidth:       maxWidth,
		maxHeight:      maxHeight,
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
		img, err := inliner.fetcher.NewImageFromURL(url, int(inliner.maxContentSize))
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
		img, err := inliner.fetcher.NewImageFromURL(url, int(inliner.maxContentSize))
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
		img, err := inliner.fetcher.NewImageFromURL(url, int(inliner.maxContentSize))
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
		urls, descriptors = ExtractSrcset(attr.Val)
	}

	return urls, descriptors
}

// ExtractSrcset attempts to match all srcset URLs including their descriptors,
// accounting for commas within the URLs.
func ExtractSrcset(srcset string) ([]string, []string) {
	// This regex captures the URL and the descriptor as separate groups
	// - (https://[^\s]+) is a capturing group that matches a URL starting with https:// and continues without any space or comma.
	// - \s+\d+(?:\.\d+)?[wx] matches one or more spaces followed by one or more digits (with optional decimal) and then 'w' or 'x', which represent the descriptors.
	// - (,|\s|$) ensures that this pattern is followed by a comma, whitespace, or the end of the string, meaning it's the end of a URL/descriptor segment.
	re := regexp.MustCompile(`(https://\S+)((\s+\d+(?:\.\d+)?[wx])+)(?:,|$)`)

	// Find all matches for the pattern.
	matches := re.FindAllStringSubmatch(strings.TrimSpace(srcset), -1)

	var urls []string
	var descriptors []string

	for _, match := range matches {
		if len(match) > 2 {
			urls = append(urls, match[1])                                  // The URL is in the first capture group
			descriptors = append(descriptors, strings.TrimSpace(match[2])) // The descriptor is in the second capture group
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

		// Find the first image that starts with "https://"
		for i, descriptor := range descriptors {
			if strings.HasPrefix(urls[i], "https://") || strings.HasPrefix(urls[i], "http://") {
				smallestIndex := i
				selectedURLs = append(selectedURLs, urls[smallestIndex])
				selectedDescriptors = append(selectedDescriptors, descriptor)
				break
			}
		}

	case SrcsetLargestImage:
		// Assuming the last descriptor usually points to the largest image. This is not always true, but it's a
		// reasonable assumption.

		// Find the first image that starts with "https://", starting from the end
		for i := len(urls) - 1; i >= 0; i-- {
			if strings.HasPrefix(urls[i], "https://") || strings.HasPrefix(urls[i], "http://") {
				largestIndex := i
				selectedURLs = append(selectedURLs, urls[largestIndex])
				selectedDescriptors = append(selectedDescriptors, descriptors[largestIndex])
				break
			}
		}

	case SrcsetPreferredDescriptors:
		found := false
		for _, preferred := range preferredDescriptors {
			for i, descriptor := range descriptors {
				if descriptor == preferred {
					if strings.HasPrefix(urls[i], "https://") || strings.HasPrefix(urls[i], "http://") {
						selectedURLs = append(selectedURLs, urls[i])
						selectedDescriptors = append(selectedDescriptors, descriptor)
						found = true
						break
					}
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
