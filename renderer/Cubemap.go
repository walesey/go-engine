package renderer

import (
	"image"

	"github.com/disintegration/imaging"
)

type CubeMap struct {
	Id     uint32
	Lod    bool
	Loaded bool

	Right, Left, Top, Bottom, Back, Front image.Image
}

func NewCubemap(baseImage image.Image, lod bool) *CubeMap {
	cubeMap := new(CubeMap)

	x := baseImage.Bounds().Max.X
	y := baseImage.Bounds().Max.Y

	cubeMap.Lod = lod
	cubeMap.Right = imaging.Crop(baseImage, image.Rect(x/2, y/3, 3*x/4, 2*y/3))
	cubeMap.Left = imaging.Crop(baseImage, image.Rect(0, y/3, x/4, 2*y/3))
	cubeMap.Top = imaging.Crop(baseImage, image.Rect(x/4, 0, x/2, y/3))
	cubeMap.Bottom = imaging.Crop(baseImage, image.Rect(x/4, 2*y/3, x/2, y))
	cubeMap.Back = imaging.Crop(baseImage, image.Rect(3*x/4, y/3, x, 2*y/3))
	cubeMap.Front = imaging.Crop(baseImage, image.Rect(x/4, y/3, x/2, 2*y/3))

	return cubeMap
}

func (cm *CubeMap) Destroy(renderer Renderer) {
	if cm != nil {
		renderer.DestroyCubeMap(cm)
		cm.Loaded = false
	}
}
