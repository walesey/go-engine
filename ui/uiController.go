package ui

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/controller"
	vmath "github.com/walesey/go-engine/vectormath"
)

func NewUiController(window *Window) controller.Controller {
	c := controller.NewActionMap()
	doMouseMove := func(xpos, ypos float64) {
		window.mouseMove(vmath.Vector2{xpos, ypos})
	}
	c.BindAxisAction(doMouseMove)
	c.BindMouseAction(func() { window.mouseClick(1, false) }, glfw.MouseButton1, glfw.Press)
	c.BindMouseAction(func() { window.mouseClick(2, false) }, glfw.MouseButton2, glfw.Press)
	c.BindMouseAction(func() { window.mouseClick(3, false) }, glfw.MouseButton3, glfw.Press)
	c.BindMouseAction(func() { window.mouseClick(4, false) }, glfw.MouseButton4, glfw.Press)
	c.BindMouseAction(func() { window.mouseClick(5, false) }, glfw.MouseButton5, glfw.Press)
	c.BindMouseAction(func() { window.mouseClick(1, true) }, glfw.MouseButton1, glfw.Release)
	c.BindMouseAction(func() { window.mouseClick(2, true) }, glfw.MouseButton2, glfw.Release)
	c.BindMouseAction(func() { window.mouseClick(3, true) }, glfw.MouseButton3, glfw.Release)
	c.BindMouseAction(func() { window.mouseClick(4, true) }, glfw.MouseButton4, glfw.Release)
	c.BindMouseAction(func() { window.mouseClick(5, true) }, glfw.MouseButton5, glfw.Release)
	return c
}
