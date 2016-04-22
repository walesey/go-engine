package assets

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
)

func ImportImage(file string) image.Image {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}
	img = imaging.FlipV(img)
	return img
}
