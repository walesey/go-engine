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

	//physics engine
	physicsWorld := physics.NewPhysicsSpace()
	actorStore := actor.NewActorStore()
	for i := 0; i < 10; i = i + 1 {
		//make obj geometry with node
		geomgun := assetLib.GetGeometry("monkey")
		gunMat := assetLib.GetMaterial("monkeyMat")
		geomgun.Material = &gunMat
		gunNode := renderer.CreateNode()
		gunNode.Add(&geomgun)

		//create object with autgenerated colliders
		phyObj := physics.NewPhysicsObject()
		phyObj.Mass = 100
		assets.BoundingBoxFromGeometry(geomgun).AttachTo(&phyObj)
		assets.ConvexHullFromGeometry(geomgun).AttachTo(&phyObj)

		//attach to all the things ()
		actorStore.Add(actor.NewPhysicsActor(gunNode, &phyObj))
		physicsWorld.Add(&phyObj)
		sceneGraph.Add(gunNode)

		//set initial position
		phyObj.Position = vmath.Vector3{0, 5 * float64(i), float64(i) * 0.1}
	}

	//gravity global force
	physicsWorld.GlobalForces.AddForce("gravity", physics.GravityForce{vmath.Vector3{0, -10, 0}})

	//camera
	camera := renderer.CreateCamera(glRenderer)
	freeMoveActor := actor.CreateFreeMoveActor(camera)
	actorStore.Add(freeMoveActor)
	freeMoveActor.MoveSpeed = 3.0
	freeMoveActor.Location = vmath.Vector3{-2, 0, 0}

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
