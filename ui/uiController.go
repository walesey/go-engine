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

func tabCycle(window *Window, increment int) {
	for i, tab := range window.Tabs {
		if tab.Active() {
			tab.Deactivate()
			next := i + increment
			if next >= len(window.Tabs) {
				window.Tabs[0].Activate()
			} else if next < 0 {
				window.Tabs[len(window.Tabs)-1].Activate()
			} else {
				window.Tabs[next].Activate()
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
		window.mouseMove(vmath.Vector2{X: xpos, Y: ypos})
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
	c.BindKeyAction(func() { shift = true }, controller.KeyLeftShift, controller.Press)
	c.BindKeyAction(func() { shift = false }, controller.KeyLeftShift, controller.Release)
	c.SetKeyAction(func(key controller.Key, action controller.Action) {
		if key == controller.KeyTab {
			if action == controller.Press {
				if shift {
					tabCycle(window, -1)
				} else {
					tabCycle(window, 1)
				}
			}
		} else if !controlKey(key) {
			keyString := convertKey(key, shift)
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

func convertKey(key controller.Key, shift bool) string {
	keyString := strings.ToLower(string(byte(key)))
	if shift {
		keyString = strings.ToUpper(keyString)
		switch {
		case key == controller.Key1:
			keyString = "!"
		case key == controller.Key2:
			keyString = "@"
		case key == controller.Key3:
			keyString = "#"
		case key == controller.Key4:
			keyString = "$"
		case key == controller.Key5:
			keyString = "%"
		case key == controller.Key6:
			keyString = "^"
		case key == controller.Key7:
			keyString = "&"
		case key == controller.Key8:
			keyString = "*"
		case key == controller.Key9:
			keyString = "("
		case key == controller.Key0:
			keyString = ")"
		case key == controller.KeyMinus:
			keyString = "_"
		case key == controller.KeyEqual:
			keyString = "+"
		case key == controller.KeyLeftBracket:
			keyString = "{"
		case key == controller.KeyRightBracket:
			keyString = "}"
		case key == controller.KeyBackslash:
			keyString = "|"
		case key == controller.KeyGraveAccent:
			keyString = "~"
		case key == controller.KeySemicolon:
			keyString = ":"
		case key == controller.KeyApostrophe:
			keyString = "\""
		case key == controller.KeyComma:
			keyString = "<"
		case key == controller.KeyPeriod:
			keyString = ">"
		case key == controller.KeySlash:
			keyString = "?"
		}
	}
	return keyString
}
