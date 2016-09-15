package physicsAPI

import (
	"github.com/go-gl/mathgl/mgl32"
)

type PhysicsObject interface {
	Delete()
	ApplyForce(force, position mgl32.Vec3)
	ApplyTorque(torque mgl32.Vec3)

	GetPosition() mgl32.Vec3
	GetVelocity() mgl32.Vec3
	GetOrientation() mgl32.Quat
	GetAngularVelocityVector() mgl32.Vec3
	GetMass() float32
	GetRadius() float32
	GetFriction() float32
	GetRestitution() float32
	InertiaTensor() mgl32.Mat3
	IsStatic() bool

	SetPosition(position mgl32.Vec3)
	SetVelocity(velocity mgl32.Vec3)
	SetOrientation(orientation mgl32.Quat)
	SetAngularVelocityVector(av mgl32.Vec3)
	SetMass(mass float32)
	SetRadius(radius float32)
	SetFriction(friction float32)
	SetRestitution(restitution float32)
}

type PhysicsObject2D interface {
	KineticEnergy() float32
	SetMass(mass float32)
	SetMoment(moment float32)
	GetMoment() float32
	SetAngle(angle float32)
	AddAngle(angle float32)
	GetMass() float32
	SetPosition(pos mgl32.Vec2)
	AddForce(force mgl32.Vec2)
	SetForce(force mgl32.Vec2)
	AddVelocity(velocity mgl32.Vec2)
	SetVelocity(velocity mgl32.Vec2)
	AddTorque(t float32)
	GetTorque() float32
	GetAngularVelocity() float32
	SetTorque(t float32)
	AddAngularVelocity(w float32)
	SetAngularVelocity(w float32)
	GetVelocity() mgl32.Vec2
	GetPosition() mgl32.Vec2
	GetAngle() float32
}
