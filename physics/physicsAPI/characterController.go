package physicsAPI

import vmath "github.com/walesey/go-engine/vectormath"

type CharacterController interface {
	Delete()
	Warp(position vmath.Vector3)
	Jump()
	SetWalkDirection(dir vmath.Vector3)
	SetVelocityForTimeInterval(speed vmath.Vector3, time float64)

	SetUpAxis(axis int)
	SetFallSpeed(speed float64)
	SetJumpSpeed(speed float64)
	SetMaxJumpHeight(height float64)
	SetGravity(gravity float64)
	SetMaxSlope(radian float64)

	GetPosition() vmath.Vector3
	CanJump() bool
	GetGravity() float64
	GetMaxSlope() float64
	OnGround() bool
}
