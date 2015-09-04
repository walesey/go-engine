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
	var c = Controller{make(map[glfw.Key]func()), nil}
	c.BindAction(testAction, glfw.KeyW)
	c.BindAction(otherTestAction, glfw.KeyE)
	fmt.Println("About to trigger actions")
	c.TriggerAction(glfw.KeyW)
	c.TriggerAction(glfw.KeyE)
}