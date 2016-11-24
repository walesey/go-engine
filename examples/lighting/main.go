package main

import (
	"image/color"
	"math"
	"runtime"

	"github.com/disintegration/imaging"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/effects"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/opengl"
	"github.com/walesey/go-engine/renderer"
)

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
}

//
func main() {

	glRenderer := &opengl.OpenglRenderer{
		WindowTitle: "Lighting",
		FullScreen:  true,
	}
	gameEngine := engine.NewEngine(glRenderer)
	camera := gameEngine.Camera()
	gameEngine.InitFpsDial()

	gameEngine.Start(func() {

		transparentNode := renderer.NewNode()
		gameEngine.AddSpatialTransparent(transparentNode)
		transparentNode.RendererParams = &renderer.RendererParams{
			DepthTest:    true,
			Unlit:        true,
			Transparency: renderer.EMISSIVE,
		}

		shader, err := assets.ImportShader("build/shaders/pbr.vert", "build/shaders/pbr.frag")
		if err != nil {
			panic("error importing shader")
		}

		gameEngine.DefaultShader(shader)

		// Sky cubemap
		skyImg, err := assets.ImportImage("TestAssets/Files/skybox/cloudSky.jpg")
		if err == nil {
			skyImg = imaging.AdjustBrightness(skyImg, -30)
			skyImg = imaging.AdjustContrast(skyImg, 30)
			geom := renderer.CreateSkyBox()
			geom.Transform(mgl32.Scale3D(10000, 10000, 10000))
			skyNode := renderer.NewNode()
			skyNode.SetOrientation(mgl32.QuatRotate(1.57, mgl32.Vec3{0, 1, 0}))
			skyNode.Material = renderer.NewMaterial(renderer.NewTexture("diffuseMap", skyImg, false))
			skyNode.RendererParams = renderer.NewRendererParams()
			skyNode.RendererParams.CullBackface = false
			skyNode.RendererParams.Unlit = true
			skyNode.Add(geom)
			gameEngine.AddSpatial(skyNode)
			// create an environmentMap using the skybox texture
			envCubeMap := renderer.NewCubemap("environmentMap", skyImg, true)
			gameEngine.DefaultCubeMap(envCubeMap)
		}

		// load scene objs
		objs := []string{
			"TestAssets/wellScene/floor.obj",
			"TestAssets/wellScene/frame1.obj",
			"TestAssets/wellScene/frame2.obj",
			"TestAssets/wellScene/well.obj",
			"TestAssets/wellScene/torches.obj",
		}
		for _, objFile := range objs {
			if geom, mat, err := assets.ImportObjCached(objFile); err == nil {
				sceneNode := renderer.NewNode()
				sceneNode.Add(geom)
				sceneNode.Material = mat
				sceneNode.RendererParams = renderer.NewRendererParams()
				sceneNode.RendererParams.CullBackface = false
				gameEngine.AddSpatial(sceneNode)
			}
		}

		fireImage, _ := assets.ImportImageCached("resources/fire.png")
		fireMat := renderer.NewMaterial(renderer.NewTexture("diffuseMap", fireImage, false))

		torchLocation := mgl32.Vec3{0.86, 1.75, 1.05}
		fire := fireParticles()
		torchParticles := effects.NewParticleGroup(camera, fire)
		torchParticles.Node.Material = fireMat
		torchParticles.SetTranslation(torchLocation)
		transparentNode.Add(torchParticles)
		gameEngine.AddUpdatable(torchParticles)

		var x float64
		gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			x += dt
			mag := float32(math.Abs(0.6*math.Sin(3*x)+0.3*math.Sin(4*x)+0.15*math.Sin(7*x)+0.1*math.Sin(15*x))) + 0.5
			mag *= 10
			lightPos := torchLocation.Add(mgl32.Vec3{0, 0.05, 0})

			// shader Lighting
			shader.Uniforms["nbPointLights"] = 2
			shader.Uniforms["pointLightPositions"] = []float32{
				lightPos.X(), lightPos.Y(), lightPos.Z(), 0,
				lightPos.X(), lightPos.Y(), -lightPos.Z(), 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			}
			shader.Uniforms["pointLightValues"] = []float32{
				0.03 * mag, 0.02 * mag, 0.003 * mag, 0,
				0.03 * mag, 0.02 * mag, 0.003 * mag, 0,
				0, 0, 0, 0,
				0, 0, 0, 0,
			}
		}))

		// input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		// camera + wasd controls
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.Location = mgl32.Vec3{-6, 2, -6}
		freeMoveActor.MoveSpeed = 1.5
		mainController := controller.NewBasicMovementController(freeMoveActor, false)
		controllerManager.AddController(mainController.(glfwController.Controller))
		gameEngine.AddUpdatable(freeMoveActor)

		//lock the cursor
		glRenderer.LockCursor(true)

		// custom key bindings
		customController := controller.CreateController()
		controllerManager.AddController(customController.(glfwController.Controller))

		// close window and exit on escape
		customController.BindKeyAction(func() {
			glRenderer.Window.SetShouldClose(true)
		}, controller.KeyEscape, controller.Press)
	})
}

func fireParticles() *effects.ParticleSystem {
	return effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     12,
		ParticleEmitRate: 2,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		TotalFrames:      36,
		FramesX:          6,
		FramesY:          6,
		MaxLife:          0.8,
		MinLife:          1.5,
		StartColor:       color.NRGBA{254, 54, 0, 255},
		EndColor:         color.NRGBA{254, 100, 20, 255},
		StartSize:        mgl32.Vec3{0.1, 0.1, 0.1},
		EndSize:          mgl32.Vec3{0.15, 0.15, 0.15},
		MinTranslation:   mgl32.Vec3{-0.01, 0.01, -0.01},
		MaxTranslation:   mgl32.Vec3{0.01, 0.01, 0.01},
		MinStartVelocity: mgl32.Vec3{-0.02, -0.02, -0.02},
		MaxStartVelocity: mgl32.Vec3{0.02, 0.02, 0.02},
		Acceleration:     mgl32.Vec3{0.0, 0.0, 0.0},
		MinRotation:      3.0,
		MaxRotation:      3.6,
	})
}
