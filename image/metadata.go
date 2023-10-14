package image

import (
	"fmt"
	_ "image/gif"  // This is required to initialize the GIF decoder
	_ "image/jpeg" // This is required to initialize the JPEG decoder
	_ "image/png"  // This is required to initialize the PNG decoder
	"io"
	"net/http"

	_ "github.com/biessek/golang-ico"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type Metadata struct {
	Width       int
	Height      int
	ContentSize int64
	ContentType string
}

var defaultMaxBytes = 524288 // 512 KB

// FetchMetadataFromHeader returns the image metadata from the given URL. However, it does not download the image.
// It will attempt to get the ContentType, Width, and Height from downloading at most maxBytes. If the metadata cannot
// be extracted within the first maxBytes, an error is returned. This is useful for validating images without
// downloading the entire image.
func FetchMetadataFromHeader(imgURL string, maxBytes int) (Metadata, error) {
	resp, err := http.Get(imgURL)
	if err != nil {
		return Metadata{}, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var imgData []byte
	buf := make([]byte, 8192) // 8 KB chunks
	totalRead := 0

	for totalRead < maxBytes {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return Metadata{}, err
		}
		if n == 0 {
			break
		}

		imgData = append(imgData, buf[:n]...)
		totalRead += n

		img, err := ImageFromBytes(imgData)

		if err == nil {
			//fmt.Printf("Found in first %d bytes\n", totalRead)
			return img.Metadata, nil
		}
	}

	return Metadata{}, fmt.Errorf("could not extract metadata within the first %d KB", maxBytes/1024)
}
