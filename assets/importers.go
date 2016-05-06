package assets

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/disintegration/imaging"
)

func ImportImage(file string) (image.Image, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening image file: %v\n", err)
		return nil, err
	}
	return DecodeImage(imgFile)
}

func DecodeImage(data io.Reader) (image.Image, error) {
	img, _, err := image.Decode(data)
	if err != nil {
		fmt.Printf("Error decoding geometry file: %v\n", err)
		return nil, err
	}
	img = imaging.FlipV(img)
	return img, nil
}
