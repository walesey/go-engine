package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type FPSController struct {
	actor FPSActor
	actionMap map[KeyAction]func()
}

func (c *FPSController) BindAction(function func(), key glfw.Key, action glfw.Action) {
	ka := KeyAction{key, action}
	c.actionMap[ka] = function
}

func (c *FPSController) TriggerAction(key glfw.Key, action glfw.Action) {
	ka := KeyAction{key, action}
	if c.actionMap[ka] != nil {
		c.actionMap[ka]()
	}
}

func (c *FPSController) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		c.TriggerAction(key, action)
}

func NewFPSController(actor FPSActor) *FPSController {
	c := &FPSController{actor, make(map[KeyAction]func())}
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