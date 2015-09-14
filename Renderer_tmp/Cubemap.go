package renderer

import (
	"image"
	"github.com/disintegration/imaging"
)

type CubeMap struct {
	Right, Left, Top, Bottom, Back, Front image.Image
}

func CreateCubemap( baseImage image.Image ) *CubeMap {
	cubeMap := new(CubeMap)

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

func (cm *CubeMap) Clone() *CubeMap{
	return &CubeMap{Right: imaging.Clone(cm.Right), Left: imaging.Clone(cm.Left), Top: imaging.Clone(cm.Top), Bottom: imaging.Clone(cm.Bottom), Back: imaging.Clone(cm.Back), Front: imaging.Clone(cm.Front)}
}

func (cm *CubeMap) Resize(size int) {
	cm.Front = imaging.Resize(cm.Front, size, size, imaging.Gaussian)
	cm.Back = imaging.Resize(cm.Back, size, size, imaging.Gaussian)
	cm.Bottom = imaging.Resize(cm.Bottom, size, size, imaging.Gaussian)
	cm.Top = imaging.Resize(cm.Top, size, size, imaging.Gaussian)
	cm.Left = imaging.Resize(cm.Left, size, size, imaging.Gaussian)
	cm.Right = imaging.Resize(cm.Right, size, size, imaging.Gaussian)
}