package ui

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/controller"
	vmath "github.com/walesey/go-engine/vectormath"
)

func ClickAndDragWindow(window *Window, hitbox Hitbox, actionMap *controller.ActionMap) {
	grabbed := false
	grabOffset := vmath.Vector2{}
	hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
		grabOffset = position
		grabbed = !release
	})
	actionMap.BindMouseAction(func() {
		grabbed = false
	}, glfw.MouseButton1, glfw.Release)
	actionMap.BindAxisAction(func(xpos, ypos float64) {
		if grabbed {
			position := vmath.Vector2{xpos, ypos}
			window.SetTranslation(position.Subtract(grabOffset).ToVector3())
		}
	})
}
