package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

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
	c := NewActionMap()
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
	c.BindMouseAction(mouseClick, glfw.MouseButtonLeft, glfw.Press)
	c.BindMouseAction(mouseRelease, glfw.MouseButtonLeft, glfw.Release)

	c.BindAction(actor.StartMovingUp, glfw.KeyUp, glfw.Press)
	c.BindAction(actor.StartMovingDown, glfw.KeyDown, glfw.Press)
	c.BindAction(actor.StartMovingLeft, glfw.KeyLeft, glfw.Press)
	c.BindAction(actor.StartMovingRight, glfw.KeyRight, glfw.Press)
	c.BindAction(actor.StopMovingUp, glfw.KeyUp, glfw.Release)
	c.BindAction(actor.StopMovingDown, glfw.KeyDown, glfw.Release)
	c.BindAction(actor.StopMovingLeft, glfw.KeyLeft, glfw.Release)
	c.BindAction(actor.StopMovingRight, glfw.KeyRight, glfw.Release)

	c.BindAction(actor.StartMovingUp, glfw.KeyW, glfw.Press)
	c.BindAction(actor.StartMovingDown, glfw.KeyS, glfw.Press)
	c.BindAction(actor.StartMovingLeft, glfw.KeyA, glfw.Press)
	c.BindAction(actor.StartMovingRight, glfw.KeyD, glfw.Press)
	c.BindAction(actor.StopMovingUp, glfw.KeyW, glfw.Release)
	c.BindAction(actor.StopMovingDown, glfw.KeyS, glfw.Release)
	c.BindAction(actor.StopMovingLeft, glfw.KeyA, glfw.Release)
	c.BindAction(actor.StopMovingRight, glfw.KeyD, glfw.Release)

	return c
}
