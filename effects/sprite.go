package effects

import (
	"image/color"

	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/vectormath"
)

type Sprite struct {
	node                                 *renderer.Node
	geometry                             *renderer.Geometry
	transform                            renderer.Transform
	location, scale                      vectormath.Vector3
	Rotation                             float64
	FaceCamera                           bool
	frame, totalFrames, framesX, framesY int
}

func CreateSprite(totalFrames, framesX, framesY int, material *renderer.Material) *Sprite {
	sprite := Sprite{
		transform:   renderer.CreateTransform(),
		location:    vectormath.Vector3{0, 0, 0},
		scale:       vectormath.Vector3{1, 1, 1},
		frame:       0,
		FaceCamera:  false,
		totalFrames: totalFrames,
		framesX:     framesX,
		framesY:     framesY,
	}
	geometry := renderer.CreateBox(1, 1)
	sprite.geometry = geometry
	sprite.geometry.Material = material
	sprite.geometry.CullBackface = false
	sprite.node = renderer.CreateNode()
	sprite.node.Add(sprite.geometry)
	return &sprite
}

func BoxFlipbook(geometry *renderer.Geometry, frame, framesX, framesY int) {
	frameSizeX := 1.0 / float32(framesX)
	frameSizeY := 1.0 / float32(framesY)
	indexX := float32(frame % framesX)
	indexY := float32(framesY - (frame / framesY) - 1)
	u1, u2 := frameSizeX*indexX, frameSizeX*(indexX+1.0)
	v1, v2 := frameSizeY*indexY, frameSizeY*(indexY+1.0)
	geometry.SetUVs(u1, v1, u2, v1, u2, v2, u1, v2)
}

func (sprite *Sprite) Draw(r renderer.Renderer) {
	if sprite.FaceCamera {
		orientation := vectormath.FacingOrientation(sprite.Rotation, r.CameraLocation().Subtract(sprite.location), vectormath.Vector3{0, 0, 1}, vectormath.Vector3{-1, 0, 0})
		sprite.node.Transform.From(sprite.scale, sprite.location, orientation)
	}
	sprite.node.Draw(r)
}

func (sprite *Sprite) Destroy(renderer renderer.Renderer) {
	sprite.node.Destroy(renderer)
}

func (sprite *Sprite) Centre() vectormath.Vector3 {
	return sprite.location
}

func (sprite *Sprite) NextFrame() {
	sprite.frame = sprite.frame + 1
	if sprite.frame >= sprite.totalFrames {
		sprite.frame = 0
	}
	BoxFlipbook(sprite.geometry, sprite.frame, sprite.framesX, sprite.framesY)
}

func (sprite *Sprite) SetColor(color color.NRGBA) {
	sprite.geometry.SetColor(color)
}

func (sprite *Sprite) Optimize(geometry *renderer.Geometry, transform renderer.Transform) {
	sprite.geometry.Optimize(geometry, transform)
}

func (sprite *Sprite) SetTranslation(translation vectormath.Vector3) {
	sprite.location = translation
	sprite.node.SetTranslation(translation)
}

func (sprite *Sprite) SetScale(scale vectormath.Vector3) {
	sprite.scale = scale
	sprite.node.SetScale(scale)
}

func (sprite *Sprite) SetOrientation(orientation vectormath.Quaternion) {
	sprite.node.SetOrientation(orientation)
}
