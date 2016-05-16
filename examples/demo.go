package examples

import (
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"

	"github.com/codegangsta/cli"
	vmath "github.com/walesey/go-engine/vectormath"
)

//
func Demo(c *cli.Context) {
	fps := renderer.CreateFPSMeter(1.0)
	fps.FpsCap = 60

	glRenderer := &opengl.OpenglRenderer{
		WindowTitle:  "GoEngine",
		WindowWidth:  1900,
		WindowHeight: 1000,
	}
	gameEngine := engine.NewEngine(glRenderer)

	gameEngine.Start(func() {
		//lighting
		glRenderer.CreateLight(0.0, 0.0, 0.0, 0.5, 0.5, 0.5, 0.7, 0.7, 0.7, true, vmath.Vector3{0.3, -1, 0.2}, 0)

		//Sky
		skyImg, err := assets.ImportImage("TestAssets/Files/skybox/cubemap.png")
		if err == nil {
			gameEngine.Sky(assets.CreateMaterial(skyImg, nil, nil, nil), 999999)
		}

		mapNode := renderer.CreateNode()
		mapModel := assets.LoadMap("TestAssets/map.map")
		mapModel.Root = assets.LoadMapToNode(mapModel.Root, mapNode)
		gameEngine.AddSpatial(mapNode)

		//input/controller manager
		controllerManager := controller.NewControllerManager(glRenderer.Window)

		//camera + player
		camera := gameEngine.Camera()
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.MoveSpeed = 20.0
		freeMoveActor.LookSpeed = 0.002
		mainController := controller.NewBasicMovementController(freeMoveActor, true)
		controllerManager.AddController(mainController)
		gameEngine.AddUpdatable(freeMoveActor)
	})
}
