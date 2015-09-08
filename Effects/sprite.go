package effects

import (
	"image/color"
	"github.com/Walesey/goEngine/renderer"
)

type Sprite struct {
	geometry renderer.Geometry
	frame, totalFrames, framesX, framesY int
}

func CreateSprite( totalFrames, framesX, framesY int, material *renderer.Material ) Sprite{
	sprite := Sprite{
		frame: 0,
		totalFrames: totalFrames,
		framesX: framesX,
		framesY: framesY,
	}
	geometry := renderer.CreateBox(1,1)
	sprite.geometry = geometry
	sprite.geometry.Material = material
    sprite.geometry.CullBackface = false
    sprite.geometry.Flipbook.FrameSizeX = 1.0/float32(framesX)
    sprite.geometry.Flipbook.FrameSizeY = 1.0/float32(framesY)
    return sprite
}

func (sprite *Sprite) Draw( renderer renderer.Renderer ) {
    //flipbook
    sprite.geometry.Flipbook.IndexX = sprite.frame % sprite.framesX
    sprite.geometry.Flipbook.IndexY = sprite.framesY - (sprite.frame / sprite.framesY) - 1
	sprite.geometry.Draw(renderer)
}

//pick a frame as a ratio (0.0 to 1.0)
func (sprite *Sprite) FrameLerp( progress float64 ) {
	sprite.frame = int( progress * float64(sprite.totalFrames) )
}

func (sprite *Sprite) NextFrame() {
	sprite.frame = sprite.frame + 1
	if sprite.frame > sprite.totalFrames {
		sprite.frame = 0
	}
}

func (sprite *Sprite) SetColor(color color.NRGBA) {
	sprite.geometry.Color = color
}