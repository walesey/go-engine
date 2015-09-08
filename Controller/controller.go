package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Controller interface {
	BindAction()
	TriggerAction()
	KeyCallback()
}

type KeyAction struct {
	key glfw.Key
	action glfw.Action
}

func Poll() {
	glfw.PollEvents()
}