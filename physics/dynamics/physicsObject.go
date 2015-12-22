package dynamics

import (
	"github.com/walesey/go-engine/physics/collision"
	vmath "github.com/walesey/go-engine/vectormath"
)

type PhysicsObject interface {
	GetPosition() vmath.Vector3
	GetVelocity() vmath.Vector3
	GetOrientation() vmath.Quaternion
	GetAngularVelocityVector() vmath.Vector3
	GetMass() float64
	GetRadius() float64
	GetFriction() float64

	SetPosition(position vmath.Vector3)
	SetVelocity(velocity vmath.Vector3)
	SetOrientation(orientation vmath.Quaternion)
	SetAngularVelocityVector(av vmath.Vector3)
	SetMass(mass float64)
	SetRadius(radius float64)
	SetFriction(friction float64)
	SetStatic(static bool)
	SetBroadPhase(broadphase collision.Collider)
	SetNarrowPhase(broadphase collision.Collider)
}

type PhysicsObjectImpl struct {
	Position, Velocity           vmath.Vector3
	Orientation, AngularVelocity vmath.Quaternion
	Mass, Radius                 float64
	Friction                     float64 // (0.0 to 1.0)
	ActiveVelocity               float64 // > 0.0
	Static, Active               bool    //disables movement
	ForceStore                   *ForceStore
	BroadPhase, NarrowPhase      collision.Collider
}

type PhysicsObjectPool struct {
	pool []*PhysicsObjectImpl
}

func NewPhysicsObjectPool() *PhysicsObjectPool {
	return &PhysicsObjectPool{
		pool: make([]*PhysicsObjectImpl, 0, 0),
	}
}

func (objPool *PhysicsObjectPool) GetPhysicsObject() *PhysicsObjectImpl {
	if len(objPool.pool) > 0 {
		obj := objPool.pool[len(objPool.pool)-1]
		objPool.pool = objPool.pool[:len(objPool.pool)-1]
		return obj
	}
	return newPhysicsObject()
}

func (objPool *PhysicsObjectPool) ReleasePhysicsObject(obj *PhysicsObjectImpl) {
	objPool.pool = append(objPool.pool, obj)
}

func newPhysicsObject() *PhysicsObjectImpl {
	return &PhysicsObjectImpl{
		Position:        vmath.Vector3{0, 0, 0},
		Velocity:        vmath.Vector3{0, 0, 0},
		Orientation:     vmath.IdentityQuaternion(),
		AngularVelocity: vmath.Quaternion{1, 0, 0, 0},
		Mass:            1.0,
		Radius:          1.0,
		ActiveVelocity:  1.0,
		Active:          true,
		ForceStore:      NewForceStore(),
	}
}

//NarrowPhaseOverlap
func (obj *PhysicsObjectImpl) NarrowPhaseOverlap(other *PhysicsObjectImpl) bool {
	if obj.NarrowPhase == nil || other.NarrowPhase == nil {
		return false
	}
	obj.NarrowPhase.Offset(obj.Position, obj.Orientation)
	other.NarrowPhase.Offset(other.Position, other.Orientation)
	return obj.NarrowPhase.Overlap(other.NarrowPhase)
}

//BroadPhaseOverlap
func (obj *PhysicsObjectImpl) BroadPhaseOverlap(other *PhysicsObjectImpl) bool {
	if obj.BroadPhase == nil || other.BroadPhase == nil {
		return false
	}
	obj.BroadPhase.Offset(obj.Position, obj.Orientation)
	other.BroadPhase.Offset(other.Position, other.Orientation)
	return obj.BroadPhase.Overlap(other.BroadPhase)
}

func (obj *PhysicsObjectImpl) PenetrationVector(other *PhysicsObjectImpl) vmath.Vector3 {
	if obj.NarrowPhase == nil || other.NarrowPhase == nil {
		return vmath.Vector3{}
	}
	obj.NarrowPhase.Offset(obj.Position, obj.Orientation)
	other.NarrowPhase.Offset(other.Position, other.Orientation)
	return obj.NarrowPhase.PenetrationVector(other.NarrowPhase)
}

func (obj *PhysicsObjectImpl) ContactPoint(other *PhysicsObjectImpl) vmath.Vector3 {
	if obj.NarrowPhase == nil || other.NarrowPhase == nil {
		return vmath.Vector3{}
	}
	obj.NarrowPhase.Offset(obj.Position, obj.Orientation)
	other.NarrowPhase.Offset(other.Position, other.Orientation)
	return obj.NarrowPhase.ContactPoint(other.NarrowPhase)
}

// AngularVelocityVector Get angular velocity as a Vector3
func (obj *PhysicsObjectImpl) GetAngularVelocityVector() vmath.Vector3 {
	w := vmath.Vector3{X: obj.AngularVelocity.X, Y: obj.AngularVelocity.Y, Z: obj.AngularVelocity.Z}
	if !vmath.ApproxEqual(w.LengthSquared(), 1.0, 0.00001) {
		if w.LengthSquared() < 0.00001 {
			w.X = 1
		} else {
			w = w.Normalize()
		}
	}
	return w.MultiplyScalar(obj.AngularVelocity.W)
}

// SetAngularVelocityVector set angular velocity as a vector3
func (obj *PhysicsObjectImpl) SetAngularVelocityVector(av vmath.Vector3) {
	w := av.Length()
	obj.AngularVelocity.W = w
	obj.AngularVelocity.X = av.X
	obj.AngularVelocity.Y = av.Y
	obj.AngularVelocity.Z = av.Z
}

// InertiaTensor returns the inertial Tensor of a sphere
func (obj *PhysicsObjectImpl) InertiaTensor() vmath.Matrix3 {
	mR := 0.4 * obj.Mass * obj.Radius
	return vmath.Matrix3{
		mR, 0.0, 0.0,
		0.0, mR, 0.0,
		0.0, 0.0, mR,
	}
}

func (obj *PhysicsObjectImpl) DoStep(dt float64) {

	//process forces and acceleration
	obj.ForceStore.DoStep(dt, obj)

	//apply position increment
	obj.Position = obj.Position.Add(obj.Velocity.MultiplyScalar(dt))

	//apply orientation increment
	axis := vmath.Vector3{obj.AngularVelocity.X, obj.AngularVelocity.Y, obj.AngularVelocity.Z}
	obj.Orientation = vmath.AngleAxis(dt*obj.AngularVelocity.W, axis).Multiply(obj.Orientation)

	// normalize orientation to prevent precision error
	if !vmath.ApproxEqual(obj.Orientation.MagnitudeSquared(), 1.0, 0.1) {
		obj.Orientation = obj.Orientation.Normalize()
	}
}

func (obj *PhysicsObjectImpl) GetPosition() vmath.Vector3       { return obj.Position }
func (obj *PhysicsObjectImpl) GetVelocity() vmath.Vector3       { return obj.Velocity }
func (obj *PhysicsObjectImpl) GetOrientation() vmath.Quaternion { return obj.Orientation }
func (obj *PhysicsObjectImpl) GetMass() float64                 { return obj.Mass }
func (obj *PhysicsObjectImpl) GetRadius() float64               { return obj.Radius }
func (obj *PhysicsObjectImpl) GetFriction() float64             { return obj.Friction }

func (obj *PhysicsObjectImpl) SetPosition(position vmath.Vector3) { obj.Position = position }
func (obj *PhysicsObjectImpl) SetVelocity(velocity vmath.Vector3) { obj.Velocity = velocity }
func (obj *PhysicsObjectImpl) SetOrientation(orientation vmath.Quaternion) {
	obj.Orientation = orientation
}
func (obj *PhysicsObjectImpl) SetMass(mass float64)         { obj.Mass = mass }
func (obj *PhysicsObjectImpl) SetRadius(radius float64)     { obj.Radius = radius }
func (obj *PhysicsObjectImpl) SetFriction(friction float64) { obj.Friction = friction }
func (obj *PhysicsObjectImpl) SetStatic(static bool)        { obj.Static = static }
func (obj *PhysicsObjectImpl) SetBroadPhase(broadphase collision.Collider) {
	obj.BroadPhase = broadphase
}
func (obj *PhysicsObjectImpl) SetNarrowPhase(narrowPhase collision.Collider) {
	obj.NarrowPhase = narrowPhase
}
