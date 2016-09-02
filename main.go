package main

import (
	"runtime"

	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/editor"
	"github.com/walesey/go-engine/glfwController"
)

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
}

func main() {
	editor.New().Start()
}
