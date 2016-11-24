package ui

import (
	"image"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

type ImageElement struct {
	id            string
	Hitbox        Hitbox
	percentWidth  bool
	percentHeight bool
	width, height float32
	rotation      float32
	size, offset  mgl32.Vec2
	node          *renderer.Node
	img           image.Image
}

func (ie *ImageElement) Render(size, offset mgl32.Vec2) mgl32.Vec2 {
	ie.size, ie.offset = size, offset
	width, height := ie.getWidth(size.X()), ie.getHeight(size.X())
	if ie.img != nil {
		if width <= 0 && height <= 0 {
			width = float32(ie.img.Bounds().Size().X)
			height = float32(ie.img.Bounds().Size().Y)
		} else if width <= 0 {
			width = height * float32(ie.img.Bounds().Size().X) / float32(ie.img.Bounds().Size().Y)
		} else if height <= 0 {
			height = width * float32(ie.img.Bounds().Size().Y) / float32(ie.img.Bounds().Size().X)
		}
	}
	imgSize := mgl32.Vec2{width, height}
	ie.node.SetScale(imgSize.Vec3(0))
	ie.node.SetTranslation(offset.Vec3(0))
	ie.offset = offset
	ie.Hitbox.SetSize(imgSize)
	return imgSize
}

func (ie *ImageElement) ReRender() {
	ie.Render(ie.size, ie.offset)
}

func (ie *ImageElement) Spatial() renderer.Spatial {
	return ie.node
}

func (ie *ImageElement) GetId() string {
	return ie.id
}

func (ie *ImageElement) SetWidth(width float32) {
	ie.width = width
}

func (ie *ImageElement) UsePercentWidth(usePercent bool) {
	ie.percentWidth = usePercent
}

func (ie *ImageElement) SetHeight(height float32) {
	ie.height = height
}

func (ie *ImageElement) UsePercentHeight(usePercent bool) {
	ie.percentHeight = usePercent
}

func (ie *ImageElement) getWidth(parentWidth float32) float32 {
	if ie.percentWidth {
		return parentWidth * ie.width / 100.0
	}
	return ie.width
}

func (ie *ImageElement) getHeight(parentWidth float32) float32 {
	if ie.percentHeight {
		return parentWidth * ie.height / 100.0
	}
	return ie.height
}

func (ie *ImageElement) SetRotation(rotation float32) {
	ie.rotation = rotation
}

func (ie *ImageElement) SetImage(img image.Image) {
	mat := renderer.NewMaterial(renderer.NewTexture("diffuseMap", img, false))
	ie.node.Material = mat
	ie.img = img
}

func (ie *ImageElement) mouseMove(position mgl32.Vec2) {
	offsetPos := position.Sub(ie.offset)
	ie.Hitbox.MouseMove(offsetPos)
}

func (ie *ImageElement) mouseClick(button int, release bool, position mgl32.Vec2) {
	offsetPos := position.Sub(ie.offset)
	ie.Hitbox.MouseClick(button, release, offsetPos)
}

func (ie *ImageElement) keyClick(key string, release bool) {}

func NewImageElement(img image.Image) *ImageElement {
	imageElement := &ImageElement{
		rotation: 0,
		Hitbox:   NewHitbox(),
		node:     renderer.NewNode(),
	}
	box := renderer.CreateBoxWithOffset(1, 1, 0, 0)
	imageElement.node.Add(box)
	imageElement.SetImage(img)
	return imageElement
}
