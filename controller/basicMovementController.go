package controller

type BasicMovementActor interface {
	Look(dx, dy float64)
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
	x, y := 0.0, 0.0
	drag := false
	doLook := func(xpos, ypos float64) {
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

	c.BindAction(actor.StartMovingUp, KeyUp, Press)
	c.BindAction(actor.StartMovingDown, KeyDown, Press)
	c.BindAction(actor.StartMovingLeft, KeyLeft, Press)
	c.BindAction(actor.StartMovingRight, KeyRight, Press)
	c.BindAction(actor.StopMovingUp, KeyUp, Release)
	c.BindAction(actor.StopMovingDown, KeyDown, Release)
	c.BindAction(actor.StopMovingLeft, KeyLeft, Release)
	c.BindAction(actor.StopMovingRight, KeyRight, Release)

	c.BindAction(actor.StartMovingUp, KeyW, Press)
	c.BindAction(actor.StartMovingDown, KeyS, Press)
	c.BindAction(actor.StartMovingLeft, KeyA, Press)
	c.BindAction(actor.StartMovingRight, KeyD, Press)
	c.BindAction(actor.StopMovingUp, KeyW, Release)
	c.BindAction(actor.StopMovingDown, KeyS, Release)
	c.BindAction(actor.StopMovingLeft, KeyA, Release)
	c.BindAction(actor.StopMovingRight, KeyD, Release)
	return c
}
