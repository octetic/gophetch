package media_test

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"
)

var (
	imgData = map[string][]byte{
		"mark.bmp":  {},
		"mark.gif":  {},
		"mark.jpg":  {},
		"mark.png":  {},
		"mark.tif":  {},
		"mark.webp": {},
		"mark.ico":  {},
	}
)

func loadImageData() error {
	for k := range imgData {
		data, err := os.ReadFile("../testdata/" + k)
		if err != nil {
			return err
		}
		imgData[k] = data
	}
	return nil
}

func setup() {
	err := loadImageData()
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()         // Setup before running tests
	code := m.Run() // Run tests
	// You can add teardown here if needed
	os.Exit(code)
}

// ReadAndEncodeImage reads an image file and returns its Base64 encoding
func ReadAndEncodeImage(filePath string) (string, error) {
	imageBytes, err := os.ReadFile("../testdata/" + filePath)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(imageBytes)

	// Get the extension after the last dot
	extension := strings.Split(filePath, ".")[1]

	switch extension {
	case "bmp":
		return "data:image/bmp;base64," + encoded, nil
	case "gif":
		return "data:image/gif;base64," + encoded, nil
	case "jpg", "jpeg":
		return "data:image/jpeg;base64," + encoded, nil
	case "png":
		return "data:image/png;base64," + encoded, nil
	case "tif", "tiff":
		return "data:image/tiff;base64," + encoded, nil
	case "webp":
		return "data:image/webp;base64," + encoded, nil
	case "ico":
		return "data:image/x-icon;base64," + encoded, nil
	default:
		return "", nil
	}
}
