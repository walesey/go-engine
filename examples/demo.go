package examples

import (
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/renderer"

	"github.com/codegangsta/cli"
)

//
func Demo(c *cli.Context) {
	fps := renderer.CreateFPSMeter(1.0)
	fps.FpsCap = 60

	glRenderer := &renderer.OpenglRenderer{
		WindowTitle:  "GoEngine",
		WindowWidth:  1900,
		WindowHeight: 1000,
	}

	//setup scenegraph
	skyBox := assets.ImportObj("TestAssets/Files/skybox/skybox.obj")
	skyBox.Material.LightingMode = renderer.MODE_UNLIT
	skyBox.CullBackface = false
	skyNode := renderer.CreateNode()
	skyNode.Add(&skyBox)
	sceneGraph := renderer.CreateSceneGraph()
	sceneGraph.AddBackGround(skyNode)

	//camera
	camera := renderer.CreateCamera(glRenderer)
	freeMoveActor := actor.NewFreeMoveActor(camera)
	freeMoveActor.MoveSpeed = 3.0

	glRenderer.Init = func() {
		//input/controller manager
		controllerManager := controller.NewControllerManager(glRenderer.Window)
		//camera free move actor
		mainController := controller.NewBasicMovementController(freeMoveActor)
		controllerManager.AddController(mainController)
	}

	glRenderer.Update = func() {
		//update things that need updating
		fps.UpdateFPSMeter()
		freeMoveActor.Update(0.018)
	}

	glRenderer.Render = func() {
		//render the whole scene
		sceneGraph.RenderScene(glRenderer)
	}

	//start the renderer and launch the window
	glRenderer.Start()
}
