package controller

import (
	"fmt"
	"testing"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func testAction() {
	fmt.Println("Test action triggered successfully")
}
func otherTestAction() {
	fmt.Println("Other Test action triggered successfully")
}

func TestMain(m *testing.M) {
	var c = BasicMovementController{nil, make(map[KeyAction]func()), nil}
	c.BindAction(testAction, glfw.KeyW, glfw.Press)
	c.BindAction(otherTestAction, glfw.KeyE, glfw.Release)
	fmt.Println("About to trigger actions")
	c.TriggerAction(glfw.KeyW, glfw.Press)
	c.TriggerAction(glfw.KeyE, glfw.Release)
}