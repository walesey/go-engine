package examples

import (
	"fmt"

	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/physics"
	"github.com/walesey/go-engine/physics/collision"
	"github.com/walesey/go-engine/physics/dynamics"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"

	"github.com/codegangsta/cli"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//
func PhysicsDemo(c *cli.Context) {
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

	//geom := assetLib.GetGeometry("nightskybox")
	//skyboxMat := assetLib.GetMaterial("nightskyboxMat")
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
	radiusMonkey := assets.RadiusFromGeometry(geomMonkey)
	pointsMonkey := assets.PointsFromGeometry(geomMonkey, 0.01)

	//physics engine
	physicsWorld := physics.NewPhysicsSpace()
	actorStore := actor.NewActorStore()

	spawn := func() dynamics.PhysicsObject {
		monkeyNode := renderer.CreateNode()
		monkeyNode.Add(&geomMonkey)

		//create object with autgenerated colliders
		phyObj := physicsWorld.CreateObject()
		phyObj.SetBroadPhase(assets.BoundingBoxFromRadius(radiusMonkey))
		phyObj.SetNarrowPhase(collision.NewConvexSet(pointsMonkey))
		phyObj.SetMass(100)
		phyObj.SetRadius(radiusMonkey)
		phyObj.SetFriction(100.0)

		//attach to all the things ()
		actorStore.Add(actor.NewPhysicsActor(monkeyNode, phyObj))
		sceneGraph.Add(monkeyNode)

		return phyObj
	}

	for i := 0; i < 1; i = i + 1 {
		phyObj := spawn()

		//set initial position
		phyObj.SetPosition(vmath.Vector3{0.2 * float64(i), 4 * float64(i), 0.2 * float64(i)})

		if i == 0 {
			// phyObj.Static = true
			// phyObj.Velocity = vmath.Vector3{0, 5.6, 0}
		}
	}

	terrain := assetLib.GetGeometry("terrain")
	terrainMat := assetLib.GetMaterial("terrainMat")
	terrain.Material = &terrainMat

	for i := 0; i < 5; i = i + 1 {
		terrainNode := renderer.CreateNode()
		terrainNode.Add(&terrain)

		phyObj := physicsWorld.CreateObject()
		phyObj.SetFriction(300.0)
		phyObj.SetStatic(true)
		phyObj.SetBroadPhase(assets.BoundingBoxFromGeometry(terrain))
		phyObj.SetNarrowPhase(assets.ConvexSetFromGeometry(terrain, 2.0))

		actorStore.Add(actor.NewPhysicsActor(terrainNode, phyObj))
		sceneGraph.Add(terrainNode)
		if i == 0 {
			phyObj.SetPosition(vmath.Vector3{0, -4, 0})
		} else if i == 1 {
			phyObj.SetPosition(vmath.Vector3{14, 0, 0})
			phyObj.SetOrientation(vmath.AngleAxis(0.7, vmath.Vector3{0, 0, 1}))
		} else if i == 2 {
			phyObj.SetPosition(vmath.Vector3{0, 0, 14})
			phyObj.SetOrientation(vmath.AngleAxis(0.7, vmath.Vector3{-1, 0, 0}))
		} else if i == 3 {
			phyObj.SetPosition(vmath.Vector3{-14, 0, 0})
			phyObj.SetOrientation(vmath.AngleAxis(0.7, vmath.Vector3{0, 0, -1}))
		} else {
			phyObj.SetPosition(vmath.Vector3{0, 0, -14})
			phyObj.SetOrientation(vmath.AngleAxis(0.7, vmath.Vector3{1, 0, 0}))
		}
	}

	//gravity global force
	physicsWorld.SetGravity(vmath.Vector3{0, -20, 0})

	//debug
	// physicsWorld.OnEvent = func(event physics.Event) {
	// 	testNode := renderer.CreateNode()
	// 	testNode.Add(&geomMonkey)
	// 	sceneGraph.Add(testNode)
	// 	testNode.SetTranslation(event.Data.(map[string]interface{})["globalContact"].(vmath.Vector3))
	// 	testNode.SetScale(vmath.Vector3{0.3, 0.3, 0.3})
	// 	go func() {
	// 		time.Sleep(3000 * time.Millisecond)
	// 		sceneGraph.Remove(testNode)
	// 	}()
	// }

	//camera
	camera := renderer.CreateCamera(glRenderer)
	freeMoveActor := actor.CreateFreeMoveActor(camera)
	actorStore.Add(freeMoveActor)
	freeMoveActor.MoveSpeed = 10.0
	freeMoveActor.Location = vmath.Vector3{10, 0, -10}

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

		//camera free move actor
		mainController := controller.NewBasicMovementController(freeMoveActor)
		controllerManager.AddController(mainController)

		//custom controller
		customController := controller.NewActionMap()
		controllerManager.AddController(customController)

		//spawn objects
		customController.BindAction(func() {
			phyObj := spawn()
			fmt.Println(phyObj)
			phyObj.SetPosition(camera.Translation)
			phyObj.SetVelocity(camera.GetDirection().MultiplyScalar(30))
		}, glfw.KeyR, glfw.Press)

		//close window and exit on escape
		customController.BindAction(func() {
			glRenderer.Window.SetShouldClose(true)
		}, glfw.KeyEscape, glfw.Press)
	}

	glRenderer.Update = func() {
		fps.UpdateFPSMeter()
		physicsWorld.DoStep()
		actorStore.UpdateAll(0.018)
	}

	glRenderer.Render = func() {
		sceneGraph.RenderScene(glRenderer)
	}

	glRenderer.Start()
}
