package controller

import (
	"fmt"
	"testing"
	"os"
)

func testAction() {
	fmt.Println("Test action triggered successfully")
}
func otherTestAction() {
	fmt.Println("Other Test action triggered successfully")
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}