package controller 

type FPSActor interface {
	StartMovingForward()
	StartMovingBackward()
	StartStrafingLeft()
	StartStrafingRight()
	StopMovingForward()
	StopMovingBackward()
	StopStrafingLeft()
	StopStrafingRight()
	Jump()
	StandUp()
	Crouch()
	Prone()
	StartSprinting()
	StopSprinting()
}