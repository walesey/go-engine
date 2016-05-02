package assets

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/disintegration/imaging"
)

func ImportImage(file string) (image.Image, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening image file: %v\n", err)
		return nil, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Printf("Error decoding geometry file: %v\n", err)
		return nil, err
	}
	img = imaging.FlipV(img)
	return img, nil
}
