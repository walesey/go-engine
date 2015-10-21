package physics

import (
	vmath "github.com/walesey/go-engine/vectormath"
)

type PhysicsObject struct {
	Position, Velocity           vmath.Vector3
	Orientation, AngularVelocity vmath.Quaternion
	broadPhase, narrowPhase      Collider
}

func NewPhysicsObject() PhysicsObject {
	return PhysicsObject{
		Position:        vmath.Vector3{0, 0, 0},
		Velocity:        vmath.Vector3{0, 0, 0},
		Orientation:     vmath.IdentityQuaternion(),
		AngularVelocity: vmath.Quaternion{1, 0, 0, 0},
	}
}

//NarrowPhaseOverlap
func (obj PhysicsObject) NarrowPhaseOverlap(other PhysicsObject) bool {
	return obj.narrowPhase.Overlap(other.narrowPhase)
}

//BroadPhaseOverlap
func (obj PhysicsObject) BroadPhaseOverlap(other PhysicsObject) bool {
	return obj.broadPhase.Overlap(other.broadPhase)
}

func (obj *PhysicsObject) doStep(dt float64) {
	//apply position increment
	obj.Position = obj.Position.Add(obj.Velocity.MultiplyScalar(dt))

	//apply orientation increment
	axis := vmath.Vector3{obj.AngularVelocity.X, obj.AngularVelocity.Y, obj.AngularVelocity.Z}
	obj.Orientation = vmath.AngleAxis(dt*obj.AngularVelocity.W, axis).Multiply(obj.Orientation)
}
