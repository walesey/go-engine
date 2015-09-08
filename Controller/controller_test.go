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

type TestBasicMovementActor struct {

}
func (tbma *TestBasicMovementActor) StartMovingUp() {
	fmt.Println("Start moving up")
}
func (tbma *TestBasicMovementActor) StartMovingDown() {
	fmt.Println("Start moving down")
}
func (tbma *TestBasicMovementActor) StartMovingLeft() {
	fmt.Println("Start moving left")
}
func (tbma *TestBasicMovementActor) StartMovingRight() {
	fmt.Println("Start moving right")
}
func (tbma *TestBasicMovementActor) StopMovingUp() {
	fmt.Println("Stop moving up")
}
func (tbma *TestBasicMovementActor) StopMovingDown() {
	fmt.Println("Stop moving down")
}
func (tbma *TestBasicMovementActor) StopMovingLeft() {
	fmt.Println("Stop moving left")
}
func (tbma *TestBasicMovementActor) StopMovingRight() {
	fmt.Println("Stop moving right")
}

func TestMain(m *testing.M) {
	var controllerList []Controller
	var manager = ControllerManager{controllerList}

	actor := &TestBasicMovementActor{}
	var c = NewBasicMovementController(actor)
	c.BindAction(testAction, glfw.KeyW, glfw.Press)
	c.BindAction(otherTestAction, glfw.KeyE, glfw.Release)
	manager.AddController(c)

	fmt.Println("About to trigger custom actions")
	manager.KeyCallback(nil, glfw.KeyW, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyE, 0, glfw.Release, 0)

	fmt.Println("About to test basic movement")
	manager.KeyCallback(nil, glfw.KeyUp, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyUp, 0, glfw.Release, 0)
	manager.KeyCallback(nil, glfw.KeyDown, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyDown, 0, glfw.Release, 0)
	manager.KeyCallback(nil, glfw.KeyLeft, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyLeft, 0, glfw.Release, 0)
	manager.KeyCallback(nil, glfw.KeyRight, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyRight, 0, glfw.Release, 0)

	fmt.Println("Test unbound key, this should do nothing")
	manager.KeyCallback(nil, glfw.KeyX, 0, glfw.Press, 0)

}