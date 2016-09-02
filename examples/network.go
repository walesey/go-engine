package examples

import (
	"fmt"
	"time"

	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/networking"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"

	"github.com/codegangsta/cli"
	vmath "github.com/walesey/go-engine/vectormath"
)

//
func Network(c *cli.Context) {
	server := len(c.Args()) > 0 && c.Args()[0] == "server"
	var gameEngine engine.Engine
	var glRenderer *opengl.OpenglRenderer
	network := networking.NewNetwork()

	if server {
		gameEngine = engine.NewHeadlessEngine()
		network.StartServer(9999)
	} else {
		glRenderer = &opengl.OpenglRenderer{
			WindowTitle:  "GoEngine",
			WindowWidth:  1900,
			WindowHeight: 1000,
		}
		gameEngine = engine.NewEngine(glRenderer)
		network.ConnectClient("127.0.0.1:9999")
	}

	//Networked Game events
	network.RegisterEvent("spawn", func(clientId string, data []byte) {
		if network.IsServer() {
			// gameEngine.AddUpdatableKey(id, nil)
			fmt.Println(data)
		}
		if network.IsClient() {
			geomMonkey, _ := assets.ImportObjCached("TestAssets/Files/physicsMonkey/phyMonkey.obj")
			monkeyNode := renderer.CreateNode()
			monkeyNode.Add(geomMonkey)
			monkeyNode.SetTranslation(vmath.Vector3{})
			gameEngine.AddSpatial(monkeyNode)
		}
	})

	network.RegisterEvent("clientEvent", func(clientId string, data []byte) {

	})

	gameEngine.AddUpdatable(network)
	gameEngine.Start(func() {
		if network.IsServer() {
			ids := 0
			position := vmath.Vector3{}
			gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
				fmt.Println("test3 ", ids)
				ids++
				position = position.Add(vmath.Vector3{5, 0, 0})
				network.TriggerOnServerAndClients("spawn", []byte("test data"))
				time.Sleep(1500 * time.Millisecond)
			}))
		}

		if network.IsClient() {
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

			//custom controller
			customController := controller.CreateController()
			controllerManager.AddController(customController.(glfwController.Controller))

			//spawn objects
			customController.BindAction(func() {
				network.TriggerOnServerAndClients("clientEvent", []byte{})
			}, controller.KeyR, controller.Press)
		}
	})
}
