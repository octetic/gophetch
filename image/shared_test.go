package image_test

import (
	"os"
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
