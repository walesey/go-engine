package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type ControllerManager struct {
	controllerList []Controller
}

//Key Callback
func (c *ControllerManager) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	for i := range c.controllerList {
		c.controllerList[i].KeyCallback(window, key, scancode, action, mods)
	}
}

//Mouse movement callback
func (c *ControllerManager) CursorPosCallback(window *glfw.Window, xpos, ypos float64) {
	for i := range c.controllerList {
		c.controllerList[i].CursorPosCallback(window, xpos, ypos)
	}
}

//Mouse scrollwheel callback
func (c *ControllerManager) ScrollCallback(window *glfw.Window, xoffset, yoffset float64) {
	for i := range c.controllerList {
		c.controllerList[i].ScrollCallback(window, xoffset, yoffset)
	}
}

func (c *ControllerManager) AddController(newCont Controller) {
	c.controllerList = append(c.controllerList, newCont)
}

func NewControllerManager(window *glfw.Window) *ControllerManager {
	var controllerList []Controller
	c := &ControllerManager{controllerList}
	window.SetKeyCallback(c.KeyCallback)
	window.SetCursorPosCallback(c.CursorPosCallback)
	window.SetScrollCallback(c.ScrollCallback)
	return c
}
