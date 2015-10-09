package controller

import (
	"fmt"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestFPSActor struct {
	moveForward, moveBackward, strafLeft, strafRight bool
	jump, crouch, prone, sprint                      bool
}

func (tbma *TestFPSActor) StartMovingForward() {
	tbma.moveForward = true
}
func (tbma *TestFPSActor) StartMovingBackward() {
	tbma.moveBackward = true
}
func (tbma *TestFPSActor) StartStrafingLeft() {
	tbma.strafLeft = true
}
func (tbma *TestFPSActor) StartStrafingRight() {
	tbma.strafRight = true
}
func (tbma *TestFPSActor) StopMovingForward() {
	tbma.moveForward = false
}
func (tbma *TestFPSActor) StopMovingBackward() {
	tbma.moveBackward = false
}
func (tbma *TestFPSActor) StopStrafingLeft() {
	tbma.strafLeft = false
}
func (tbma *TestFPSActor) StopStrafingRight() {
	tbma.strafRight = false
}
func (tbma *TestFPSActor) Jump() {
	tbma.jump = true
}
func (tbma *TestFPSActor) Crouch() {
	tbma.crouch = true
}
func (tbma *TestFPSActor) StandUp() {
	tbma.crouch = false
}
func (tbma *TestFPSActor) Prone() {
	tbma.prone = true
}
func (tbma *TestFPSActor) StartSprinting() {
	tbma.sprint = true
}
func (tbma *TestFPSActor) StopSprinting() {
	tbma.sprint = false
}

func TestFPSController(t *testing.T) {
	var controllerList []Controller
	manager := &ControllerManager{controllerList}

	actor := &TestFPSActor{false, false, false, false, false, false, false, false}
	var c = NewFPSController(actor)
	manager.AddController(c)

	fmt.Println("About to test basic movement")
	manager.KeyCallback(nil, glfw.KeyW, 0, glfw.Press, 0)
	assert.True(t, actor.moveForward, "start moving forward")
	manager.KeyCallback(nil, glfw.KeyW, 0, glfw.Release, 0)
	assert.False(t, actor.moveForward, "stop moving forward")
	manager.KeyCallback(nil, glfw.KeyS, 0, glfw.Press, 0)
	assert.True(t, actor.moveBackward, "start moving backward")
	manager.KeyCallback(nil, glfw.KeyS, 0, glfw.Release, 0)
	assert.False(t, actor.moveBackward, "stop moving backward")
	manager.KeyCallback(nil, glfw.KeyA, 0, glfw.Press, 0)
	assert.True(t, actor.strafLeft, "start strafing left")
	manager.KeyCallback(nil, glfw.KeyA, 0, glfw.Release, 0)
	assert.False(t, actor.strafLeft, "stop strafing left")
	manager.KeyCallback(nil, glfw.KeyD, 0, glfw.Press, 0)
	assert.True(t, actor.strafRight, "start strafing right")
	manager.KeyCallback(nil, glfw.KeyD, 0, glfw.Release, 0)
	assert.False(t, actor.strafRight, "stop strafing right")
	manager.KeyCallback(nil, glfw.KeyLeftControl, 0, glfw.Press, 0)
	assert.True(t, actor.crouch, "crouch")
	manager.KeyCallback(nil, glfw.KeyLeftControl, 0, glfw.Release, 0)
	assert.False(t, actor.crouch, "standup")
	manager.KeyCallback(nil, glfw.KeySpace, 0, glfw.Press, 0)
	assert.True(t, actor.jump, "jump")
	manager.KeyCallback(nil, glfw.KeyZ, 0, glfw.Press, 0)
	assert.True(t, actor.prone, "Prone")
	manager.KeyCallback(nil, glfw.KeyLeftShift, 0, glfw.Press, 0)
	assert.True(t, actor.sprint, "start sprinting")
	manager.KeyCallback(nil, glfw.KeyLeftShift, 0, glfw.Release, 0)
	assert.False(t, actor.sprint, "stop sprinting")

	fmt.Println("Test unbound key, this should do nothing")
	manager.KeyCallback(nil, glfw.KeyX, 0, glfw.Press, 0)

}
