package ui

import (
	"image"

	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type ImageElement struct {
	id            string
	Hitbox        Hitbox
	width, height float64
	rotation      float64
	offset        vmath.Vector2
	node          *renderer.Node
	img           image.Image
}

func (ie *ImageElement) Render(size, offset vmath.Vector2) vmath.Vector2 {
	width, height := ie.width, ie.height
	if width <= 0 && height <= 0 {
		width = float64(ie.img.Bounds().Size().X)
		height = float64(ie.img.Bounds().Size().Y)
	} else if width <= 0 {
		width = height * float64(ie.img.Bounds().Size().X) / float64(ie.img.Bounds().Size().Y)
	} else if height <= 0 {
		height = width * float64(ie.img.Bounds().Size().Y) / float64(ie.img.Bounds().Size().X)
	}
	imgSize := vmath.Vector2{float64(width), float64(height)}
	ie.node.SetScale(imgSize.ToVector3())
	ie.node.SetTranslation(offset.ToVector3())
	ie.offset = offset
	ie.Hitbox.SetSize(imgSize)
	return imgSize
}

func (ie *ImageElement) Spatial() renderer.Spatial {
	return ie.node
}

func (ie *ImageElement) GetId() string {
	return ie.id
}

func (ie *ImageElement) SetWidth(width float64) {
	ie.width = width
}

func (ie *ImageElement) SetHeight(height float64) {
	ie.height = height
}

func (ie *ImageElement) SetRotation(rotation float64) {
	ie.rotation = rotation
}

func (ie *ImageElement) SetImage(img image.Image) {
	ie.node.RemoveAll(true)
	box := renderer.CreateBoxWithOffset(1, 1, 0, 0)
	mat := renderer.CreateMaterial()
	mat.Diffuse = img
	mat.LightingMode = renderer.MODE_UNLIT
	box.Material = mat
	ie.img = img
	ie.node.Add(box)
}

func (ie *ImageElement) mouseMove(position vmath.Vector2) {
	offsetPos := position.Subtract(ie.offset)
	ie.Hitbox.MouseMove(offsetPos)
}

func (ie *ImageElement) mouseClick(button int, release bool, position vmath.Vector2) {
	offsetPos := position.Subtract(ie.offset)
	ie.Hitbox.MouseClick(button, release, offsetPos)
}

func (ie *ImageElement) keyClick(key string, release bool) {}

func NewImageElement(img image.Image) *ImageElement {
	imageElement := &ImageElement{
		rotation: 0,
		Hitbox:   NewHitbox(),
		node:     renderer.CreateNode(),
	}
	imageElement.SetImage(img)
	return imageElement
}
