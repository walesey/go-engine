package ui

import (
	"strings"

	"github.com/walesey/go-engine/controller"
	vmath "github.com/walesey/go-engine/vectormath"
)

func controlKey(key controller.Key) bool {
	return key == controller.KeyLeftAlt ||
		key == controller.KeyLeftControl ||
		key == controller.KeyLeftShift ||
		key == controller.KeyRightAlt ||
		key == controller.KeyRightControl ||
		key == controller.KeyRightShift
}

func tabCycle(window *Window) {
	for i, tab := range window.Tabs {
		if tab.Active() {
			tab.Deactivate()
			if (i + 1) >= len(window.Tabs) {
				window.Tabs[0].Activate()
			} else {
				window.Tabs[i+1].Activate()
			}
			return
		}
	}
	if len(window.Tabs) > 0 {
		window.Tabs[0].Activate()
	}
}

func NewUiController(window *Window) controller.Controller {
	c := controller.CreateController()
	doMouseMove := func(xpos, ypos float64) {
		window.mouseMove(vmath.Vector2{xpos, ypos})
	}
	c.BindAxisAction(doMouseMove)
	c.BindMouseAction(func() { window.mouseClick(1, false) }, controller.MouseButton1, controller.Press)
	c.BindMouseAction(func() { window.mouseClick(2, false) }, controller.MouseButton2, controller.Press)
	c.BindMouseAction(func() { window.mouseClick(3, false) }, controller.MouseButton3, controller.Press)
	c.BindMouseAction(func() { window.mouseClick(4, false) }, controller.MouseButton4, controller.Press)
	c.BindMouseAction(func() { window.mouseClick(5, false) }, controller.MouseButton5, controller.Press)
	c.BindMouseAction(func() { window.mouseClick(1, true) }, controller.MouseButton1, controller.Release)
	c.BindMouseAction(func() { window.mouseClick(2, true) }, controller.MouseButton2, controller.Release)
	c.BindMouseAction(func() { window.mouseClick(3, true) }, controller.MouseButton3, controller.Release)
	c.BindMouseAction(func() { window.mouseClick(4, true) }, controller.MouseButton4, controller.Release)
	c.BindMouseAction(func() { window.mouseClick(5, true) }, controller.MouseButton5, controller.Release)

	var shift bool
	c.BindAction(func() { shift = true }, controller.KeyLeftShift, controller.Press)
	c.BindAction(func() { shift = false }, controller.KeyLeftShift, controller.Release)
	c.BindKeyAction(func(key controller.Key, action controller.Action) {
		if key == controller.KeyTab {
			if action == controller.Press {
				tabCycle(window)
			}
		} else if !controlKey(key) {
			keyString := string(byte(key))
			if shift {
				keyString = strings.ToUpper(keyString)
			} else {
				keyString = strings.ToLower(keyString)
			}
			switch {
			case key == controller.KeyBackspace:
				keyString = "backspace"
			case key == controller.KeyLeft:
				keyString = "leftArrow"
			case key == controller.KeyRight:
				keyString = "rightArrow"
			case key == controller.KeyUp:
				keyString = "upArrow"
			case key == controller.KeyDown:
				keyString = "downArrow"
			}
			window.keyClick(keyString, action == controller.Release)
		}
	})
	return c
}
