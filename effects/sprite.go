package effects

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/util"
)

type Sprite struct {
	node                                 *renderer.Node
	geometry                             *renderer.Geometry
	Rotation                             float32
	FaceCamera                           bool
	cameraPosition                       mgl32.Vec3
	frame, totalFrames, framesX, framesY int
}

func CreateSprite(totalFrames, framesX, framesY int, material *renderer.Material) *Sprite {
	sprite := Sprite{
		frame:       0,
		FaceCamera:  true,
		totalFrames: totalFrames,
		framesX:     framesX,
		framesY:     framesY,
	}
	geometry := renderer.CreateBox(1, 1)
	sprite.geometry = geometry
	sprite.node = renderer.NewNode()
	sprite.node.Material = material
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

func (sprite *Sprite) Draw(r renderer.Renderer, transform mgl32.Mat4) {
	sprite.node.Draw(r, transform)
}

func (sprite *Sprite) Destroy(renderer renderer.Renderer) {
	sprite.node.Destroy(renderer)
}

func (sprite *Sprite) Centre() mgl32.Vec3 {
	return sprite.node.Centre()
}

func (sprite *Sprite) SetParent(parent *renderer.Node) {
	sprite.node.SetParent(parent)
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

func (sprite *Sprite) Optimize(geometry *renderer.Geometry, transform mgl32.Mat4) {
	sprite.geometry.Optimize(geometry, transform)
}

func (sprite *Sprite) SetTranslation(translation mgl32.Vec3) {
	sprite.node.SetTranslation(translation)
}

func (sprite *Sprite) SetScale(scale mgl32.Vec3) {
	sprite.node.SetScale(scale)
}

func (sprite *Sprite) SetOrientation(orientation mgl32.Quat) {
	sprite.node.SetOrientation(orientation)
}

func (sprite *Sprite) SetCameraLocation(cameraLocation mgl32.Vec3) {
	sprite.cameraPosition = cameraLocation
}

func (sprite *Sprite) Update(dt float64) {
	if sprite.FaceCamera {
		orientation := util.FacingOrientation(sprite.Rotation, sprite.cameraPosition.Sub(sprite.node.Translation), mgl32.Vec3{0, 0, 1}, mgl32.Vec3{-1, 0, 0})
		sprite.node.SetOrientation(orientation)
	}
}
