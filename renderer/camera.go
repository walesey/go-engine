package renderer

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
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

func (c *Camera) GetMouseVector(mouse vmath.Vector2) vmath.Vector3 {
	window := c.renderer.WindowDimensions()
	v, err := mgl32.UnProject(
		mgl32.Vec3{float32(mouse.X), float32(window.Y - mouse.Y), 0.5},
		mgl32.LookAtV(convertVector(c.translation), convertVector(c.lookat), convertVector(c.up)),
		mgl32.Perspective(mgl32.DegToRad(float32(c.angle)), float32(window.X)/float32(window.Y), float32(c.near), float32(c.far)),
		0, 0, int(window.X), int(window.Y),
	)
	if err == nil {
		return vmath.Vector3{float64(v.X()), float64(v.Y()), float64(v.Z())}.Subtract(c.translation).Normalize()
	} else {
		log.Println("Error converting camera vector: ", err)
	}
	return c.lookat
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

func (c *Camera) SetLookat(lookat vmath.Vector3) {
	c.lookat = lookat
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
