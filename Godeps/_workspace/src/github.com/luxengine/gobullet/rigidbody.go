package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

import (
	"reflect"
	"unsafe"
)

/*************************
********Rigid Body********
*************************/
type RigidBody struct {
	handle C.plRigidBody
}

//note: think about putting the userdata in the RigidBody struct and giving nil to bullet
//will check once callbacks are functional
func CreateRigidBody(userdata unsafe.Pointer, mass float32, shape CollisionShape) RigidBody {
	return RigidBody{C.plCreateRigidBody(userdata, C.float(mass), shape.handle)}
}

//overload for compound shape
func CreateRigidBodyCompound(userdata unsafe.Pointer, mass float32, shape CompoundShape) RigidBody {
	return RigidBody{C.plCreateRigidBody(userdata, C.float(mass), shape.handle)}
}

//overload for convex hull
func CreateRigidBodyConvex(userdata unsafe.Pointer, mass float32, shape ConvexHull) RigidBody {
	return RigidBody{C.plCreateRigidBody(userdata, C.float(mass), shape.handle)}
}

//overload for triangle mesh
func CreateRigidBodyConcave(userdata unsafe.Pointer, mass float32, shape TriangleMeshShape) RigidBody {
	return RigidBody{C.plCreateRigidBody(userdata, C.float(mass), shape.handle)}
}

var collisionshapetype = reflect.TypeOf(CollisionShape{})
var compoundshapetype = reflect.TypeOf(CompoundShape{})
var convexhulltype = reflect.TypeOf(ConvexHull{})
var trianglemeshshapetype = reflect.TypeOf(TriangleMeshShape{})

//func CreateRigidBodys(userdata *interface{}, mass float32, shape interface{}) {}

//Free the RigidBody
//NOTE: stepping a world with a non-removed, deleted RigidBody will crash.
func (this RigidBody) Delete() {
	C.plDeleteRigidBody(this.handle)
}

//UNTESTED
//Fill dest with the linear velocity of the rigid body (in meter/second)
func (this RigidBody) GetLinearVelocity(dest *[3]float32) {
	C.plRigidBodyGetLinearVelocity(this.handle, (*C.plReal)(&dest[0]))
}

//UNTESTED
func (this RigidBody) SetLinearVelocity(vel *[3]float32) {
	C.plRigidBodySetLinearVelocity(this.handle, (*C.plReal)(&vel[0]))
}

//UNTESTED
func (this RigidBody) SetLinearFactor(linearFactor *[3]float32) {
	C.plRigidBodySetLinearFactor(this.handle, (*C.plReal)(&linearFactor[0]))
}

//UNTESTED
func (this RigidBody) GetLinearFactor(linearFactor *[3]float32) {
	C.plRigidBodyGetLinearFactor(this.handle, (*C.plReal)(&linearFactor[0]))
}

//UNTESTED
func (this RigidBody) ApplyGravity() {
	C.plRigidBodyApplyGravity(this.handle)
}

//UNTESTED
func (this RigidBody) SetGravity(acceleration *[3]float32) {
	C.plRigidBodySetGravity(this.handle, (*C.plReal)(&acceleration[0]))
}

//UNTESTED
func (this RigidBody) GetGravity(acceleration *[3]float32) {
	C.plRigidBodyGetGravity(this.handle, (*C.plReal)(&acceleration[0]))
}

//UNTESTED
func (this RigidBody) SetDamping(linear, angular float32) {
	C.plRigidBodySetDamping(this.handle, C.plReal(linear), C.plReal(angular))
}

//UNTESTED
func (this RigidBody) GetLinearDamping() float32 {
	return float32(C.plRigidBodyGetLinearDamping(this.handle))
}

//UNTESTED
func (this RigidBody) GetAngularDamping() float32 {
	return float32(C.plRigidBodyGetAngularDamping(this.handle))
}

//UNTESTED
func (this RigidBody) SetRestitution(factor float32) {
	C.plRigidBodySetRestitution(this.handle, C.plReal(factor))
}

//UNTESTED
func (this RigidBody) GetRestitution() float32 {
	return float32(C.plRigidBodyGetRestitution(this.handle))
}

//UNTESTED
func (this RigidBody) SetFriction(factor float32) {
	C.plRigidBodySetFriction(this.handle, C.plReal(factor))
}

//UNTESTED
func (this RigidBody) GetFriction() float32 {
	return float32(C.plRigidBodyGetFriction(this.handle))
}

//UNTESTED
func (this RigidBody) SetRollingFriction(factor float32) {
	C.plRigidBodySetRollingFriction(this.handle, C.plReal(factor))
}

//UNTESTED
func (this RigidBody) GetRollingFriction() float32 {
	return float32(C.plRigidBodyGetRollingFriction(this.handle))
}

//Fills the provided [16]float32 with a matrix transform
//dest: the destination matrix
func (this RigidBody) GetOpenGLMatrix(dest *[16]float32) {
	C.plGetOpenGLMatrix(this.handle, (*C.plReal)(&dest[0]))
}

//Fills the provided *[3]float32 with the position of the RigidBody
//dest: the destination position
func (this RigidBody) GetPosition(dest *[3]float32) {
	C.plGetPosition(this.handle, (*C.plReal)(&dest[0]))
}

//Fills the provided *[3]float32 with the position of the RigidBody
//dest: the destination quaternion {w,{x,y,z}}
func (this RigidBody) GetOrientation(dest *[4]float32) {
	C.plGetOrientation(this.handle, (*C.plReal)(&dest[0]))
}

//Moves the RigidBody to that position
//position: {x,y,z}
func (this RigidBody) SetPosition(position *[3]float32) {
	C.plSetPosition(this.handle, (*C.plReal)(&position[0]))
}

//Sets the orientation of this RigidBody
//quaternion: {w,{x,y,z}}
func (this RigidBody) SetOrientation(quaternion *[4]float32) {
	C.plSetOrientation(this.handle, (*C.plReal)(&quaternion[0]))
}

//Sets the transform of this RigidBody
//mat: the transform matrix, column major
func (this RigidBody) SetOpenGLMatrix(mat *[16]float32) {
	C.plSetOpenGLMatrix(this.handle, (*C.plReal)(&mat[0]))
}

//UNTESTED
func (this RigidBody) SetAngularFactor(angfac *[3]float32) {
	C.plRigidBodySetAngularFactor(this.handle, (*C.plReal)(&angfac[0]))
}

//UNTESTED
func (this RigidBody) GetAngularFactor(angfac *[3]float32) {
	C.plRigidBodyGetAngularFactor(this.handle, (*C.plReal)(&angfac[0]))
}

//FORCES

//UNTESTED
func (this RigidBody) ApplyTorque(torque *[3]float32) {
	C.plRigidBodyApplyTorque(this.handle, (*C.plReal)(&torque[0]))
}

//UNTESTED
func (this RigidBody) ApplyForce(force, rel_pos *[3]float32) {
	C.plRigidBodyApplyForce(this.handle, (*C.plReal)(&force[0]), (*C.plReal)(&rel_pos[0]))
}

//UNTESTED
func (this RigidBody) ApplyCentralImpulse(impulse *[3]float32) {
	C.plRigidBodyApplyCentralImpulse(this.handle, (*C.plReal)(&impulse[0]))
}

//UNTESTED
func (this RigidBody) ApplyTorqueImpulse(torque *[3]float32) {
	C.plRigidBodyApplyTorqueImpulse(this.handle, (*C.plReal)(&torque[0]))
}

//UNTESTED
func (this RigidBody) ApplyImpulse(impulse, rel_pos *[3]float32) {
	C.plRigidBodyApplyImpulse(this.handle, (*C.plReal)(&impulse[0]), (*C.plReal)(&rel_pos[0]))
}

//UNTESTED
func (this RigidBody) ClearForces() {
	C.plRigidBodyClearForces(this.handle)
}

//UNTESTED
func (this RigidBody) UpdateInertiaTensor() {
	C.plRigidBodyUpdateInertiaTensor(this.handle)
}

//Collision Object methods

//UNTESTED
func (this RigidBody) GetActivationState() ActivationTag {
	return ActivationTag(C.plRigidBodyGetActivationState(this.handle))
}

//UNTESTED
func (this RigidBody) SetActivationState(newState ActivationTag) {
	C.plRigidBodySetActivationState(this.handle, C.int(newState))
}

//UNTESTED
func (this RigidBody) SetDeactivationTime(time float32) {
	C.plRigidBodySetDeactivationTime(this.handle, C.plReal(time))
}

//UNTESTED
func (this RigidBody) GetDeactivationTime() float32 {
	return float32(C.plRigidBodyGetDeactivationTime(this.handle))
}

//UNTESTED
func (this RigidBody) ForceActivationState(newState ActivationTag) {
	C.plRigidBodyForceActivationState(this.handle, C.int(newState))
}

//UNTESTED
func (this RigidBody) Activate(forceActivation bool) {
	C.plRigidBodyActivate(this.handle, C.bool(forceActivation))
}

//UNTESTED
func (this RigidBody) IsActive() bool {
	return bool(C.plRigidBodyIsActive(this.handle))
}

//UNTESTED
func (this RigidBody) SetUserPointer(ptr unsafe.Pointer) {
	C.plRigidBodySetUserPointer(this.handle, ptr)
}

//UNTESTED
func (this RigidBody) GetUserPointer() unsafe.Pointer {
	return C.plRigidBodyGetUserPointer(this.handle)
}
