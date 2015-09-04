package controller

import (
	"fmt"
	"testing"
)

func testAction() {
	fmt.Println("Test action triggered successfully")
}
func otherTestAction() {
	fmt.Println("Other Test action triggered successfully")
}

func TestMain(m *testing.M) {
	var c = Controller{make(map[string]func())}
	c.BindAction("w", testAction)
	c.BindAction("e", otherTestAction)
	fmt.Println("About to trigger actions")
	c.TriggerAction("w")
	c.TriggerAction("e")
}