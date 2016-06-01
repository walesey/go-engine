package physicsAPI

import vmath "github.com/walesey/go-engine/vectormath"

type PhysicsObject interface {
	Delete()
	ApplyForce(force, position vmath.Vector3)
	ApplyTorque(torque vmath.Vector3)

	GetPosition() vmath.Vector3
	GetVelocity() vmath.Vector3
	GetOrientation() vmath.Quaternion
	GetAngularVelocityVector() vmath.Vector3
	GetMass() float64
	GetRadius() float64
	GetFriction() float64
	GetRestitution() float64
	InertiaTensor() vmath.Matrix3
	IsStatic() bool

	SetPosition(position vmath.Vector3)
	SetVelocity(velocity vmath.Vector3)
	SetOrientation(orientation vmath.Quaternion)
	SetAngularVelocityVector(av vmath.Vector3)
	SetMass(mass float64)
	SetRadius(radius float64)
	SetFriction(friction float64)
	SetRestitution(restitution float64)
}

type PhysicsObject2D interface {
	KineticEnergy() float64
	SetMass(mass float64)
	SetMoment(moment float64)
	GetMoment() float64
	SetAngle(angle float64)
	AddAngle(angle float64)
	GetMass() float64
	SetPosition(pos vmath.Vector2)
	AddForce(force vmath.Vector2)
	SetForce(force vmath.Vector2)
	AddVelocity(velocity vmath.Vector2)
	SetVelocity(velocity vmath.Vector2)
	AddTorque(t float64)
	GetTorque() float64
	GetAngularVelocity() float64
	SetTorque(t float64)
	AddAngularVelocity(w float64)
	SetAngularVelocity(w float64)
	GetVelocity() vmath.Vector2
	GetPosition() vmath.Vector2
	GetAngle() float64
}
