package examples

import (
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/physics/dynamics"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"

	"github.com/codegangsta/cli"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

// Collisions - Collisions
func Collisions(c *cli.Context) {
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
	physicsWorld := dynamics.NewPhysicsSpace()
	actorStore := actor.NewActorStore()
	transformNode := renderer.CreateNode()
	transformNode.Add(&geomMonkey)
	transformNode.SetOrientation(vmath.AngleAxis(0.07, vmath.Vector3{0, 1, 0}))
	monkeyNode := renderer.CreateNode()
	monkeyNode.Add(transformNode)

	//create object with autgenerated colliders
	phyObj1 := physicsWorld.CreateObject()
	phyObj1.SetMass(100)
	phyObj1.SetBroadPhase(assets.BoundingBoxFromGeometry(geomMonkey))
	phyObj1.SetNarrowPhase(assets.ConvexSetFromGeometry(geomMonkey, 0.1))

	//attach to all the things ()
	freeMoveActor := actor.NewFreeMoveActor(monkeyNode)
	actorStore.Add(freeMoveActor)
	sceneGraph.Add(monkeyNode)

	//set initial position
	phyObj1.SetPosition(vmath.Vector3{0, 5, 0})

	terrain := assetLib.GetGeometry("terrain")
	terrainMat := assetLib.GetMaterial("terrainMat")
	terrain.Material = &terrainMat

	terrainNode := renderer.CreateNode()
	terrainNode.Add(&terrain)

	phyObj2 := physicsWorld.CreateObject()
	phyObj2.SetBroadPhase(assets.BoundingBoxFromGeometry(terrain))
	phyObj2.SetNarrowPhase(assets.ConvexSetFromGeometry(terrain, 0.1))
	phyObj2.SetPosition(vmath.Vector3{-10, -10, -10})

	actorStore.Add(actor.NewPhysicsActor(terrainNode, phyObj2))
	sceneGraph.Add(terrainNode)

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

		customController.BindAction(func() {
			phyObj1.SetPosition(freeMoveActor.Location)
			phyObj1.SetOrientation(monkeyNode.Orientation)
			physicsWorld.SimulateStep(0.018, 1)
			freeMoveActor.Location = phyObj1.GetPosition()
		}, glfw.KeyR, glfw.Press)

		//close window and exit on escape
		customController.BindAction(func() {
			glRenderer.Window.SetShouldClose(true)
		}, glfw.KeyEscape, glfw.Press)
	}

	glRenderer.Update = func() {
		fps.UpdateFPSMeter()
		actorStore.UpdateAll(0.018)
	}

	glRenderer.Render = func() {
		sceneGraph.RenderScene(glRenderer)
	}

	glRenderer.Start()
}
