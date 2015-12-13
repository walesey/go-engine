package physicsObject

import vmath "github.com/walesey/go-engine/vectormath"

type PhysicsObject struct {
	Position, Velocity           vmath.Vector3
	Orientation, AngularVelocity vmath.Quaternion
	Mass, Radius                 float64
	Friction                     float64 // (0.0 to 1.0)
	Static                       bool    //disables movement
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
	return newPhysicsObject()
}

func (objPool PhysicsObjectPool) ReleasePhysicsObject(obj *PhysicsObject) {
	objPool.pool = append(objPool.pool, obj)
}

func newPhysicsObject() *PhysicsObject {
	return &PhysicsObject{
		Position:        vmath.Vector3{0, 0, 0},
		Velocity:        vmath.Vector3{0, 0, 0},
		Orientation:     vmath.IdentityQuaternion(),
		AngularVelocity: vmath.Quaternion{1, 0, 0, 0},
		Mass:            1.0,
		Radius:          1.0,
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
	obj.NarrowPhase.Offset(obj.Position, obj.Orientation)
	other.NarrowPhase.Offset(other.Position, other.Orientation)
	return obj.NarrowPhase.PenetrationVector(other.NarrowPhase)
}

func (obj *PhysicsObject) ContactPoint(other *PhysicsObject) vmath.Vector3 {
	if obj.NarrowPhase == nil || other.NarrowPhase == nil {
		return vmath.Vector3{}
	}
	obj.NarrowPhase.Offset(obj.Position, obj.Orientation)
	other.NarrowPhase.Offset(other.Position, other.Orientation)
	return obj.NarrowPhase.ContactPoint(other.NarrowPhase)
}

// AngularVelocityVector Get angular velocity as a Vector3
func (obj *PhysicsObject) AngularVelocityVector() vmath.Vector3 {
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
func (obj *PhysicsObject) SetAngularVelocityVector(av vmath.Vector3) {
	w := av.Length()
	obj.AngularVelocity.W = w
	obj.AngularVelocity.X = av.X
	obj.AngularVelocity.Y = av.Y
	obj.AngularVelocity.Z = av.Z
}

func (obj *PhysicsObject) doStep(dt float64) {
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
