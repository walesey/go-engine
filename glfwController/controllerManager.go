package glfwController

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type ControllerManager struct {
	controllerList []Controller
}

//Key Callback
func (c *ControllerManager) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	for _, cont := range c.controllerList {
		cont.KeyCallback(window, key, scancode, action, mods)
	}
}

//Mouse click callback
func (c *ControllerManager) MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	for _, cont := range c.controllerList {
		cont.MouseButtonCallback(window, button, action, mods)
	}
}

//Mouse movement callback
func (c *ControllerManager) CursorPosCallback(window *glfw.Window, xpos, ypos float64) {
	for _, cont := range c.controllerList {
		cont.CursorPosCallback(window, float32(xpos), float32(ypos))
	}
}

//Mouse scrollwheel callback
func (c *ControllerManager) ScrollCallback(window *glfw.Window, xoffset, yoffset float64) {
	for _, cont := range c.controllerList {
		cont.ScrollCallback(window, float32(xoffset), float32(yoffset))
	}
}

func (c *ControllerManager) AddController(newCont Controller) {
	c.controllerList = append(c.controllerList, newCont)
}

func (c *ControllerManager) RemoveController(controller Controller) {
	for index, cont := range c.controllerList {
		if cont == controller {
			c.controllerList = append(c.controllerList[:index], c.controllerList[index+1:]...)
		}
	}
}

func NewControllerManager(window *glfw.Window) *ControllerManager {
	var controllerList []Controller
	c := &ControllerManager{controllerList}
	window.SetKeyCallback(c.KeyCallback)
	window.SetMouseButtonCallback(c.MouseButtonCallback)
	window.SetCursorPosCallback(c.CursorPosCallback)
	window.SetScrollCallback(c.ScrollCallback)
	return c
}
