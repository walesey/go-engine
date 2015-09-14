package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type ControllerManager struct {
	controllerList []Controller
}

func (c *ControllerManager) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	for i := range c.controllerList {
		c.controllerList[i].KeyCallback(window, key, scancode, action, mods)
	}
}

func (c *ControllerManager) AddController(newCont Controller) {
	c.controllerList = append(c.controllerList, newCont)
} 

func NewControllerManager(window *glfw.Window) *ControllerManager {
	var controllerList []Controller
	c := &ControllerManager{controllerList}
	window.SetKeyCallback(c.KeyCallback)
	return c
}