package examples

import (
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/physics"
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
		WindowWidth:  1900,
		WindowHeight: 1000,
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

	//physics engine
	physicsWorld := physics.NewPhysicsSpace()
	actorStore := actor.NewActorStore()
	for i := 0; i < 10; i = i + 1 {
		monkeyNode := renderer.CreateNode()
		monkeyNode.Add(&geomMonkey)

		//create object with autgenerated colliders
		phyObj := physics.NewPhysicsObject()
		phyObj.Mass = 100
		phyObj.Friction = 0.5
		phyObj.BroadPhase = assets.BoundingBoxFromGeometry(geomMonkey)
		phyObj.NarrowPhase = assets.ConvexSetFromGeometry(geomMonkey, 0.3)

		//attach to all the things ()
		actorStore.Add(actor.NewPhysicsActor(monkeyNode, phyObj))
		physicsWorld.Add(phyObj)
		sceneGraph.Add(monkeyNode)

		//set initial position
		phyObj.Position = vmath.Vector3{0, 5 * float64(i), 0}

		if i == 0 {
			phyObj.Static = true
		}
	}

	terrain := assetLib.GetGeometry("terrain")
	terrainMat := assetLib.GetMaterial("terrainMat")
	terrain.Material = &terrainMat

	for i := 0; i < 5; i = i + 1 {
		terrainNode := renderer.CreateNode()
		terrainNode.Add(&terrain)

		phyObj := physics.NewPhysicsObject()
		phyObj.Friction = 0.5
		phyObj.Static = true
		phyObj.BroadPhase = assets.BoundingBoxFromGeometry(terrain)
		phyObj.NarrowPhase = assets.ConvexSetFromGeometry(terrain, 2.0)

		actorStore.Add(actor.NewPhysicsActor(terrainNode, phyObj))
		physicsWorld.Add(phyObj)
		sceneGraph.Add(terrainNode)
		if i == 0 {
			phyObj.Position = vmath.Vector3{0, -15, 0}
			phyObj.Orientation = vmath.AngleAxis(0.57, vmath.Vector3{0, 0, 1})
		} else if i == 1 {
			phyObj.Position = vmath.Vector3{15, 0, 0}
			phyObj.Orientation = vmath.AngleAxis(0.7, vmath.Vector3{0, 0, 1})
		} else if i == 2 {
			phyObj.Position = vmath.Vector3{0, 0, 15}
			phyObj.Orientation = vmath.AngleAxis(0.7, vmath.Vector3{-1, 0, 0})
		} else if i == 3 {
			phyObj.Position = vmath.Vector3{-15, 0, 0}
			phyObj.Orientation = vmath.AngleAxis(0.7, vmath.Vector3{0, 0, -1})
		} else {
			phyObj.Position = vmath.Vector3{0, 0, -15}
			phyObj.Orientation = vmath.AngleAxis(0.7, vmath.Vector3{1, 0, 0})
		}
	}

	//gravity global force
	physicsWorld.GlobalForces.AddForce("gravity", physics.GravityForce{vmath.Vector3{0, -10, 0}})

	//camera
	camera := renderer.CreateCamera(glRenderer)
	freeMoveActor := actor.CreateFreeMoveActor(camera)
	actorStore.Add(freeMoveActor)
	freeMoveActor.MoveSpeed = 3.0
	freeMoveActor.Location = vmath.Vector3{-10, 0, -10}

	glRenderer.Init = func() {
		//lighting
		glRenderer.CreateLight(0.1, 0.1, 0.1, 3, 3, 3, 2, 2, 2, true, vmath.Vector3{0.3, -1, 0.2}, 0)

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
