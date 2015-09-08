package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type BasicMovementController struct {
	actor BasicMovementActor
	ActionMap map[KeyAction]func()
}

func (c *BasicMovementController) BindAction(function func(), key glfw.Key, action glfw.Action) {
	ka := KeyAction{key, action}
	c.ActionMap[ka] = function
}

func (c *BasicMovementController) TriggerAction(key glfw.Key, action glfw.Action) {
	ka := KeyAction{key, action}
	c.ActionMap[ka]()
}

func (c *BasicMovementController) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch {
		case key == glfw.KeyUp:
			if action == glfw.Press { c.actor.StartMovingUp() 
			} else if action == glfw.Release {c.actor.StopMovingUp() } 
		case key == glfw.KeyDown: 
			if action == glfw.Press { c.actor.StartMovingDown() 
			} else if action == glfw.Release {c.actor.StopMovingDown() } 
		case key == glfw.KeyLeft: 
			if action == glfw.Press { c.actor.StartMovingLeft() 
			} else if action == glfw.Release {c.actor.StopMovingLeft() } 
		case key == glfw.KeyRight: 
			if action == glfw.Press { c.actor.StartMovingRight() 
			} else if action == glfw.Release {c.actor.StopMovingRight() } 
		default: 
			c.TriggerAction(key, action) 
	}
}

func NewBasicMovementController(actor BasicMovementActor) *BasicMovementController {
	c := &BasicMovementController{actor, make(map[KeyAction]func())}
	return c
}