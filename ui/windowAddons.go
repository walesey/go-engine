package ui

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/controller"
)

func ClickAndDragWindow(window *Window, hitbox Hitbox, c controller.Controller) {
	grabbed := false
	grabOffset := mgl32.Vec2{}
	hitbox.AddOnClick(func(button int, release bool, position mgl32.Vec2) {
		grabOffset = position
		grabbed = !release
	})
	c.BindMouseAction(func() {
		grabbed = false
	}, controller.MouseButton1, controller.Release)
	c.BindAxisAction(func(xpos, ypos float32) {
		if grabbed {
			position := mgl32.Vec2{xpos, ypos}
			window.SetTranslation(position.Sub(grabOffset).Vec3(0))
		}
	})
}
