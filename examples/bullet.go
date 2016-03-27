package examples

import (
	"log"

	"github.com/luxengine/gobullet"
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/physics/bullet"
	"github.com/walesey/go-engine/physics/physicsAPI"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"

	"github.com/codegangsta/cli"
	"github.com/go-gl/glfw/v3.1/glfw"
)

//
func BulletDemo(c *cli.Context) {
	//Setup renderer and game Engine
	glRenderer := &renderer.OpenglRenderer{
		WindowTitle:  "GoEngine",
		WindowWidth:  1900,
		WindowHeight: 1000,
	}
	gameEngine := engine.NewEngine(glRenderer)

	//physics engine
	sdk := gobullet.NewBulletSDK()
	defer sdk.Delete()
	physicsWorld := bullet.NewBtDynamicsWorld(sdk)
	defer physicsWorld.Delete()
	physicsWorld.SetGravity(vmath.Vector3{0, -10, 0})
	gameEngine.AddUpdatable(physicsWorld)

	// assets
	assetLib, err := assets.LoadAssetLibrary("TestAssets/physics.asset")
	if err != nil {
		panic(err)
	}

	//geometry for physics objects
	geomMonkey := assetLib.GetGeometry("monkey")
	monkeyMat := assetLib.GetMaterial("monkeyMat")
	geomMonkey.Material = monkeyMat

	monkeyCollision := assets.CollisionShapeFromGeometry(geomMonkey, 0.3)

	spawn := func() physicsAPI.PhysicsObject {
		monkeyNode := renderer.CreateNode()
		monkeyNode.Add(geomMonkey)

		//create object with autgenerated colliders
		phyObj := bullet.NewBtRigidBody(100, monkeyCollision)
		physicsWorld.AddObject(phyObj)
		gameEngine.AddUpdatable(actor.NewPhysicsActor(monkeyNode, phyObj))
		gameEngine.AddSpatial(monkeyNode)

		return phyObj
	}

	for i := 0; i < 4; i = i + 1 {
		for j := 0; j < 4; j = j + 1 {
			for k := 0; k < 3; k = k + 1 {
				phyObj := spawn()

				//set initial position
				phyObj.SetPosition(vmath.Vector3{3*float64(i) - 5, 3*float64(k) + 15, 3*float64(j) - 5})
			}
		}
	}

	terrain := assetLib.GetGeometry("terrain")
	terrainMat := assetLib.GetMaterial("terrainMat")
	terrain.Material = terrainMat
	terrainCollision := assets.TriangleMeshShapeFromGeometry(assetLib.GetGeometry("terrain_lowpoli"))
	if err != nil {
		log.Printf("Error loading collision shape: %v\n", err)
	}

	terrainNode := renderer.CreateNode()
	terrainNode.Add(terrain)

	phyObj := bullet.NewBtRigidBodyConcave(0, terrainCollision)
	physicsWorld.AddObject(phyObj)

	gameEngine.
		AddUpdatable(actor.NewPhysicsActor(terrainNode, phyObj))
	gameEngine.AddSpatial(terrainNode)

	gameEngine.Start(func() {
		//lighting
		glRenderer.CreateLight(0.0, 0.0, 0.0, 0.5, 0.5, 0.5, 0.7, 0.7, 0.7, true, vmath.Vector3{0.3, -1, 0.2}, 0)

		//Sky
		gameEngine.Sky(assetLib.GetMaterial("skyboxMat"), 999999)

		//input/controller manager
		controllerManager := controller.NewControllerManager(glRenderer.Window)

		//lock the cursor
		glRenderer.LockCursor(true)

		//camera + player
		camera := renderer.CreateCamera(glRenderer)
		playerCollision := gobullet.NewBoxShape(1, 3, 1)
		playerController := bullet.NewBtCharacterController(playerCollision, 5)
		playerController.SetJumpSpeed(50)
		playerController.Warp(vmath.Vector3{0, 10, 0})
		physicsWorld.AddCharacterController(playerController)
		fpsActor := actor.NewFPSActor(camera, playerController)
		gameEngine.
			AddUpdatable(fpsActor)

		//fps controller
		mainController := controller.NewFPSController(fpsActor)
		controllerManager.AddController(mainController)

		//custom controller
		customController := controller.NewActionMap()
		controllerManager.AddController(customController)

		//spawn objects
		customController.BindAction(func() {
			phyObj := spawn()
			phyObj.SetPosition(camera.GetTranslation().Add(camera.GetDirection().MultiplyScalar(4)))
			phyObj.SetVelocity(camera.GetDirection().MultiplyScalar(30))
		}, glfw.KeyR, glfw.Press)

		//close window and exit on escape
		customController.BindAction(func() {
			glRenderer.Window.SetShouldClose(true)
		}, glfw.KeyEscape, glfw.Press)
	})
}
