package physicsAPI

import "github.com/go-gl/mathgl/mgl32"

type CharacterController interface {
	Delete()
	Warp(position mgl32.Vec3)
	Jump()
	SetWalkDirection(dir mgl32.Vec3)
	SetVelocityForTimeInterval(speed mgl32.Vec3, time float32)

	SetUpAxis(axis int)
	SetFallSpeed(speed float32)
	SetJumpSpeed(speed float32)
	SetMaxJumpHeight(height float32)
	SetGravity(gravity float32)
	SetMaxSlope(radian float32)

	GetPosition() mgl32.Vec3
	CanJump() bool
	GetGravity() float32
	GetMaxSlope() float32
	OnGround() bool
}
