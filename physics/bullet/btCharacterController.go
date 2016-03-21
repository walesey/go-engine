package bullet

import (
	"github.com/luxengine/gobullet"
	vmath "github.com/walesey/go-engine/vectormath"
)

const CF_CHARACTER_OBJECT = 16

type BtCharacterController struct {
	kcc       gobullet.KinematicCharacterController
	ghost     gobullet.PairCachingGhostObject
	transform gobullet.Transform
}

func NewBtCharacterController(shape gobullet.CollisionShape, stepHeight float64) *BtCharacterController {
	ghost := gobullet.NewPairCachingGhostObject()
	ghost.SetCollisionShape(shape)
	kcc := gobullet.NewKinematicCharacterController(ghost, shape, float32(stepHeight))
	return &BtCharacterController{
		kcc:       kcc,
		ghost:     ghost,
		transform: gobullet.NewTransform(),
	}
}

func (cc BtCharacterController) Delete() {
	cc.kcc.Delete()
}

func (cc BtCharacterController) Warp(position vmath.Vector3) {
	setVector(cc.kcc.Warp, position)
}

func (cc BtCharacterController) Jump() {
	cc.kcc.Jump()
}

func (cc BtCharacterController) SetUpAxis(axis int) {
	cc.kcc.SetUpAxis(axis)
}

func (cc BtCharacterController) SetWalkDirection(dir vmath.Vector3) {
	setVector(cc.kcc.SetWalkDirection, dir)
}

func (cc BtCharacterController) SetVelocityForTimeInterval(speed vmath.Vector3, time float64) {
	//TODO
}

func (cc BtCharacterController) SetFallSpeed(speed float64) {
	cc.kcc.SetFallSpeed(float32(speed))
}

func (cc BtCharacterController) SetJumpSpeed(speed float64) {
	cc.kcc.SetJumpSpeed(float32(speed))
}

func (cc BtCharacterController) SetMaxJumpHeight(height float64) {
	cc.kcc.SetMaxJumpHeight(float32(height))
}

func (cc BtCharacterController) SetGravity(gravity float64) {
	cc.kcc.SetGravity(float32(gravity))
}

func (cc BtCharacterController) SetMaxSlope(radian float64) {
	cc.kcc.SetMaxSlope(float32(radian))
}

func (cc BtCharacterController) GetPosition() vmath.Vector3 {
	cc.ghost.GetWorldTransform(cc.transform)
	return getVector(cc.transform.GetOrigin)
}

func (cc BtCharacterController) CanJump() bool {
	return cc.kcc.CanJump()
}

func (cc BtCharacterController) GetGravity() float64 {
	return float64(cc.kcc.GetGravity())
}

func (cc BtCharacterController) GetMaxSlope() float64 {
	return float64(cc.kcc.GetMaxSlope())
}

func (cc BtCharacterController) OnGround() bool {
	return cc.kcc.OnGround()
}
