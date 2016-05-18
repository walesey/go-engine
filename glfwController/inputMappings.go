package glfwController

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/controller"
)

func getJoystick(joystick glfw.Joystick) controller.Joystick {
	return controller.Joystick(joystick)
}

func getKey(key glfw.Key) controller.Key {
	return controller.Key(key)
}

func getModifierKey(key glfw.ModifierKey) controller.ModifierKey {
	return controller.ModifierKey(key)
}

func getMouseButton(mouseButton glfw.MouseButton) controller.MouseButton {
	return controller.MouseButton(mouseButton)
}

func getAction(mouseButton glfw.Action) controller.Action {
	return controller.Action(mouseButton)
}
