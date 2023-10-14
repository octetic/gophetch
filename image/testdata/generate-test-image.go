package main

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	//"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func main() {
	img := createImage()

	saveImage(img, "test_image.bmp", bmp.Encode)
	saveImage(img, "test_image.jpeg", encodeJPEG)
	saveImage(img, "test_image.gif", encodeGIF)
	saveImage(img, "test_image.png", png.Encode)
	saveImage(img, "test_image.tiff", encodeTIFF)
	//saveImage(img, "test_image.webp", encodeWEBP)
}

func createImage() image.Image {
	const width, height = 50, 50
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.RGBA{R: 255, A: 255}) // red
			} else {
				img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255}) // white
			}
		}
	}
	return img
}

func saveImage(img image.Image, filename string, encodeFunc func(w io.Writer, m image.Image) error) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	err = encodeFunc(file, img)
	if err != nil {
		panic(err)
	}
}

func encodeGIF(w io.Writer, m image.Image) error {
	return gif.Encode(w, m, &gif.Options{NumColors: 256})
}

func encodeJPEG(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, &jpeg.Options{Quality: 100})
}

func encodeTIFF(w io.Writer, m image.Image) error {
	return tiff.Encode(w, m, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
}

//func encodeWEBP(w io.Writer, m image.Image) error {
//	return webp.Encode(w, m, &webp.Options{Lossless: true})
//}
