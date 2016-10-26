package renderer

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

// The camera Entity
type Camera struct {
	renderer                Renderer
	translation, lookat, up mgl32.Vec3
	angle, near, far        float32
	ortho                   bool
}

func CreateCamera(renderer Renderer) *Camera {
	cam := Camera{
		renderer:    renderer,
		translation: mgl32.Vec3{1, 0, 0},
		lookat:      mgl32.Vec3{0, 0, 0},
		up:          mgl32.Vec3{0, 1, 0},
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
	c.renderer.Perspective(c.translation, c.lookat, c.up, float32(c.angle), float32(c.near), float32(c.far))
}

func (c *Camera) GetDirection() mgl32.Vec3 {
	return c.lookat.Sub(c.translation).Normalize()
}

func (c *Camera) GetTranslation() mgl32.Vec3 {
	return c.translation
}

func (c *Camera) GetMouseVector(mouse mgl32.Vec2) mgl32.Vec3 {
	window := c.renderer.WindowDimensions()
	v, err := mgl32.UnProject(
		mgl32.Vec3{mouse.X(), window.Y() - mouse.Y(), 0.5},
		mgl32.LookAtV(c.translation, c.lookat, c.up),
		mgl32.Perspective(mgl32.DegToRad(float32(c.angle)), float32(window.X())/float32(window.Y()), float32(c.near), float32(c.far)),
		0, 0, int(window.X()), int(window.Y()),
	)
	if err == nil {
		return v.Sub(c.translation).Normalize()
	} else {
		log.Println("Error converting camera vector: ", err)
	}
	return c.lookat
}

func (c *Camera) SetScale(scale mgl32.Vec3) {} //na

func (c *Camera) SetTranslation(translation mgl32.Vec3) {
	c.translation = translation
	c.Perspective()
}

func (c *Camera) SetOrientation(orientation mgl32.Quat) {
	direction := orientation.Rotate(mgl32.Vec3{9999999999, 0, 0})
	c.lookat = c.translation.Add(direction)
	c.Perspective()
}

func (c *Camera) SetLookat(lookat mgl32.Vec3) {
	c.lookat = lookat
	c.Perspective()
}

func (c *Camera) SetUp(up mgl32.Vec3) {
	c.up = up
	c.Perspective()
}

func (c *Camera) SetAngle(angle float32) {
	c.angle = angle
	c.Perspective()
}

func (c *Camera) SetNear(near float32) {
	c.near = near
	c.Perspective()
}

func (c *Camera) SetFar(far float32) {
	c.far = far
	c.Perspective()
}
