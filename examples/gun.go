package examples

import (
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/vectormath"

	"github.com/codegangsta/cli"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	// "github.com/go-gl/mathgl/mgl32"
)

//
func GunDemo(c *cli.Context) {
	fps := renderer.CreateFPSMeter(1.0)
	fps.FpsCap = 1960

	glRenderer := &opengl.OpenglRenderer{
		WindowTitle:  "GoEngine",
		WindowWidth:  1900,
		WindowHeight: 1000,
	}

	//setup scenegraph
	geom, _ := assets.ImportObj("TestAssets/Files/skybox/skybox.obj")
	geom.Material.LightingMode = renderer.MODE_UNLIT
	geom.CullBackface = false
	skyNode := renderer.CreateNode()
	skyNode.Add(geom)
	skyNode.SetRotation(1.57, vectormath.Vector3{0, 1, 0})
	skyNode.SetScale(vectormath.Vector3{5000, 5000, 5000})

	geomgun, _ := assets.ImportObj("TestAssets/Files/gun/rifle.obj")
	gunNode := renderer.CreateNode()
	gunNode.Add(geomgun)
	gunNode.SetTranslation(vectormath.Vector3{0, 0, 0})

	sceneGraph := renderer.CreateSceneGraph()
	sceneGraph.AddBackGround(skyNode)
	sceneGraph.Add(gunNode)

	//camera
	camera := renderer.CreateCamera(glRenderer)
	freeMoveActor := actor.NewFreeMoveActor(camera)
	freeMoveActor.MoveSpeed = 3.0
	freeMoveActor.Location = vectormath.Vector3{-2, 0, 0}

	glRenderer.Init(func() {
		//lighting
		glRenderer.CreateLight(0.1, 0.1, 0.1, 3, 3, 3, 2, 2, 2, true, vectormath.Vector3{0.3, -1, 0.2}, 0)

		//setup reflection map
		cubeMap := renderer.CreateCubemap(geom.Material.Diffuse)
		glRenderer.ReflectionMap(cubeMap)

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
		mainController := controller.NewBasicMovementController(freeMoveActor, false)
		controllerManager.AddController(mainController)

		//custom controller
		customController := controller.NewActionMap()
		controllerManager.AddController(customController)

		//close window and exit on escape
		customController.BindAction(func() {
			glRenderer.Window.SetShouldClose(true)
		}, glfw.KeyEscape, glfw.Press)
	})

	glRenderer.Update(func() {
		fps.UpdateFPSMeter()
		freeMoveActor.Update(0.018)
	})

	glRenderer.Render(func() {
		sceneGraph.RenderScene(glRenderer)
	})

	glRenderer.Start()
}
