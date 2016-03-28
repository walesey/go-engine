package renderer

import (
	vmath "github.com/walesey/go-engine/vectormath"
)

//A camera is an Entity
type Camera struct {
	renderer                Renderer
	translation, lookat, up vmath.Vector3
	angle, near, far        float64
	ortho                   bool
}

func CreateCamera(renderer Renderer) *Camera {
	cam := Camera{
		renderer:    renderer,
		translation: vmath.Vector3{1, 0, 0},
		lookat:      vmath.Vector3{0, 0, 0},
		up:          vmath.Vector3{0, 1, 0},
		angle:       45.0,
		near:        0.1,
		far:         999999999.0,
	}

	return &cam
}

func (c *Camera) Ortho() {
	c.renderer.Ortho()
}

func (c *Camera) Perspective() {
	c.renderer.Projection(float32(c.angle), float32(c.near), float32(c.far))
	c.renderer.Camera(c.translation, c.lookat, c.up)
}

func (c *Camera) GetDirection() vmath.Vector3 {
	return c.lookat.Subtract(c.translation).Normalize()
}

func (c *Camera) GetTranslation() vmath.Vector3 {
	return c.translation
}

func (c *Camera) SetScale(scale vmath.Vector3) {} //na

func (c *Camera) SetTranslation(translation vmath.Vector3) {
	c.translation = translation
	c.Perspective()
}

func (c *Camera) SetOrientation(orientation vmath.Quaternion) {
	direction := orientation.Apply(vmath.Vector3{9999999999, 0, 0})
	c.lookat = c.translation.Add(direction)
	c.Perspective()
}

func (c *Camera) SetUp(up vmath.Vector3) {
	c.up = up
	c.Perspective()
}

func (c *Camera) SetAngle(angle float64) {
	c.angle = angle
	c.Perspective()
}

func (c *Camera) SetNear(near float64) {
	c.near = near
	c.Perspective()
}

func (c *Camera) SetFar(far float64) {
	c.far = far
	c.Perspective()
}
