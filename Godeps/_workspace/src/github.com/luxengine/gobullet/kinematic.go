package gobullet

/*
#include <bullet/Bullet-C-Api.h>
#cgo LDFLAGS: -lBulletSoftBody -lBulletCollision -lBulletDynamics -lLinearMath
*/
import "C"

/*********************************************************
***Pair Caching Ghost Object (for kinematic controller)***
*********************************************************/
type PairCachingGhostObject struct {
	handle C.plPairCachingGhostObject
}

//UNTESTED
func NewPairCachingGhostObject() PairCachingGhostObject {
	return PairCachingGhostObject{C.plNewPairCachingGhostObject()}
}

//UNTESTED
func (this PairCachingGhostObject) Delete() {
	C.plPairCachingGhostObjectDelete(this.handle)
}

//UNTESTED
func (this PairCachingGhostObject) SetWorldTransform(t Transform) {
	C.plPairCachingGhostObjectSetWorldTransform(this.handle, t.handle)
}

//UNTESTED
func (this PairCachingGhostObject) SetCollisionShape(shape CollisionShape) {
	C.plPairCachingGhostObjectSetCollisionShape(this.handle, shape.handle)
}

//UNTESTED
func (this PairCachingGhostObject) SetCollisionFlags(flags int16) {
	C.plPairCachingGhostObjectSetCollisionFlags(this.handle, C.short(flags))
}

//UNTESTED
func (this PairCachingGhostObject) GetWorldTransform(transform Transform) {
	C.plPairCachingGhostObjectGetWorldTransform(this.handle, transform.handle)
}

/*************************************
****Kinematic Character Controller****
*************************************/
type KinematicCharacterController struct {
	handle C.plKinematicCharacterController
}

//UNTESTED
func NewKinematicCharacterController(ghost PairCachingGhostObject, shape CollisionShape, stepHeight float32) KinematicCharacterController {
	return KinematicCharacterController{C.plNewKinematicCharacterController(ghost.handle, shape.handle, C.plReal(stepHeight))}
}

//UNTESTED
func (this KinematicCharacterController) Delete() {
	C.plKinematicCharacterControllerDelete(this.handle)
}

//UNTESTED
func (this KinematicCharacterController) SetUpAxis(axis int) {
	C.plKinematicCharacterControllerSetUpAxis(this.handle, C.int(axis))
}

//UNTESTED
func (this KinematicCharacterController) SetWalkDirection(dir *[3]float32) {
	C.plKinematicCharacterControllerSetWalkDirection(this.handle, (*C.plReal)(&dir[0]))
}

//UNTESTED
func (this KinematicCharacterController) SetVelocityForTimeInterval(speed *[3]float32, time float32) {
	C.plKinematicCharacterControllerSetVelocityForTimeInterval(this.handle, (*C.plReal)(&speed[0]), C.plReal(time))
}

//UNTESTED
func (this KinematicCharacterController) Warp(position *[3]float32) {
	C.plKinematicCharacterControllerWarp(this.handle, (*C.plReal)(&position[0]))
}

//UNTESTED
func (this KinematicCharacterController) SetFallSpeed(speed float32) {
	C.plKinematicCharacterControllerSetFallSpeed(this.handle, C.plReal(speed))
}

//UNTESTED
func (this KinematicCharacterController) SetJumpSpeed(speed float32) {
	C.plKinematicCharacterControllerSetJumpSpeed(this.handle, C.plReal(speed))
}

//UNTESTED
func (this KinematicCharacterController) SetMaxJumpHeight(height float32) {
	C.plKinematicCharacterControllerSetMaxJumpHeight(this.handle, C.plReal(height))
}

//UNTESTED
func (this KinematicCharacterController) CanJump() bool {
	return bool(C.plKinematicCharacterControllerCanJump(this.handle))
}

//UNTESTED
func (this KinematicCharacterController) Jump() {
	C.plKinematicCharacterControllerJump(this.handle)
}

//UNTESTED
func (this KinematicCharacterController) SetGravity(gravity float32) {
	C.plKinematicCharacterControllerSetGravity(this.handle, C.plReal(gravity))
}

//UNTESTED
func (this KinematicCharacterController) GetGravity() float32 {
	return float32(C.plKinematicCharacterControllerGetGravity(this.handle))
}

//UNTESTED
func (this KinematicCharacterController) SetMaxSlope(radian float32) {
	C.plKinematicCharacterControllerSetMaxSlope(this.handle, C.plReal(radian))
}

//UNTESTED
func (this KinematicCharacterController) GetMaxSlope() float32 {
	return float32(C.plKinematicCharacterControllerGetMaxSlope(this.handle))
}

//UNTESTED
func (this KinematicCharacterController) OnGround() bool {
	return bool(C.plKinematicCharacterControllerOnGround(this.handle))
}
