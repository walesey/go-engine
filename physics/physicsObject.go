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
	broadPhase, narrowPhase      Collider
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
	newObj := NewPhysicsObject()
	return &newObj
}

func (objPool PhysicsObjectPool) ReleasePhysicsObject(obj *PhysicsObject) {
	objPool.pool = append(objPool.pool, obj)
}

func NewPhysicsObject() PhysicsObject {
	return PhysicsObject{
		Position:        vmath.Vector3{0, 0, 0},
		Velocity:        vmath.Vector3{0, 0, 0},
		Orientation:     vmath.IdentityQuaternion(),
		AngularVelocity: vmath.Quaternion{1, 0, 0, 0},
		Mass:            1.0,
		ForceStore:      NewForceStore(),
	}
}

//NarrowPhaseOverlap
func (obj PhysicsObject) NarrowPhaseOverlap(other PhysicsObject) bool {
	if obj.narrowPhase == nil || other.narrowPhase == nil {
		return false
	}
	return obj.narrowPhase.Overlap(other.narrowPhase)
}

//BroadPhaseOverlap
func (obj PhysicsObject) BroadPhaseOverlap(other PhysicsObject) bool {
	if obj.broadPhase == nil || other.broadPhase == nil {
		return false
	}
	return obj.broadPhase.Overlap(other.broadPhase)
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
