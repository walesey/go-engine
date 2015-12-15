package renderer

import (
	vmath "github.com/walesey/go-engine/vectormath"
)

//A camera is an Entity
type Camera struct {
	renderer            Renderer
	Translation, Lookat vmath.Vector3
	Up                  vmath.Vector3
}

func CreateCamera(renderer Renderer) *Camera {
	cam := Camera{
		renderer:    renderer,
		Translation: vmath.Vector3{1, 0, 0},
		Lookat:      vmath.Vector3{0, 0, 0},
		Up:          vmath.Vector3{0, 1, 0},
	}
	return &cam
}

func (c *Camera) GetDirection() vmath.Vector3 {
	return c.Lookat.Subtract(c.Translation)
}

func (c *Camera) SetScale(scale vmath.Vector3) {} //na

func (c *Camera) SetTranslation(translation vmath.Vector3) {
	c.Translation = translation
	c.renderer.Camera(c.Translation, c.Lookat, c.Up)
}

func (c *Camera) SetOrientation(orientation vmath.Quaternion) {
	direction := orientation.Apply(vmath.Vector3{1, 0, 0})
	c.Lookat = c.Translation.Add(direction)
	c.renderer.Camera(c.Translation, c.Lookat, c.Up)
}
