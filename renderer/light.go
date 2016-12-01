package renderer

import "github.com/go-gl/mathgl/mgl32"

type LightType int

const (
	_ LightType = iota
	POINT
	DIRECTIONAL
)

type Light struct {
	LightType
	Color     [3]float32 //RGB
	Position  mgl32.Vec3
	Direction mgl32.Vec3
}

func NewLight(lightType LightType) *Light {
	return &Light{
		LightType: lightType,
		Color:     [3]float32{1, 1, 1},
		Direction: mgl32.Vec3{1, 0, 0},
	}
}

func (l *Light) SetScale(scale mgl32.Vec3) {} //na

func (l *Light) SetTranslation(translation mgl32.Vec3) {
	if l.LightType == POINT {
		l.Position = translation
	}
}

func (l *Light) SetOrientation(orientation mgl32.Quat) {
	if l.LightType == DIRECTIONAL {
		l.Direction = orientation.Rotate(mgl32.Vec3{1, 0, 0})
	}
}
