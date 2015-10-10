package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

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
	c := NewActionMap()
	x, y := 0.0, 0.0
	doLook := func(xpos, ypos float64) {
		actor.Look(xpos-x, ypos-y)
		x, y = xpos, ypos
	}
	c.BindAxisAction(doLook)
	c.BindAction(actor.StartMovingForward, glfw.KeyW, glfw.Press)
	c.BindAction(actor.StartMovingBackward, glfw.KeyS, glfw.Press)
	c.BindAction(actor.StopMovingForward, glfw.KeyW, glfw.Release)
	c.BindAction(actor.StopMovingBackward, glfw.KeyS, glfw.Release)
	c.BindAction(actor.StartStrafingLeft, glfw.KeyA, glfw.Press)
	c.BindAction(actor.StopStrafingLeft, glfw.KeyA, glfw.Release)
	c.BindAction(actor.StartStrafingRight, glfw.KeyD, glfw.Press)
	c.BindAction(actor.StopStrafingRight, glfw.KeyD, glfw.Release)
	c.BindAction(actor.Jump, glfw.KeySpace, glfw.Press)
	c.BindAction(actor.Crouch, glfw.KeyLeftControl, glfw.Press)
	c.BindAction(actor.StandUp, glfw.KeyLeftControl, glfw.Release)
	c.BindAction(actor.Prone, glfw.KeyZ, glfw.Press)
	c.BindAction(actor.StartSprinting, glfw.KeyLeftShift, glfw.Press)
	c.BindAction(actor.StopSprinting, glfw.KeyLeftShift, glfw.Release)
	return c
}
