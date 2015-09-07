package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Controller struct {
	ActionMap map[glfw.Key]func()
	Window *glfw.Window 
}

func (c *Controller) BindAction(action func(), key glfw.Key) {
	c.ActionMap[key] = action
}

func (c *Controller) TriggerAction(key glfw.Key) {
	c.ActionMap[key]()
}

func (c *Controller) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		c.TriggerAction(key)
	}
}

func Poll() {
	glfw.PollEvents()
}

func NewController(window *glfw.Window) *Controller {
	var c = &Controller{make(map[glfw.Key]func()), window}
	c.Window.SetKeyCallback(c.KeyCallback)
	return c
}