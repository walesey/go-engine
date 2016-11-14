package assets

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"io/ioutil"

	"github.com/disintegration/imaging"
	"github.com/walesey/go-engine/renderer"
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

func ImportShader(vertexFile, fragmentFile string) (*renderer.Shader, error) {
	vertsrc, err := ioutil.ReadFile(vertexFile)
	if err != nil {
		fmt.Printf("Error vertex file: %v\n", err)
		return nil, err
	}

	fragsrc, err := ioutil.ReadFile(fragmentFile)
	if err != nil {
		fmt.Printf("Error fragment file: %v\n", err)
		return nil, err
	}

	shader := renderer.NewShader()
	shader.VertSrc = string(vertsrc)
	shader.FragSrc = string(fragsrc)
	return shader, nil
}
