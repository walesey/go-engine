package effects

import (
	"github.com/Walesey/goEngine/vectorMath"
	"github.com/Walesey/goEngine/renderer"
)

type Sprite struct {
	geometry *renderer.Geometry
	node renderer.Node
	frame, totalFrames, framesX, framesY int
}

func CreateSprite( totalFrames, framesX, framesY int, material renderer.Material ) Sprite{
	sprite := Sprite{
		frame: 0,
		totalFrames: totalFrames,
		framesX: framesX,
		framesY: framesY,
	}
	geometry := renderer.CreateBox(1,1)
	sprite.geometry = &geometry
	sprite.geometry.Material = material
	sprite.node = renderer.CreateNode()
    sprite.geometry.CullBackface = false
    sprite.geometry.Flipbook.FrameSizeX = 1.0/float32(framesX)
    sprite.geometry.Flipbook.FrameSizeY = 1.0/float32(framesY)
	sprite.node.Add( sprite.geometry )
    return sprite
}

func (sprite *Sprite) Draw( renderer renderer.Renderer ) {
	//face the camera
    sprite.node.SetFacing( 3.14, renderer.CameraLocation().Subtract(sprite.node.Translation).Normalize(), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,-1} )
    //flipbook
    sprite.NextFrame()
    sprite.geometry.Flipbook.IndexX = sprite.frame % sprite.framesX
    sprite.geometry.Flipbook.IndexY = sprite.framesY - (sprite.frame / sprite.framesY) - 1
	sprite.node.Draw(renderer)
}

func (sprite *Sprite) NextFrame() {
	sprite.frame = sprite.frame + 1
	if sprite.frame > sprite.totalFrames {
		sprite.frame = 0
	}
}

//GameEntity
func (sprite *Sprite) SetScale( scale vectorMath.Vector3 ) {
	sprite.node.SetScale(scale)
}

func (sprite *Sprite) SetTranslation( translation vectorMath.Vector3 ) {
	sprite.node.SetTranslation(translation)
}

func (sprite *Sprite) SetOrientation( orientation vectorMath.Quaternion  ) {
	sprite.node.SetOrientation(orientation)
}

func (sprite *Sprite) SetRotation( angle float64, axis vectorMath.Vector3 ) {
	sprite.node.SetRotation(angle, axis)
}

func (sprite *Sprite) SetFacing( rotation float64, newNormal, normal, tangent vectorMath.Vector3 ) {
	//not implemented
}