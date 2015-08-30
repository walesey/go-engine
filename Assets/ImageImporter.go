package assets

import (
	"os"
	"log"
	"image"
	_ "image/png"
	_ "image/jpeg"

	"github.com/disintegration/imaging"
)

type CubeMapData struct {
	Right, Left, Top, Bottom, Back, Front image.Image
}

func ImportImage( file string ) image.Image {
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

func ImportCubemap( file string ) *CubeMapData {
	cubeMap := new(CubeMapData)

	baseImage := ImportImage(file)
	x := baseImage.Bounds().Max.X 
	y := baseImage.Bounds().Max.Y

	cubeMap.Right = imaging.Crop( baseImage, image.Rect(x/2, y/3, 3*x/4, 2*y/3) )
	cubeMap.Left = imaging.Crop( baseImage, image.Rect(0, y/3, x/4, 2*y/3) )
	cubeMap.Top = imaging.Crop( baseImage, image.Rect(x/4, 0, x/2, y/3) )
	cubeMap.Bottom = imaging.Crop( baseImage, image.Rect(x/4, 2*y/3, x/2, y) )
	cubeMap.Back = imaging.Crop( baseImage, image.Rect(3*x/4, y/3, x, 2*y/3) )
	cubeMap.Front = imaging.Crop( baseImage, image.Rect(x/4, y/3, x/2, 2*y/3) )

	return cubeMap
}

func (cm *CubeMapData) Clone() *CubeMapData{
	return &CubeMapData{Right: imaging.Clone(cm.Right), Left: imaging.Clone(cm.Left), Top: imaging.Clone(cm.Top), Bottom: imaging.Clone(cm.Bottom), Back: imaging.Clone(cm.Back), Front: imaging.Clone(cm.Front)}
}

func (cm *CubeMapData) Resize(size int) {
	cm.Front = imaging.Resize(cm.Front, size, size, imaging.Gaussian)
	cm.Back = imaging.Resize(cm.Back, size, size, imaging.Gaussian)
	cm.Bottom = imaging.Resize(cm.Bottom, size, size, imaging.Gaussian)
	cm.Top = imaging.Resize(cm.Top, size, size, imaging.Gaussian)
	cm.Left = imaging.Resize(cm.Left, size, size, imaging.Gaussian)
	cm.Right = imaging.Resize(cm.Right, size, size, imaging.Gaussian)
}