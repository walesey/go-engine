package renderer

import (
	"image"
)

type Texture struct {
	TextureId   uint32
	TextureName string
	Img         image.Image
	Loaded      bool
	CubeMap     *CubeMap
}

type Material struct {
	Textures []*Texture
}

func NewTexture(name string, img image.Image) *Texture {
	return &Texture{
		TextureName: name,
		Img:         img,
	}
}

func NewMaterial(textures ...*Texture) *Material {
	return &Material{
		Textures: textures,
	}
}
