package renderer

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

// The camera Entity
type Camera struct {
	Translation, Lookat, Up mgl32.Vec3
	Angle, Near, Far        float32
	Ortho                   bool
}

func CreateCamera() *Camera {
	cam := Camera{
		Translation: mgl32.Vec3{0, 0, 0},
		Lookat:      mgl32.Vec3{1, 0, 0},
		Up:          mgl32.Vec3{0, 1, 0},
		Angle:       45.0,
		Near:        0.1,
		Far:         999999999.0,
	}

	return &cam
}

func (c *Camera) GetDirection() mgl32.Vec3 {
	return c.Lookat.Sub(c.Translation).Normalize()
}

func (c *Camera) GetMouseVector(windowSize mgl32.Vec2, mouse mgl32.Vec2) mgl32.Vec3 {
	v, err := mgl32.UnProject(
		mgl32.Vec3{mouse.X(), windowSize.Y() - mouse.Y(), 0.5},
		mgl32.LookAtV(c.Translation, c.Lookat, c.Up),
		mgl32.Perspective(mgl32.DegToRad(c.Angle), windowSize.X()/windowSize.Y(), c.Near, c.Far),
		0, 0, int(windowSize.X()), int(windowSize.Y()),
	)
	if err == nil {
		return v.Sub(c.Translation).Normalize()
	} else {
		log.Println("Error converting camera vector: ", err)
	}
	return c.Lookat
}

func (c *Camera) SetScale(scale mgl32.Vec3) {} //na

func (c *Camera) SetTranslation(translation mgl32.Vec3) {
	c.Translation = translation
}

func (c *Camera) SetOrientation(orientation mgl32.Quat) {
	direction := orientation.Rotate(mgl32.Vec3{1, 0, 0})
	c.Lookat = c.Translation.Add(direction)
}
