package main

import (
	"image/color"
	"runtime"

	"github.com/walesey/go-engine/actor"
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

//
func main() {

	glRenderer := &opengl.OpenglRenderer{
		WindowTitle:  "Simple",
		WindowWidth:  800,
		WindowHeight: 800,
	}
	gameEngine := engine.NewEngine(glRenderer)
	gameEngine.InitFpsDial()

	gameEngine.Start(func() {

		//lighting
		glRenderer.CreateLight(
			0.3, 0.3, 0.3, //ambient
			0.5, 0.5, 0.5, //diffuse
			0.7, 0.7, 0.7, //specular
			false, vmath.Vector3{0.7, 0.2, 0.7}, //position
			0, //index
		)

		// Sky cubemap
		skyImg, err := assets.ImportImage("resources/cubemap.png")
		if err == nil {
			gameEngine.Sky(assets.CreateMaterial(skyImg, nil, nil, nil), 999999)
		}

		// input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		// camera + wasd controls
		camera := gameEngine.Camera()
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.Location = vmath.Vector3{}
		mainController := controller.NewBasicMovementController(freeMoveActor, false)
		controllerManager.AddController(mainController.(glfwController.Controller))
		gameEngine.AddUpdatable(freeMoveActor)

		//lock the cursor
		glRenderer.LockCursor(true)

		// Create a red box geometry, attach to a node, add the node to the scenegraph
		boxGeometry := renderer.CreateBox(10, 10)
		boxGeometry.Material = renderer.CreateMaterial()
		boxGeometry.SetColor(color.NRGBA{254, 0, 0, 254})
		boxGeometry.CullBackface = false
		boxNode := renderer.CreateNode()
		boxNode.SetTranslation(vmath.Vector3{X: 30})
		boxNode.Add(boxGeometry)
		gameEngine.AddSpatial(boxNode)

		// make the box spin
		var angle float64
		gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			angle += dt
			boxNode.SetOrientation(vmath.AngleAxis(angle, vmath.Vector3{Y: 1}))
		}))

		// custom key bindings
		customController := controller.CreateController()
		controllerManager.AddController(customController.(glfwController.Controller))

		// close window and exit on escape
		customController.BindKeyAction(func() {
			glRenderer.Window.SetShouldClose(true)
		}, controller.KeyEscape, controller.Press)
	})
}
