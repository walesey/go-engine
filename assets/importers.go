package assets

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/walesey/freetype/truetype"
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

func LoadFont(fontfile string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Printf("Error Reading from font file: %v\n", err)
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Printf("Error Parsing font file: %v\n", err)
		return nil, err
	}
	return f, nil
}
