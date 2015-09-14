package physics

import (
	"github.com/walesey/go-engine/vectormath"
)

type Collider interface {
	Overlap(other Collider) bool
}

type BoundingBox struct {
	bounds vectormath.Vector3
}

type PhysicsObject struct {
	Position, Velocity           vectormath.Vector3
	Orientation, AngularVelocity vectormath.Quaternion
	broadPhase, narrowPhase      Collider
}

type PhysicsSpace struct {
	objects []*PhysicsObject
	StepDt  float64
}

func CreatePhysicsSpace() PhysicsSpace {
	return PhysicsSpace{
		StepDt:  0.018,
		objects: make([]*PhysicsObject, 0, 500),
	}
}

func CreatePhysicsObject() PhysicsObject {
	return PhysicsObject{
		Position:        vectormath.Vector3{0, 0, 0},
		Velocity:        vectormath.Vector3{0, 0, 0},
		Orientation:     vectormath.IdentityQuaternion(),
		AngularVelocity: vectormath.Quaternion{1, 0, 0, 0},
	}
}

func (ps *PhysicsSpace) PhysicsStep() {
	for _, object := range ps.objects {
		object.doStep(ps.StepDt)
	}
}

func (obj *PhysicsObject) doStep(dt float64) {
	obj.Position = obj.Position.Add(obj.Velocity.MultiplyScalar(dt))
	axis := vectormath.Vector3{obj.AngularVelocity.X, obj.AngularVelocity.Y, obj.AngularVelocity.Z}
	obj.Orientation = obj.Orientation.Multiply(vectormath.AngleAxis(dt*obj.AngularVelocity.W, axis))
}
