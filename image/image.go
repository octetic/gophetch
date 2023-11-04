package image

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
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

// DefaultMaxImageSize represents the default maximum number of bytes we are willing to download for an image. (10 MB)
const DefaultMaxImageSize = 10 * 1024 * 1024

type Image struct {
	Metadata
	Bytes       []byte
	Cache       Cache
	ContentSize int64
	Format      string
	Image       image.Image
	URL         string
	Extension   string
}

// ShouldCacheImage takes a http.Header and returns whether the image should be cached
func (i *Image) ShouldCacheImage() bool {
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
func (i *Image) ShouldRevalidateImage() bool {
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
func (i *Image) ShouldRefreshImage() bool {

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
func (i *Image) GenerateUniqueFilename() string {
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

// NewImageFromBytes returns the image and metadata from the given bytes. It will attempt to get the ContentType, Width,
// Height, Format, and ContentSize from the given bytes. If the metadata cannot be extracted, an error is returned.
func NewImageFromBytes(data []byte) (*Image, error) {
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

	return &Image{
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
	}, nil
}

// NewImageFromURLOLD will download the image from the given URL and return the image and metadata. It will attempt to get the
// ContentType, Width, Height, Format, and ContentSize from the downloaded bytes as well.
func NewImageFromURLOLD(imgURL string) (*Image, error) {
	resp, err := http.Get(imgURL)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	img, err := NewImageFromBytes(imgData)
	if err != nil {
		return nil, err
	}

	img.Cache = ParseCacheHeader(resp.Header)
	img.URL = imgURL
	return img, nil
}

// NewImageFromURL will download the image from the given URL and return the image and metadata, but only if it is within the MaxImageSize.
func NewImageFromURL(imgURL string, maxSize int) (*Image, error) {
	if maxSize <= 0 {
		maxSize = DefaultMaxImageSize
	}

	resp, err := http.Get(imgURL)
	if err != nil {
		return nil, err
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

// NewImageFromDataURI will parse the data URI and return the image and metadata. It will attempt to get the
// ContentType, Width, Height, Format, and ContentSize from the data URI as well.
func NewImageFromDataURI(dataURI string) (*Image, error) {
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
