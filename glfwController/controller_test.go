package glfwController

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/stretchr/testify/assert"
	"github.com/walesey/go-engine/controller"
)

type TestObject struct {
	testState bool
}

func (to *TestObject) testAction() {
	to.testState = true
}
func (to *TestObject) otherTestAction() {
	to.testState = false
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestActionMap(t *testing.T) {
	var controllerList []Controller
	manager := &ControllerManager{controllerList}

	to := TestObject{false}
	c := NewActionMap()
	c.BindKeyAction(to.testAction, controller.KeyW, controller.Press)
	c.BindKeyAction(to.otherTestAction, controller.KeyE, controller.Release)
	manager.AddController(c.(Controller))

	fmt.Println("About to trigger custom actions")
	manager.KeyCallback(nil, glfw.KeyW, 0, glfw.Press, 0)
	assert.True(t, to.testState, "test state not triggered by key binding")
	manager.KeyCallback(nil, glfw.KeyE, 0, glfw.Release, 0)
	assert.False(t, to.testState, "test state not triggered by key binding")
}
