package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type BasicMovementActor interface {
	StartMovingUp()
	StartMovingDown()
	StartMovingRight()
	StartMovingLeft()
	StopMovingUp()
	StopMovingDown()
	StopMovingRight()
	StopMovingLeft()
}

func NewBasicMovementController(actor BasicMovementActor) Controller {
	c := NewActionMap()
	c.BindAction(actor.StartMovingUp, glfw.KeyUp, glfw.Press)
	c.BindAction(actor.StartMovingDown, glfw.KeyDown, glfw.Press)
	c.BindAction(actor.StartMovingLeft, glfw.KeyLeft, glfw.Press)
	c.BindAction(actor.StartMovingRight, glfw.KeyRight, glfw.Press)
	c.BindAction(actor.StopMovingUp, glfw.KeyUp, glfw.Release)
	c.BindAction(actor.StopMovingDown, glfw.KeyDown, glfw.Release)
	c.BindAction(actor.StopMovingLeft, glfw.KeyLeft, glfw.Release)
	c.BindAction(actor.StopMovingRight, glfw.KeyRight, glfw.Release)
	return c
}
