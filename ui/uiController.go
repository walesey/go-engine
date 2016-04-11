package ui

import (
	"strings"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/controller"
	vmath "github.com/walesey/go-engine/vectormath"
)

func controlKey(key glfw.Key) bool {
	return key == glfw.KeyLeftAlt ||
		key == glfw.KeyLeftControl ||
		key == glfw.KeyLeftShift ||
		key == glfw.KeyRightAlt ||
		key == glfw.KeyRightControl ||
		key == glfw.KeyRightShift
}

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

	var shift bool
	c.BindAction(func() { shift = true }, glfw.KeyLeftShift, glfw.Press)
	c.BindAction(func() { shift = false }, glfw.KeyLeftShift, glfw.Release)
	c.BindKeyAction(func(key glfw.Key, action glfw.Action) {
		if !controlKey(key) {
			keyString := string(byte(key))
			if shift {
				keyString = strings.ToUpper(keyString)
			} else {
				keyString = strings.ToLower(keyString)
			}
			if key == glfw.KeyBackspace {
				keyString = "backspace"
			}
			window.keyClick(keyString, action == glfw.Release)
		}
	})
	return c
}
