package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Controller interface {
	BindAction(function func(), key glfw.Key, action glfw.Action)
	TriggerAction(key glfw.Key, action glfw.Action)
	KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
}

type KeyAction struct {
	key glfw.Key
	action glfw.Action
}

func Poll() {
	glfw.PollEvents()
}