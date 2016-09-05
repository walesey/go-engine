package main

import (
	"image/color"
	"runtime"

	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
}

func main() {

	glRenderer := &opengl.OpenglRenderer{
		WindowTitle:  "Platformer",
		WindowWidth:  800,
		WindowHeight: 800,
	}
	gameEngine := engine.NewEngine(glRenderer)

	gameEngine.Start(func() {

		// Sky cubemap
		skyImg, err := assets.ImportImage("resources/cubemap.png")
		if err == nil {
			gameEngine.Sky(assets.CreateMaterial(skyImg, nil, nil, nil), 999999)
		}

		// input/controller manager
		// controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		// The player is a green box
		boxGeometry := renderer.CreateBox(20, 20)
		boxGeometry.Material = renderer.CreateMaterial()
		boxGeometry.SetColor(color.NRGBA{0, 254, 0, 254})
		boxGeometry.CullBackface = false
		boxNode := renderer.CreateNode()
		boxNode.SetTranslation(vmath.Vector2{X: 400, Y: 400}.ToVector3())
		boxNode.Add(boxGeometry)
		gameEngine.AddOrtho(boxNode)
	})
}
