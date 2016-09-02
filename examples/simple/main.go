package examples

import (
	"image/color"

	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"

	"github.com/codegangsta/cli"
	vmath "github.com/walesey/go-engine/vectormath"
)

//
func Demo(c *cli.Context) {

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

		//input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		//camera + player
		camera := gameEngine.Camera()
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.MoveSpeed = 20.0
		freeMoveActor.LookSpeed = 0.002
		mainController := controller.NewBasicMovementController(freeMoveActor, true)
		controllerManager.AddController(mainController.(glfwController.Controller))
		gameEngine.AddUpdatable(freeMoveActor)

		customController := controller.CreateController()
		controllerManager.AddController(customController.(glfwController.Controller))

		mouse := vmath.Vector2{}
		customController.BindAxisAction(func(xpos, ypos float64) {
			mouse = vmath.Vector2{X: xpos, Y: ypos}
		})

		customController.BindAction(func() {
			vector := camera.GetMouseVector(mouse)
			rayGeom := renderer.CreateBeam(0.1, vector.MultiplyScalar(1000))
			rayGeom.SetColor(color.RGBA{R: 255, G: 0, B: 0, A: 130})
			mat := renderer.CreateMaterial()
			mat.LightingMode = renderer.MODE_UNLIT
			rayGeom.Material = mat
			rayNode := renderer.CreateNode()
			rayNode.Add(rayGeom)
			rayNode.SetTranslation(camera.GetTranslation().Subtract(vmath.Vector3{Y: 0.2}))
			gameEngine.AddSpatialTransparent(rayNode)
		}, controller.KeyQ, controller.Press)

		//Map loader
		assetLoader := assets.NewLoader()
		gameEngine.AddUpdatable(assetLoader)
		assetLoader.LoadMap("TestAssets/map.json", func(mapNode *renderer.Node, mapModel *editorModels.NodeModel) {
			gameEngine.AddSpatial(mapNode)
			// Define map class behaviours
			spiners := assets.FindNodeByClass("spiner", mapModel)
			for _, spiner := range spiners {
				node := spiner.GetNode()
				gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
					worldPosition, _ := mapModel.GetNode().RelativePosition(node)
					if freeMoveActor.Location.Subtract(worldPosition).LengthSquared() < 30.0 {
						node.SetOrientation(vmath.AngleAxis(dt, vmath.Vector3{0, 1, 0}).Multiply(node.Orientation))
					}
				}))
			}
		})
	})
}
