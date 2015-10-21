package physics

import (
	"github.com/stretchr/testify/assert"
	"github.com/walesey/go-engine/vectormath"
	"testing"
)

func TestPhysicsObject(t *testing.T) {
	object := NewPhysicsObject()
	object.Velocity = vectormath.Vector3{1, 2, 3}
	object.AngularVelocity = vectormath.Quaternion{0, 1, 0, 5}
	object.doStep(1)
	assert.True(t, vectormath.Vector3{1, 2, 3}.ApproxEqual(object.Position, 0.001), "physics object velocity")
	object.doStep(1)
	assert.True(t, vectormath.AngleAxis(10, vectormath.Vector3{0, 1, 0}).ApproxEqual(object.Orientation, 0.001), "physics object angular veloctiy")
}

func TestAddRemoveObjects(t *testing.T) {
	object1 := NewPhysicsObject()
	object1.Velocity = vectormath.Vector3{1, 0, 0}
	object1.Position = vectormath.Vector3{0, 0, 0}

	object2 := NewPhysicsObject()
	object2.Velocity = vectormath.Vector3{1, 0, 0}
	object2.Position = vectormath.Vector3{0, 0, 0}

	world := NewPhysicsSpace()
	world.StepDt = 1
	world.Add(&object1, &object2)
	world.PhysicsStep()

	assert.True(t, vectormath.Vector3{1, 0, 0}.ApproxEqual(object1.Position, 0.001), "physics object should be updated by physicsSpace")
	assert.True(t, vectormath.Vector3{1, 0, 0}.ApproxEqual(object2.Position, 0.001), "physics object should be updated by physicsSpace")

	world.Remove(&object1, &object2)
	world.PhysicsStep()

	assert.True(t, vectormath.Vector3{1, 0, 0}.ApproxEqual(object1.Position, 0.001), "physics object should not be updated by physicsSpace")
	assert.True(t, vectormath.Vector3{1, 0, 0}.ApproxEqual(object2.Position, 0.001), "physics object should not be updated by physicsSpace")
}
