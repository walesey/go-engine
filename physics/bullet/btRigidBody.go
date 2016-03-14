package bullet

import (
	"github.com/luxengine/gobullet"
	"github.com/walesey/go-engine/physics/collision"
	"github.com/walesey/go-engine/physics/physicsAPI"
	vmath "github.com/walesey/go-engine/vectormath"
)

type BtRigidBody struct {
	rigidBody gobullet.RigidBody
}

func NewBtRigidBody(mass float64, shape gobullet.CollisionShape) physicsAPI.PhysicsObject {
	return &BtRigidBody{gobullet.CreateRigidBody(nil, float32(mass), shape)}
}

func NewBtRigidBodyCompound(mass float64, shape gobullet.CompoundShape) physicsAPI.PhysicsObject {
	return &BtRigidBody{gobullet.CreateRigidBodyCompound(nil, float32(mass), shape)}
}

func NewBtRigidBodyConvex(mass float64, shape gobullet.ConvexHull) physicsAPI.PhysicsObject {
	return &BtRigidBody{gobullet.CreateRigidBodyConvex(nil, float32(mass), shape)}
}

func NewBtRigidBodyConcave(mass float64, shape gobullet.TriangleMeshShape) physicsAPI.PhysicsObject {
	return &BtRigidBody{gobullet.CreateRigidBodyConcave(nil, float32(mass), shape)}
}

func NewBtRigidBodyGImpact(mass float64, shape gobullet.GImpactMeshShape) physicsAPI.PhysicsObject {
	return &BtRigidBody{shape.RigidBody(nil, float32(mass))}
}

func (body *BtRigidBody) DoStep(timeStep float64)            {}
func (body *BtRigidBody) ApplyGravity(gravity vmath.Vector3) {}
func (body *BtRigidBody) ApplyForce(force vmath.Vector3)     {}
func (body *BtRigidBody) ApplyTorque(torque vmath.Vector3)   {}

func (body *BtRigidBody) GetPosition() vmath.Vector3 {
	return getVector(body.rigidBody.GetPosition)
}

func (body *BtRigidBody) GetVelocity() vmath.Vector3 {
	return getVector(body.rigidBody.GetLinearVelocity)
}

func (body *BtRigidBody) GetOrientation() vmath.Quaternion {
	return getQuaternion(body.rigidBody.GetOrientation)
}

func (body *BtRigidBody) GetAngularVelocityVector() vmath.Vector3 { return vmath.Vector3{} }
func (body *BtRigidBody) GetMass() float64                        { return 0 }
func (body *BtRigidBody) GetRadius() float64                      { return 0 }
func (body *BtRigidBody) GetFriction() float64                    { return 0 }
func (body *BtRigidBody) GetRestitution() float64                 { return 0 }
func (body *BtRigidBody) InertiaTensor() vmath.Matrix3            { return vmath.Matrix3{} }
func (body *BtRigidBody) IsStatic() bool                          { return false }

func (body *BtRigidBody) SetPosition(position vmath.Vector3) {
	setVector(body.rigidBody.SetPosition, position)
}

func (body *BtRigidBody) SetVelocity(velocity vmath.Vector3) {
	setVector(body.rigidBody.SetLinearVelocity, velocity)
}

func (body *BtRigidBody) SetOrientation(orientation vmath.Quaternion) {
	setQuaternion(body.rigidBody.SetOrientation, orientation)
}

func (body *BtRigidBody) SetAngularVelocityVector(av vmath.Vector3)    {}
func (body *BtRigidBody) SetMass(mass float64)                         {}
func (body *BtRigidBody) SetRadius(radius float64)                     {}
func (body *BtRigidBody) SetFriction(friction float64)                 {}
func (body *BtRigidBody) SetRestitution(restitution float64)           {}
func (body *BtRigidBody) SetBroadPhase(broadphase collision.Collider)  {}
func (body *BtRigidBody) SetNarrowPhase(broadphase collision.Collider) {}
