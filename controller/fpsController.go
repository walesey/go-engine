package controller

type FPSActor interface {
	Look(dx, dy float32)
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
	var x, y float32 = 0.0, 0.0
	doLook := func(xpos, ypos float32) {
		actor.Look(xpos-x, ypos-y)
		x, y = xpos, ypos
	}
	c.BindAxisAction(doLook)
	c.BindKeyAction(actor.StartMovingForward, KeyW, Press)
	c.BindKeyAction(actor.StartMovingBackward, KeyS, Press)
	c.BindKeyAction(actor.StopMovingForward, KeyW, Release)
	c.BindKeyAction(actor.StopMovingBackward, KeyS, Release)
	c.BindKeyAction(actor.StartStrafingLeft, KeyA, Press)
	c.BindKeyAction(actor.StopStrafingLeft, KeyA, Release)
	c.BindKeyAction(actor.StartStrafingRight, KeyD, Press)
	c.BindKeyAction(actor.StopStrafingRight, KeyD, Release)
	c.BindKeyAction(actor.Jump, KeySpace, Press)
	c.BindKeyAction(actor.Crouch, KeyLeftControl, Press)
	c.BindKeyAction(actor.StandUp, KeyLeftControl, Release)
	c.BindKeyAction(actor.Prone, KeyZ, Press)
	c.BindKeyAction(actor.StartSprinting, KeyLeftShift, Press)
	c.BindKeyAction(actor.StopSprinting, KeyLeftShift, Release)
	return c
}
