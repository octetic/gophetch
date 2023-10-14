package gophetch_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stretchr/testify/mock"

	"github.com/minsoft-io/gophetch/image"
)

// MockImageFetcher returns a mock image
type MockImageFetcher struct{}

func (m *MockImageFetcher) NewImageFromURL(url string) (*image.Image, error) {
	filename := strings.TrimPrefix(url, "https://example.com/")
	imgBytes, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		return nil, err
	}

	contentType, err := image.ContentTypeByExtension("." + strings.ToLower(filename[strings.LastIndex(filename, ".")+1:]))
	if err != nil {
		return nil, err
	}

	return &image.Image{
		Bytes: imgBytes,
		Metadata: image.Metadata{
			ContentType: contentType,
		},
	}, nil
}

type MockImageFetcherHybrid struct {
	mock.Mock
}

func (m *MockImageFetcherHybrid) NewImageFromURL(url string) (*image.Image, error) {
	//args := m.Called(url)
	filename := strings.TrimPrefix(url, "https://example.com/")
	imgBytes, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		return nil, err
	}

	contentType, err := image.ContentTypeByExtension("." + strings.ToLower(filename[strings.LastIndex(filename, ".")+1:]))
	if err != nil {
		return nil, err
	}

	return &image.Image{
		Bytes: imgBytes,
		Metadata: image.Metadata{
			ContentType: contentType,
			Width:       100,
			Height:      100,
			ContentSize: 100 * 100,
		},
	}, nil
}

var imageBase64Map map[string]string
var testImages = []string{
	"mark.bmp",
	"mark.gif",
	"mark.ico",
	"mark.jpg",
	"mark.png",
	"mark.tif",
	"mark.webp",
}

func SetupFiles() {
	imageBase64Map = make(map[string]string)
	// Assume testImages is a slice containing the paths to the test images
	for _, path := range testImages {
		pathExt := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
		imgBytes, err := os.ReadFile("testdata/" + path)
		if err != nil {
			log.Fatalf("Failed to read test image: %v", err)
		}
		contentType, err := image.ContentTypeByExtension("." + pathExt)
		if err != nil {
			log.Fatalf("Failed to get content type: %v", err)
		}
		imgBase64 := base64.StdEncoding.EncodeToString(imgBytes)
		newURL := fmt.Sprintf("data:%s;base64,%s", contentType, imgBase64)
		imageBase64Map[path] = newURL
	}
}
