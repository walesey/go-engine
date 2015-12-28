package dynamics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	vmath "github.com/walesey/go-engine/vectormath"
)

func TestPhysicsObject(t *testing.T) {
	object := newPhysicsObject()
	object.Velocity = vmath.Vector3{1, 2, 3}
	object.AngularVelocity = vmath.Quaternion{0, 1, 0, 5}
	object.doStep(1)
	assert.True(t, vmath.Vector3{1, 2, 3}.ApproxEqual(object.Position, 0.001), "physics object velocity")
	object.doStep(1)
	assert.True(t, vmath.AngleAxis(10, vmath.Vector3{0, 1, 0}).ApproxEqual(object.Orientation, 0.001), "physics object angular veloctiy")
}

func TestAddRemoveObjects(t *testing.T) {
	world := NewPhysicsSpace()
	world.StepDt = 1

	object1 := world.CreateObject()
	object1.Velocity = vmath.Vector3{1, 0, 0}
	object1.Position = vmath.Vector3{0, 0, 0}

	object2 := world.CreateObject()
	object2.Velocity = vmath.Vector3{1, 0, 0}
	object2.Position = vmath.Vector3{0, 0, 0}

	world.DoStep()

	assert.True(t, vmath.Vector3{1, 0, 0}.ApproxEqual(object1.Position, 0.001), "physics object should be updated by physicsSpace")
	assert.True(t, vmath.Vector3{1, 0, 0}.ApproxEqual(object2.Position, 0.001), "physics object should be updated by physicsSpace")

	world.Remove(object1, object2)
	world.DoStep()

	assert.True(t, vmath.Vector3{1, 0, 0}.ApproxEqual(object1.Position, 0.001), "physics object should not be updated by physicsSpace")
	assert.True(t, vmath.Vector3{1, 0, 0}.ApproxEqual(object2.Position, 0.001), "physics object should not be updated by physicsSpace")
}
