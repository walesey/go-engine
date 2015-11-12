package physics

import (
	vmath "github.com/walesey/go-engine/vectormath"
)

type PhysicsObject struct {
	Position, Velocity           vmath.Vector3
	Orientation, AngularVelocity vmath.Quaternion
	Mass                         float64
	Static                       bool //disables movement
	ForceStore                   *ForceStore
	BroadPhase, NarrowPhase      Collider
}

type PhysicsObjectPool struct {
	pool []*PhysicsObject
}

func NewPhysicsObjectPool() *PhysicsObjectPool {
	return &PhysicsObjectPool{
		pool: make([]*PhysicsObject, 0, 0),
	}
}

func (objPool PhysicsObjectPool) GetPhysicsObject() *PhysicsObject {
	if len(objPool.pool) > 0 {
		obj := objPool.pool[len(objPool.pool)-1]
		objPool.pool = objPool.pool[:len(objPool.pool)-1]
		return obj
	}
	return NewPhysicsObject()
}

func (objPool PhysicsObjectPool) ReleasePhysicsObject(obj *PhysicsObject) {
	objPool.pool = append(objPool.pool, obj)
}

func NewPhysicsObject() *PhysicsObject {
	return &PhysicsObject{
		Position:        vmath.Vector3{0, 0, 0},
		Velocity:        vmath.Vector3{0, 0, 0},
		Orientation:     vmath.IdentityQuaternion(),
		AngularVelocity: vmath.Quaternion{1, 0, 0, 0},
		Mass:            1.0,
		ForceStore:      NewForceStore(),
	}
}

//NarrowPhaseOverlap
func (obj *PhysicsObject) NarrowPhaseOverlap(other *PhysicsObject) bool {
	if obj.NarrowPhase == nil || other.NarrowPhase == nil {
		return false
	}
	obj.NarrowPhase.Offset(obj.Position, obj.Orientation)
	other.NarrowPhase.Offset(other.Position, other.Orientation)
	return obj.NarrowPhase.Overlap(other.NarrowPhase)
}

//BroadPhaseOverlap
func (obj *PhysicsObject) BroadPhaseOverlap(other *PhysicsObject) bool {
	if obj.BroadPhase == nil || other.BroadPhase == nil {
		return false
	}
	obj.BroadPhase.Offset(obj.Position, obj.Orientation)
	other.BroadPhase.Offset(other.Position, other.Orientation)
	return obj.BroadPhase.Overlap(other.BroadPhase)
}

func (obj *PhysicsObject) PenetrationVector(other *PhysicsObject) vmath.Vector3 {
	if obj.NarrowPhase == nil || other.NarrowPhase == nil {
		return vmath.Vector3{}
	}
	return obj.NarrowPhase.PenetrationVector(other.NarrowPhase)
}

func (obj *PhysicsObject) doStep(dt float64) {
	//process forces and acceleration
	obj.ForceStore.DoStep(dt, obj)

	//apply position increment
	obj.Position = obj.Position.Add(obj.Velocity.MultiplyScalar(dt))

	//apply orientation increment
	axis := vmath.Vector3{obj.AngularVelocity.X, obj.AngularVelocity.Y, obj.AngularVelocity.Z}
	obj.Orientation = vmath.AngleAxis(dt*obj.AngularVelocity.W, axis).Multiply(obj.Orientation)
}

//handleCollision
func (obj *PhysicsObject) handleCollision(other *PhysicsObject) {

}
