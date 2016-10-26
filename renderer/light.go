package renderer

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Light struct {
	Directional                bool
	Ambient, Diffuse, Specular [3]float32
	translation, direction     mgl32.Vec3
	dirtyFlag                  bool
}

func CreateLight() *Light {
	light := &Light{
		translation: mgl32.Vec3{},
		direction:   mgl32.Vec3{1, 0, 0},
		Ambient:     [3]float32{0.3, 0.3, 0.3},
		Diffuse:     [3]float32{0.3, 0.3, 0.3},
		Specular:    [3]float32{0.3, 0.3, 0.3},
		dirtyFlag:   true,
	}
	return light
}

func (l *Light) Render(renderer Renderer, index int) {
	if l.dirtyFlag {
		l.dirtyFlag = false

		position := l.translation
		if l.Directional {
			position = l.direction
		}

		renderer.CreateLight(
			l.Ambient[0], l.Ambient[1], l.Ambient[2],
			l.Diffuse[0], l.Diffuse[1], l.Diffuse[2],
			l.Specular[0], l.Specular[1], l.Specular[2],
			l.Directional,
			position,
			index,
		)
	}
}

func (l *Light) SetScale(scale mgl32.Vec3) {} // na

func (l *Light) SetTranslation(translation mgl32.Vec3) {
	l.translation = translation
	l.dirtyFlag = true
}

func (l *Light) SetOrientation(orientation mgl32.Quat) {
	l.direction = orientation.Rotate(mgl32.Vec3{1, 0, 0})
	l.dirtyFlag = true
}
