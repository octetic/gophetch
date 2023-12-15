package media

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"
	"time"
)

// DefaultMaxMediaSize represents the default maximum number of bytes we are willing to download for an image. (10 MB)
const DefaultMaxMediaSize = 30 * 1024 * 1024

type Type int

const (
	ImageType Type = iota
	VideoType
	AudioType
	VectorImageType // For SVGs, ICOs, etc.
)

type Media struct {
	Metadata
	Bytes       []byte
	Cache       Cache
	ContentSize int64
	Format      string
	Image       image.Image // To hold the image.Image object if it is an image
	URL         string
	Extension   string
	MediaType   Type
}

// ShouldCacheImage takes a http.Header and returns whether the image should be cached
func (i *Media) ShouldCacheImage() bool {
	if i.Cache.NoStore {
		return false
	}

	if i.Cache.MaxAge > 0 {
		return true
	}

	if !i.Cache.Expires.IsZero() {
		return true
	}

	return false
}

// ShouldRevalidateImage takes a http.Header and returns whether the image should be revalidated
func (i *Media) ShouldRevalidateImage() bool {
	if !i.Cache.Available {
		return true
	}

	if i.Cache.NoCache {
		return true
	}

	if i.Cache.MustRevalidate {
		return true
	}

	if i.Cache.MaxAge == 0 {
		return true
	}

	if !i.Cache.Expires.IsZero() && i.Cache.Expires.Before(time.Now()) {
		return true
	}

	return false
}

// ShouldRefreshImage takes a http.Header and returns whether the image should be refreshed
func (i *Media) ShouldRefreshImage() bool {

	// Should refresh if the max age is > 0 and the expires time is in the past
	if i.Cache.MaxAge > 0 && !i.Cache.Expires.IsZero() && !i.Cache.Expires.Before(time.Now()) {
		return true
	}

	if i.Cache.MaxAge >= 0 {
		return false
	}

	// Should not refresh if the max age is not present and the expires time is still in the future.
	// This is a special case where the max age is not set, but the expires time is set.
	if !i.Cache.Expires.IsZero() && i.Cache.Expires.After(time.Now()) {
		return false
	}

	return true
}

// GenerateUniqueFilename generates a unique filename based on the Image properties
func (i *Media) GenerateUniqueFilename() string {
	var hashString string

	// If URL is empty, generate a random hash
	if i.URL == "" {
		randomData := make([]byte, 16) // using 16 bytes for the random part
		_, err := rand.Read(randomData)
		if err != nil {
			// handle error
			return ""
		}
		hashString = hex.EncodeToString(randomData)
	} else {
		// Hash the URL using SHA-1
		h := sha1.New()
		h.Write([]byte(i.URL))
		hashed := h.Sum(nil)
		hashString = fmt.Sprintf("%x", hashed)
	}

	// Append extension if available
	if i.Extension != "" {
		return fmt.Sprintf("%s.%s", hashString, i.Extension)
	}

	return hashString
}

// GenerateCacheKey generates a cache key for the given URL
func GenerateCacheKey(url string) string {
	// Create a new hash.Hash object for SHA-256
	hasher := sha256.New()

	// Write the URL into the hasher
	hasher.Write([]byte(url)) // Ignoring error for simplicity; Write on sha256 never returns an error

	// Compute the SHA-256 checksum of the URL
	sum := hasher.Sum(nil)

	// Return the hexadecimal encoding of the checksum
	return hex.EncodeToString(sum)
}

// NewImageFromBytes returns the image and metadata from the given bytes. It will attempt to get the ContentType, Width,
// Height, Format, and ContentSize from the given bytes. If the metadata cannot be extracted, an error is returned.
func NewImageFromBytes(data []byte) (*Media, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	contentType := http.DetectContentType(data)
	extension, err := ExtensionByContentType(contentType)
	if err != nil {
		extension = ""
	}

	return &Media{
		Bytes:       data,
		ContentSize: int64(len(data)),
		Metadata: Metadata{
			Width:       width,
			Height:      height,
			ContentSize: int64(len(data)),
			ContentType: contentType,
		},
		Format:    format,
		Image:     img,
		Extension: extension,
		MediaType: ImageType,
	}, nil
}

// NewMediaFromURL will download the media from the given URL and return the media and metadata, but only if it is within the MaxMediaSize.
func NewMediaFromURL(mediaURL string, maxSize int) (*Media, error) {
	if maxSize <= 0 {
		maxSize = DefaultMaxMediaSize
	}

	resp, err := http.Get(mediaURL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	contentType := resp.Header.Get("Content-Type")

	media := &Media{
		URL: mediaURL,
	}

	switch {
	case strings.HasPrefix(contentType, "image/svg"):
		media.MediaType = VectorImageType
		return NewVectorImageFromHTTPResponse(resp, mediaURL, maxSize)

	case strings.HasPrefix(contentType, "image/"):
		return NewImageFromHTTPResponse(resp, mediaURL, maxSize)

	case strings.HasPrefix(contentType, "audio/"), strings.HasPrefix(contentType, "video/"):
		return NewAudioVideoFromHTTPResponse(resp, mediaURL, maxSize)

	default:
		return nil, fmt.Errorf("unsupported media type: %s", contentType)
	}

}

// NewImageFromURL will download the image from the given URL and return the image and metadata, but only if it is within the MaxImageSize.
func NewImageFromURL(imgURL string, maxSize int) (*Media, error) {
	resp, err := http.Get(imgURL)
	if err != nil {
		return nil, err
	}
	return NewImageFromHTTPResponse(resp, imgURL, maxSize)
}

// NewImageFromHTTPResponse will download the image from the given URL and return the image and metadata, but only if it is within the MaxImageSize.
func NewImageFromHTTPResponse(resp *http.Response, imgURL string, maxSize int) (*Media, error) {
	if maxSize <= 0 {
		maxSize = DefaultMaxMediaSize
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var imgData bytes.Buffer
	// Create a buffer to read into and count the bytes read.
	buf := make([]byte, 1024) // 1KB chunks
	totalRead := 0

	for {
		// Read a chunk.
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		totalRead += n

		// Check if we have read too much.
		if totalRead > maxSize {
			return nil, fmt.Errorf("image exceeds max size of %d bytes", maxSize)
		}

		// Write the bytes to the buffer.
		_, writeErr := imgData.Write(buf[:n])
		if writeErr != nil {
			return nil, writeErr
		}

		// Check for end of file.
		if err == io.EOF {
			break
		}
	}

	img, imgErr := NewImageFromBytes(imgData.Bytes())
	if imgErr != nil {
		return nil, imgErr
	}

	img.Cache = ParseCacheHeader(resp.Header)
	img.URL = imgURL

	return img, nil
}

// NewVectorImageFromHTTPResponse retrieves an SVG or ICO image from the given URL and returns it as raw bytes
func NewVectorImageFromHTTPResponse(resp *http.Response, imgURL string, maxSize int) (*Media, error) {
	if maxSize <= 0 {
		maxSize = DefaultMaxMediaSize
	}

	var imgData bytes.Buffer
	// Create a buffer to read into and count the bytes read.
	buf := make([]byte, 1024) // 1KB chunks
	totalRead := 0

	for {
		// Read a chunk.
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		totalRead += n

		// Check if we have read too much.
		if totalRead > maxSize {
			return nil, fmt.Errorf("vector image exceeds max size of %d bytes", maxSize)
		}

		// Write the bytes to the buffer.
		_, writeErr := imgData.Write(buf[:n])
		if writeErr != nil {
			return nil, writeErr
		}

		// Check for end of file.
		if err == io.EOF {
			break
		}
	}

	contentType := resp.Header.Get("Content-Type")
	extension, err := ExtensionByContentType(contentType)
	if err != nil {
		extension = ""
	}

	media := &Media{
		Bytes:       imgData.Bytes(),
		ContentSize: int64(totalRead),
		Metadata: Metadata{
			Width:       0,
			Height:      0,
			ContentSize: int64(totalRead),
			ContentType: contentType,
		},
		Format:    "",
		Extension: extension,
		URL:       imgURL,
		Cache:     ParseCacheHeader(resp.Header),
		MediaType: VectorImageType,
	}

	return media, nil
}

// NewAudioVideoFromHTTPResponse retrieves an audio or video file from the given URL and returns it as raw bytes
func NewAudioVideoFromHTTPResponse(resp *http.Response, mediaURL string, maxSize int) (*Media, error) {
	if maxSize <= 0 {
		maxSize = DefaultMaxMediaSize
	}

	var mediaData bytes.Buffer
	buf := make([]byte, 1024) // 1KB chunks
	totalRead := 0

	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		totalRead += n

		if totalRead > maxSize {
			return nil, fmt.Errorf("media file exceeds max size of %d bytes", maxSize)
		}

		_, writeErr := mediaData.Write(buf[:n])
		if writeErr != nil {
			return nil, writeErr
		}

		if err == io.EOF {
			break
		}
	}

	contentType := resp.Header.Get("Content-Type")
	extension, err := ExtensionByContentType(contentType)
	if err != nil {
		extension = ""
	}

	mediaType := AudioType
	if strings.HasPrefix(contentType, "video/") {
		mediaType = VideoType
	}

	media := &Media{
		Bytes:       mediaData.Bytes(),
		ContentSize: int64(totalRead),
		Metadata: Metadata{
			Width:       0,
			Height:      0,
			ContentSize: int64(totalRead),
			ContentType: contentType,
		},
		Format:    "",
		Extension: extension,
		URL:       mediaURL,
		Cache:     ParseCacheHeader(resp.Header),
		MediaType: mediaType,
	}

	return media, nil
}

// NewImageFromDataURI will parse the data URI and return the image and metadata. It will attempt to get the
// ContentType, Width, Height, Format, and ContentSize from the data URI as well.
func NewImageFromDataURI(dataURI string) (*Media, error) {
	imgData, err := DataURIToBytes(dataURI)
	if err != nil {
		return nil, err
	}

	img, err := NewImageFromBytes(imgData)
	if err != nil {
		return nil, err
	}

	img.URL = dataURI
	return img, nil
}

// DataURIToBytes takes a Data URI and returns the byte data it contains.
func DataURIToBytes(dataURI string) ([]byte, error) {
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) < 2 {
		return nil, errors.New("invalid Data URI")
	}

	// The second part contains the actual data
	data := parts[1]

	// Decode the Base64 data
	return base64.StdEncoding.DecodeString(data)
}
