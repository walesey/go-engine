package controller

import (
	"fmt"
	"testing"
	"github.com/go-gl/glfw/v3.1/glfw"
)

type TestFPSActor struct {

}
func (tbma *TestFPSActor) StartMovingForward() {
	fmt.Println("Start moving up")
}
func (tbma *TestFPSActor) StartMovingBackward() {
	fmt.Println("Start moving down")
}
func (tbma *TestFPSActor) StartStrafingLeft() {
	fmt.Println("Start moving left")
}
func (tbma *TestFPSActor) StartStrafingRight() {
	fmt.Println("Start moving right")
}
func (tbma *TestFPSActor) StopMovingForward() {
	fmt.Println("Stop moving up")
}
func (tbma *TestFPSActor) StopMovingBackward() {
	fmt.Println("Stop moving down")
}
func (tbma *TestFPSActor) StopStrafingLeft() {
	fmt.Println("Stop moving left")
}
func (tbma *TestFPSActor) StopStrafingRight() {
	fmt.Println("Stop moving right")
}
func (tbma *TestFPSActor) Jump() {
	fmt.Println("Jump")
}
func (tbma *TestFPSActor) StandUp() {
	fmt.Println("Stand up")
}
func (tbma *TestFPSActor) Crouch() {
	fmt.Println("Crouch")
}
func (tbma *TestFPSActor) Prone() {
	fmt.Println("Prone")
}
func (tbma *TestFPSActor) StartSprinting() {
	fmt.Println("Start Sprinting")
}
func (tbma *TestFPSActor) StopSprinting() {
	fmt.Println("Stop Sprinting")
}

func TestFPSController(t *testing.T) {
	var controllerList []Controller
	var manager = ControllerManager{controllerList}

	actor := &TestFPSActor{}
	var c = NewFPSController(actor)
	c.BindAction(testAction, glfw.KeyC, glfw.Press)
	c.BindAction(otherTestAction, glfw.KeyV, glfw.Release)
	manager.AddController(c)

	fmt.Println("About to trigger custom actions")
	manager.KeyCallback(nil, glfw.KeyC, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyV, 0, glfw.Release, 0)

	fmt.Println("About to test basic movement")
	manager.KeyCallback(nil, glfw.KeyW, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyW, 0, glfw.Release, 0)
	manager.KeyCallback(nil, glfw.KeyS, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyS, 0, glfw.Release, 0)
	manager.KeyCallback(nil, glfw.KeyA, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyA, 0, glfw.Release, 0)
	manager.KeyCallback(nil, glfw.KeyD, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyD, 0, glfw.Release, 0)
	manager.KeyCallback(nil, glfw.KeyLeftControl, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyLeftControl, 0, glfw.Release, 0)
	manager.KeyCallback(nil, glfw.KeySpace, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyZ, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyLeftShift, 0, glfw.Press, 0)
	manager.KeyCallback(nil, glfw.KeyLeftShift, 0, glfw.Release, 0)

	fmt.Println("Test unbound key, this should do nothing")
	manager.KeyCallback(nil, glfw.KeyX, 0, glfw.Press, 0)

}