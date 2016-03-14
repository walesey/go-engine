package dynamics

import (
	"math"

	"github.com/walesey/go-engine/physics/collision"
	vmath "github.com/walesey/go-engine/vectormath"
)

func sphereInertia(radius, mass float64) vmath.Matrix3 {
	return vmath.IdentityMatrix3().MultiplyScalar((2.0 / 5.0) * mass * radius * radius)
}

type PhysicsObjectImpl struct {
	position           vmath.Vector3
	velocity           vmath.Vector3
	totalForce         vmath.Vector3
	totalTorque        vmath.Vector3
	orientation        vmath.Quaternion
	angularVelocity    vmath.Quaternion
	inertiaTensor      vmath.Matrix3
	mass               float64
	radius             float64
	friction           float64
	restitution        float64
	activationVelocity float64 // > 0.0
	active             bool
	linearDamping      float64
	angularDamping     float64
	broadPhase         collision.Collider
	narrowPhase        collision.Collider
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
		position:           vmath.Vector3{0, 0, 0},
		velocity:           vmath.Vector3{0, 0, 0},
		orientation:        vmath.IdentityQuaternion(),
		angularVelocity:    vmath.Quaternion{1, 0, 0, 0},
		inertiaTensor:      vmath.IdentityMatrix3(),
		mass:               0.0,
		radius:             1.0,
		friction:           0.1,
		restitution:        0.1,
		activationVelocity: 0.01,
		active:             true,
		linearDamping:      0.01,
		angularDamping:     0.01,
	}
}

//NarrowPhaseOverlap
func (obj *PhysicsObjectImpl) NarrowPhaseOverlap(other *PhysicsObjectImpl) bool {
	if obj.narrowPhase == nil || other.narrowPhase == nil {
		return false
	}
	obj.narrowPhase.Offset(obj.position, obj.orientation)
	other.narrowPhase.Offset(other.position, other.orientation)
	return obj.narrowPhase.Overlap(other.narrowPhase)
}

//BroadPhaseOverlap
func (obj *PhysicsObjectImpl) BroadPhaseOverlap(other *PhysicsObjectImpl) bool {
	if obj.broadPhase == nil || other.broadPhase == nil {
		return false
	}
	obj.broadPhase.Offset(obj.position, obj.orientation)
	other.broadPhase.Offset(other.position, other.orientation)
	return obj.broadPhase.Overlap(other.broadPhase)
}

func (obj *PhysicsObjectImpl) PenetrationVector(other *PhysicsObjectImpl) vmath.Vector3 {
	if obj.narrowPhase == nil || other.narrowPhase == nil {
		return vmath.Vector3{}
	}
	obj.narrowPhase.Offset(obj.position, obj.orientation)
	other.narrowPhase.Offset(other.position, other.orientation)
	return obj.narrowPhase.PenetrationVector(other.narrowPhase)
}

func (obj *PhysicsObjectImpl) ContactPoint(other *PhysicsObjectImpl) vmath.Vector3 {
	if obj.narrowPhase == nil || other.narrowPhase == nil {
		return vmath.Vector3{}
	}
	obj.narrowPhase.Offset(obj.position, obj.orientation)
	other.narrowPhase.Offset(other.position, other.orientation)
	return obj.narrowPhase.ContactPoint(other.narrowPhase)
}

// AngularVelocityVector Get angular velocity as a Vector3
func (obj *PhysicsObjectImpl) GetAngularVelocityVector() vmath.Vector3 {
	w := vmath.Vector3{X: obj.angularVelocity.X, Y: obj.angularVelocity.Y, Z: obj.angularVelocity.Z}
	if !vmath.ApproxEqual(w.LengthSquared(), 1.0, 0.00001) {
		if w.LengthSquared() < 0.00001 {
			w.X = 1
		} else {
			w = w.Normalize()
			obj.angularVelocity.X, obj.angularVelocity.Y, obj.angularVelocity.Z = w.X, w.Y, w.Z
		}
	}
	return w.MultiplyScalar(obj.angularVelocity.W)
}

// SetAngularVelocityVector set angular velocity as a vector3
func (obj *PhysicsObjectImpl) SetAngularVelocityVector(av vmath.Vector3) {
	w := av.Length()
	obj.angularVelocity.W = w
	obj.angularVelocity.X = av.X
	obj.angularVelocity.Y = av.Y
	obj.angularVelocity.Z = av.Z
}

func (obj *PhysicsObjectImpl) SetRadius(radius float64) {
	obj.radius = radius
	if !obj.IsStatic() {
		obj.inertiaTensor = sphereInertia(obj.radius, obj.mass)
	}
}

func (obj *PhysicsObjectImpl) SetMass(mass float64) {
	obj.mass = mass
	if !obj.IsStatic() {
		obj.inertiaTensor = sphereInertia(obj.radius, obj.mass)
	}
}

// InertiaTensor returns the inertial Tensor of a sphere
func (obj *PhysicsObjectImpl) InertiaTensor() vmath.Matrix3 {
	mR := 0.4 * obj.mass * obj.radius
	return vmath.Matrix3{
		mR, 0.0, 0.0,
		0.0, mR, 0.0,
		0.0, 0.0, mR,
	}
}

// IsStatic object is static if it's mass is zero
func (obj *PhysicsObjectImpl) IsStatic() bool {
	return obj.mass <= 0.000000001
}

func (obj *PhysicsObjectImpl) IsActive() bool {
	return obj.active
}

// ApplyGravity Apply a force based on gravity and this object's mass
func (obj *PhysicsObjectImpl) ApplyGravity(gravity vmath.Vector3) {
	obj.totalForce = obj.totalForce.Add(gravity.MultiplyScalar(obj.mass))
}

// Apply a force on the object's centre of mass
func (obj *PhysicsObjectImpl) ApplyForce(force vmath.Vector3) {
	obj.totalForce = obj.totalForce.Add(force)
}

// Apply a torque on the object
func (obj *PhysicsObjectImpl) ApplyTorque(torque vmath.Vector3) {
	obj.totalTorque = obj.totalTorque.Add(torque)
}

// DoStep move the body assuming no constraints
func (obj *PhysicsObjectImpl) DoStep(timeStep float64) {

	//damping
	obj.velocity = obj.velocity.MultiplyScalar(math.Pow(1.0-obj.linearDamping, timeStep))
	obj.angularVelocity.W = obj.angularVelocity.W * math.Pow(1.0-obj.angularDamping, timeStep)

	//forces
	obj.velocity = obj.velocity.Add(obj.totalForce.DivideScalar(obj.mass).MultiplyScalar(timeStep))
	obj.totalForce = vmath.Vector3{0, 0, 0}

	// torques
	if obj.totalTorque.LengthSquared() > 0 {
		angularV := obj.GetAngularVelocityVector()
		angularV = angularV.Add(obj.inertiaTensor.Inverse().MultiplyVector(obj.totalTorque).MultiplyScalar(timeStep))
		obj.SetAngularVelocityVector(angularV)
		obj.totalTorque = vmath.Vector3{0, 0, 0}
	}

	//apply position increment
	obj.position = obj.position.Add(obj.velocity.MultiplyScalar(timeStep))

	//apply orientation increment
	axis := vmath.Vector3{obj.angularVelocity.X, obj.angularVelocity.Y, obj.angularVelocity.Z}
	obj.orientation = vmath.AngleAxis(timeStep*obj.angularVelocity.W, axis).Multiply(obj.orientation)

	// normalize orientation to prevent precision error
	if !vmath.ApproxEqual(obj.orientation.MagnitudeSquared(), 1.0, 0.1) {
		obj.orientation = obj.orientation.Normalize()
	}
}

func (obj *PhysicsObjectImpl) GetPosition() vmath.Vector3       { return obj.position }
func (obj *PhysicsObjectImpl) GetVelocity() vmath.Vector3       { return obj.velocity }
func (obj *PhysicsObjectImpl) GetOrientation() vmath.Quaternion { return obj.orientation }
func (obj *PhysicsObjectImpl) GetMass() float64                 { return obj.mass }
func (obj *PhysicsObjectImpl) GetRadius() float64               { return obj.radius }
func (obj *PhysicsObjectImpl) GetFriction() float64             { return obj.friction }
func (obj *PhysicsObjectImpl) GetRestitution() float64          { return obj.restitution }

func (obj *PhysicsObjectImpl) SetPosition(position vmath.Vector3) { obj.position = position }
func (obj *PhysicsObjectImpl) SetVelocity(velocity vmath.Vector3) { obj.velocity = velocity }
func (obj *PhysicsObjectImpl) SetOrientation(orientation vmath.Quaternion) {
	obj.orientation = orientation
}
func (obj *PhysicsObjectImpl) SetFriction(friction float64)       { obj.friction = friction }
func (obj *PhysicsObjectImpl) SetRestitution(restitution float64) { obj.restitution = restitution }
func (obj *PhysicsObjectImpl) SetBroadPhase(broadphase collision.Collider) {
	obj.broadPhase = broadphase
}
func (obj *PhysicsObjectImpl) SetNarrowPhase(narrowPhase collision.Collider) {
	obj.narrowPhase = narrowPhase
}
