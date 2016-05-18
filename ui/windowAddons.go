package ui

import (
	"github.com/walesey/go-engine/controller"
	vmath "github.com/walesey/go-engine/vectormath"
)

func ClickAndDragWindow(window *Window, hitbox Hitbox, c controller.Controller) {
	grabbed := false
	grabOffset := vmath.Vector2{}
	hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
		grabOffset = position
		grabbed = !release
	})
	c.BindMouseAction(func() {
		grabbed = false
	}, controller.MouseButton1, controller.Release)
	c.BindAxisAction(func(xpos, ypos float64) {
		if grabbed {
			position := vmath.Vector2{xpos, ypos}
			window.SetTranslation(position.Subtract(grabOffset).ToVector3())
		}
	})
}

func DeactivateAllTextElements(container *Container) {
	for i := 0; i < container.GetNbChildren(); i++ {
		child := container.GetChild(i)
		childContainer, ok := child.(*Container)
		if ok {
			DeactivateAllTextElements(childContainer)
		}
		text, ok := child.(*TextElement)
		if ok {
			text.Deactivate()
		}
	}
}
