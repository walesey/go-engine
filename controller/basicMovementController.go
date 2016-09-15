package controller

type BasicMovementActor interface {
	Look(dx, dy float32)
	StartMovingUp()
	StartMovingDown()
	StartMovingRight()
	StartMovingLeft()
	StopMovingUp()
	StopMovingDown()
	StopMovingRight()
	StopMovingLeft()
}

func NewBasicMovementController(actor BasicMovementActor, useDrag bool) Controller {
	c := CreateController()
	var x, y float32 = 0.0, 0.0
	drag := false
	doLook := func(xpos, ypos float32) {
		if drag || !useDrag {
			actor.Look(xpos-x, ypos-y)
		}
		x, y = xpos, ypos
	}
	mouseClick := func() {
		drag = true
	}
	mouseRelease := func() {
		drag = false
	}
	c.BindAxisAction(doLook)
	c.BindMouseAction(mouseClick, MouseButtonLeft, Press)
	c.BindMouseAction(mouseRelease, MouseButtonLeft, Release)

	c.BindKeyAction(actor.StartMovingUp, KeyUp, Press)
	c.BindKeyAction(actor.StartMovingDown, KeyDown, Press)
	c.BindKeyAction(actor.StartMovingLeft, KeyLeft, Press)
	c.BindKeyAction(actor.StartMovingRight, KeyRight, Press)
	c.BindKeyAction(actor.StopMovingUp, KeyUp, Release)
	c.BindKeyAction(actor.StopMovingDown, KeyDown, Release)
	c.BindKeyAction(actor.StopMovingLeft, KeyLeft, Release)
	c.BindKeyAction(actor.StopMovingRight, KeyRight, Release)

	c.BindKeyAction(actor.StartMovingUp, KeyW, Press)
	c.BindKeyAction(actor.StartMovingDown, KeyS, Press)
	c.BindKeyAction(actor.StartMovingLeft, KeyA, Press)
	c.BindKeyAction(actor.StartMovingRight, KeyD, Press)
	c.BindKeyAction(actor.StopMovingUp, KeyW, Release)
	c.BindKeyAction(actor.StopMovingDown, KeyS, Release)
	c.BindKeyAction(actor.StopMovingLeft, KeyA, Release)
	c.BindKeyAction(actor.StopMovingRight, KeyD, Release)
	return c
}
