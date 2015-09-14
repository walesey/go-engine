package physics

import (
	"github.com/stretchr/testify/assert"
	"github.com/walesey/go-engine/vectormath"
	"testing"
)

func TestPhysicsObject(t *testing.T) {
	object := CreatePhysicsObject()
	object.Velocity = vectormath.Vector3{1, 2, 3}
	object.AngularVelocity = vectormath.Quaternion{0, 1, 0, 5}
	object.doStep(1)
	assert.True(t, vectormath.Vector3{1, 2, 3}.ApproxEqual(object.Position, 0.001), "physics object velocity")
	object.doStep(1)
	assert.True(t, vectormath.AngleAxis(10, vectormath.Vector3{0, 1, 0}).ApproxEqual(object.Orientation, 0.001), "physics object angular veloctiy")
}
