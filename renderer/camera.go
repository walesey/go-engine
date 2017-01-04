package renderer

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/util"
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

// GetMouseVector - Returns a normal vector given by the mouse position
func (c *Camera) GetMouseVector(windowSize, mouse mgl32.Vec2) mgl32.Vec3 {
	v, err := mgl32.UnProject(
		mgl32.Vec3{mouse.X(), windowSize.Y() - mouse.Y(), 0.5},
		mgl32.LookAtV(c.Translation, c.Lookat, c.Up),
		mgl32.Perspective(mgl32.DegToRad(c.Angle), windowSize.X()/windowSize.Y(), c.Near, c.Far),
		0, 0, int(windowSize.X()), int(windowSize.Y()),
	)
	if err != nil {
		log.Println("Error converting camera vector: ", err)
		return c.Lookat
	}

	return v.Sub(c.Translation).Normalize()
}

// GetWindowVector - Returns the screen position of the given world vector
func (c *Camera) GetWindowVector(windowSize mgl32.Vec2, point mgl32.Vec3) mgl32.Vec3 {
	v := mgl32.Project(
		point,
		mgl32.LookAtV(c.Translation, c.Lookat, c.Up),
		mgl32.Perspective(mgl32.DegToRad(c.Angle), windowSize.X()/windowSize.Y(), c.Near, c.Far),
		0, 0, int(windowSize.X()), int(windowSize.Y()),
	)

	return mgl32.Vec3{v.X(), windowSize.Y() - v.Y(), v.Z()}
}

// CameraContainsSphere - determines if a sphere in contained in the frustrum given by the camera.
// sphere is given by point and radius
func (c *Camera) CameraContainsSphere(windowSize mgl32.Vec2, radius float32, point mgl32.Vec3) bool {
	return c.FrustrumContainsSphere(windowSize, mgl32.Vec2{}, windowSize, radius, point)
}

// FrustrumContainsSphere - determines if a sphere in contained in the frustrum given by start/end vectors on the screen.
// sphere is given by point and radius
func (c *Camera) FrustrumContainsSphere(windowSize, start, end mgl32.Vec2, radius float32, point mgl32.Vec3) bool {
	tlv := c.GetMouseVector(windowSize, start)
	trv := c.GetMouseVector(windowSize, mgl32.Vec2{end.X(), start.Y()})
	blv := c.GetMouseVector(windowSize, mgl32.Vec2{start.X(), end.Y()})
	brv := c.GetMouseVector(windowSize, end)

	r := radius
	delta := point.Sub(c.Translation)
	if point.ApproxEqual(c.Translation) {
		delta = mgl32.Vec3{1, 0, 0}
	}

	return delta.Dot(tlv.Cross(trv).Normalize()) > -r &&
		delta.Dot(trv.Cross(brv).Normalize()) > -r &&
		delta.Dot(brv.Cross(blv).Normalize()) > -r &&
		delta.Dot(blv.Cross(tlv).Normalize()) > -r &&
		util.Vec3LenSq(delta) < (c.Far+r)*(c.Far+r)
}

func (c *Camera) SetScale(scale mgl32.Vec3) {} //na

func (c *Camera) SetTranslation(translation mgl32.Vec3) {
	c.Translation = translation
}

func (c *Camera) SetOrientation(orientation mgl32.Quat) {
	direction := orientation.Rotate(mgl32.Vec3{1, 0, 0})
	c.Lookat = c.Translation.Add(direction)
}
