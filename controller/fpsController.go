package controller

type FPSActor interface {
	Look(dx, dy float64)
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

func NewFPSController(actor FPSActor) Controller {
	c := CreateController()
	x, y := 0.0, 0.0
	doLook := func(xpos, ypos float64) {
		actor.Look(xpos-x, ypos-y)
		x, y = xpos, ypos
	}
	c.BindAxisAction(doLook)
	c.BindAction(actor.StartMovingForward, KeyW, Press)
	c.BindAction(actor.StartMovingBackward, KeyS, Press)
	c.BindAction(actor.StopMovingForward, KeyW, Release)
	c.BindAction(actor.StopMovingBackward, KeyS, Release)
	c.BindAction(actor.StartStrafingLeft, KeyA, Press)
	c.BindAction(actor.StopStrafingLeft, KeyA, Release)
	c.BindAction(actor.StartStrafingRight, KeyD, Press)
	c.BindAction(actor.StopStrafingRight, KeyD, Release)
	c.BindAction(actor.Jump, KeySpace, Press)
	c.BindAction(actor.Crouch, KeyLeftControl, Press)
	c.BindAction(actor.StandUp, KeyLeftControl, Release)
	c.BindAction(actor.Prone, KeyZ, Press)
	c.BindAction(actor.StartSprinting, KeyLeftShift, Press)
	c.BindAction(actor.StopSprinting, KeyLeftShift, Release)
	return c
}
