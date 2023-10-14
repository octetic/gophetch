package image

import (
	"bytes"
	"image"
	"io"
	"net/http"
	"time"
)

type Image struct {
	Metadata
	Bytes       []byte
	Cache       Cache
	ContentSize int64
	Format      string
	Image       image.Image
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

// ImageFromBytes returns the image and metadata from the given bytes. It will attempt to get the ContentType, Width,
// Height, Format, and ContentSize from the given bytes. If the metadata cannot be extracted, an error is returned.
func ImageFromBytes(data []byte) (*Image, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	contentType := http.DetectContentType(data)

	return &Image{
		Bytes:       data,
		ContentSize: int64(len(data)),
		Metadata: Metadata{
			Width:       width,
			Height:      height,
			ContentSize: int64(len(data)),
			ContentType: contentType,
		},
		Format: format,
		Image:  img,
	}, nil
}

// ImageFromURL will download the image from the given URL and return the image and metadata. It will attempt to get the
// ContentType, Width, Height, Format, and ContentSize from the downloaded bytes as well.
func ImageFromURL(imgURL string) (*Image, error) {
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

	img, err := ImageFromBytes(imgData)
	if err != nil {
		return nil, err
	}

	img.Cache = ParseCacheHeader(resp.Header)
	return img, nil
}
