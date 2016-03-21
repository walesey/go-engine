package examples

import (
	"log"

	"github.com/luxengine/gobullet"
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/physics/bullet"
	"github.com/walesey/go-engine/physics/physicsAPI"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"

	"github.com/codegangsta/cli"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//
func BulletDemo(c *cli.Context) {
	fps := renderer.CreateFPSMeter(1.0)
	fps.FpsCap = 60

	glRenderer := &renderer.OpenglRenderer{
		WindowTitle:  "GoEngine",
		WindowWidth:  990,
		WindowHeight: 900,
	}

	assetLib, err := assets.LoadAssetLibrary("TestAssets/physics.asset")
	if err != nil {
		panic(err)
	}

	//setup scenegraph

	// geom := assetLib.GetGeometry("nightskybox")
	// skyboxMat := assetLib.GetMaterial("nightskyboxMat")
	geom := assetLib.GetGeometry("skybox")
	skyboxMat := assetLib.GetMaterial("skyboxMat")
	geom.Material = &skyboxMat
	geom.Material.LightingMode = renderer.MODE_UNLIT
	geom.CullBackface = false
	skyNode := renderer.CreateNode()
	skyNode.Add(&geom)
	skyNode.SetRotation(1.57, vmath.Vector3{0, 1, 0})
	skyNode.SetScale(vmath.Vector3{5000, 5000, 5000})

	sceneGraph := renderer.CreateSceneGraph()
	sceneGraph.AddBackGround(skyNode)

	//geometry for physics objects
	geomMonkey := assetLib.GetGeometry("monkey")
	monkeyMat := assetLib.GetMaterial("monkeyMat")
	geomMonkey.Material = &monkeyMat

	monkeyCollision := assets.CollisionShapeFromGeometry(geomMonkey, 0.3)

	//physics engine
	sdk := gobullet.NewBulletSDK()
	defer sdk.Delete()
	physicsWorld := bullet.NewBtDynamicsWorld(sdk)
	defer physicsWorld.Delete()
	actorStore := actor.NewActorStore()

	spawn := func() physicsAPI.PhysicsObject {
		monkeyNode := renderer.CreateNode()
		monkeyNode.Add(&geomMonkey)

		//create object with autgenerated colliders
		phyObj := bullet.NewBtRigidBody(100, monkeyCollision)
		physicsWorld.AddObject(phyObj)

		//attach to all the things
		actorStore.Add(actor.NewPhysicsActor(monkeyNode, phyObj))
		sceneGraph.Add(monkeyNode)

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
	terrain.Material = &terrainMat
	terrainCollision := assets.TriangleMeshShapeFromGeometry(assetLib.GetGeometry("terrain_lowpoli"))
	if err != nil {
		log.Printf("Error loading collision shape: %v\n", err)
	}

	terrainNode := renderer.CreateNode()
	terrainNode.Add(&terrain)

	phyObj := bullet.NewBtRigidBodyConcave(0, terrainCollision)
	physicsWorld.AddObject(phyObj)

	actorStore.Add(actor.NewPhysicsActor(terrainNode, phyObj))
	sceneGraph.Add(terrainNode)

	//gravity global force
	physicsWorld.SetGravity(vmath.Vector3{0, -10, 0})

	glRenderer.Init = func() {
		//lighting
		glRenderer.CreateLight(0.0, 0.0, 0.0, 0.5, 0.5, 0.5, 0.7, 0.7, 0.7, true, vmath.Vector3{0.3, -1, 0.2}, 0)

		//setup reflection map
		//cubeMap := renderer.CreateCubemap(assetLib.GetMaterial("nightskyboxMat").Diffuse)
		cubeMap := renderer.CreateCubemap(assetLib.GetMaterial("skyboxMat").Diffuse)
		glRenderer.ReflectionMap(*cubeMap)

		//post effects
		bloomHorizontal := renderer.Shader{
			Name: "shaders/bloom/bloomHorizontal",
			Uniforms: []renderer.Uniform{
				renderer.Uniform{"size", mgl32.Vec2{1900, 1000}},
				renderer.Uniform{"quality", 2.5},
				renderer.Uniform{"samples", 12},
				renderer.Uniform{"threshold", 0.995},
				renderer.Uniform{"intensity", 1.9},
			},
		}
		bloomVertical := renderer.Shader{
			Name: "shaders/bloom/bloomVertical",
			Uniforms: []renderer.Uniform{
				renderer.Uniform{"size", mgl32.Vec2{1900, 1000}},
				renderer.Uniform{"quality", 2.5},
				renderer.Uniform{"samples", 12},
				renderer.Uniform{"threshold", 0.995},
				renderer.Uniform{"intensity", 1.9},
			},
		}
		glRenderer.CreatePostEffect(bloomVertical)
		glRenderer.CreatePostEffect(bloomHorizontal)

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
		actorStore.Add(fpsActor)

		//fps controller
		mainController := controller.NewFPSController(fpsActor)
		controllerManager.AddController(mainController)

		//custom controller
		customController := controller.NewActionMap()
		controllerManager.AddController(customController)

		//spawn objects
		customController.BindAction(func() {
			phyObj := spawn()
			phyObj.SetPosition(camera.Translation.Add(camera.GetDirection().MultiplyScalar(4)))
			phyObj.SetVelocity(camera.GetDirection().MultiplyScalar(30))
		}, glfw.KeyR, glfw.Press)

		//close window and exit on escape
		customController.BindAction(func() {
			glRenderer.Window.SetShouldClose(true)
		}, glfw.KeyEscape, glfw.Press)
	}

	glRenderer.Update = func() {
		fps.UpdateFPSMeter()
		physicsWorld.SimulateStep(6, 1)
		actorStore.UpdateAll(0.018)

		// if playerController.GetPosition().Y < -5 {
		// 	playerController.Warp(vmath.Vector3{0, 10, 0})
		// }
	}

	glRenderer.Render = func() {
		sceneGraph.RenderScene(glRenderer)
	}

	glRenderer.Start()
}
