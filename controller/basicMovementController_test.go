package controller

import (
	"fmt"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestBasicMovementActor struct {
	moveUp, moveDown, moveLeft, moveRight bool
}

func (tbma *TestBasicMovementActor) StartMovingUp() {
	tbma.moveUp = true
}
func (tbma *TestBasicMovementActor) StartMovingDown() {
	tbma.moveDown = true
}
func (tbma *TestBasicMovementActor) StartMovingLeft() {
	tbma.moveLeft = true
}
func (tbma *TestBasicMovementActor) StartMovingRight() {
	tbma.moveRight = true
}
func (tbma *TestBasicMovementActor) StopMovingUp() {
	tbma.moveUp = false
}
func (tbma *TestBasicMovementActor) StopMovingDown() {
	tbma.moveDown = false
}
func (tbma *TestBasicMovementActor) StopMovingLeft() {
	tbma.moveLeft = false
}
func (tbma *TestBasicMovementActor) StopMovingRight() {
	tbma.moveRight = false
}

func TestBasicMovementController(t *testing.T) {
	var controllerList []Controller
	manager := &ControllerManager{controllerList}

	actor := &TestBasicMovementActor{false, false, false, false}
	var c = NewBasicMovementController(actor)
	manager.AddController(c)

	fmt.Println("About to test basic movement")
	manager.KeyCallback(nil, glfw.KeyUp, 0, glfw.Press, 0)
	assert.True(t, actor.moveUp, "start moving up")
	manager.KeyCallback(nil, glfw.KeyUp, 0, glfw.Release, 0)
	assert.False(t, actor.moveUp, "stop moving up")
	manager.KeyCallback(nil, glfw.KeyDown, 0, glfw.Press, 0)
	assert.True(t, actor.moveDown, "start moving down")
	manager.KeyCallback(nil, glfw.KeyDown, 0, glfw.Release, 0)
	assert.False(t, actor.moveDown, "stop moving down")
	manager.KeyCallback(nil, glfw.KeyLeft, 0, glfw.Press, 0)
	assert.True(t, actor.moveLeft, "start moving left")
	manager.KeyCallback(nil, glfw.KeyLeft, 0, glfw.Release, 0)
	assert.False(t, actor.moveLeft, "stop moving left")
	manager.KeyCallback(nil, glfw.KeyRight, 0, glfw.Press, 0)
	assert.True(t, actor.moveRight, "start moving right")
	manager.KeyCallback(nil, glfw.KeyRight, 0, glfw.Release, 0)
	assert.False(t, actor.moveRight, "stop moving right")

	fmt.Println("Test unbound key, this should do nothing")
	manager.KeyCallback(nil, glfw.KeyX, 0, glfw.Press, 0)

}
