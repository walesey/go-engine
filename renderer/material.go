package renderer

import (
	"image"
)

type Texture struct {
	TextureId   uint32
	TextureName string
	Img         image.Image
	Lod         bool
	Loaded      bool
}

type Material struct {
	Textures []*Texture
}

func NewTexture(name string, img image.Image, lod bool) *Texture {
	return &Texture{
		TextureName: name,
		Img:         img,
		Lod:         lod,
	}
}

func NewMaterial(textures ...*Texture) *Material {
	return &Material{
		Textures: textures,
	}
}

func (m *Material) Destroy(renderer Renderer) {
	renderer.DestroyMaterial(m)
}
