package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type BasicMovementController struct {
	actor BasicMovementActor
	actionMap map[KeyAction]func()
}

func (c *BasicMovementController) BindAction(function func(), key glfw.Key, action glfw.Action) {
	ka := KeyAction{key, action}
	c.actionMap[ka] = function
}

func (c *BasicMovementController) TriggerAction(key glfw.Key, action glfw.Action) {
	ka := KeyAction{key, action}
	if c.actionMap[ka] != nil {
		c.actionMap[ka]()
	}
}

func (c *BasicMovementController) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		c.TriggerAction(key, action)
}

func NewBasicMovementController(actor BasicMovementActor) *BasicMovementController {
	c := &BasicMovementController{actor, make(map[KeyAction]func())}
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