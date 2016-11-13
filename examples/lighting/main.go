package main

import (
	"image/color"
	"runtime"

	"io/ioutil"

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

		//TODO: make a util function for importing shaders from file
		shader := renderer.NewShader()
		vertsrc, err := ioutil.ReadFile("build/shaders/basic.vert")
		if err != nil {
			panic(err)
		}
		shader.VertSrc = string(vertsrc)

		fragsrc, err := ioutil.ReadFile("build/shaders/basic.frag")
		if err != nil {
			panic(err)
		}
		shader.FragSrc = string(fragsrc)

		// Sky cubemap
		// skyImg, err := assets.ImportImage("resources/cubemapNightSky.jpg")
		// skyImg, err := assets.ImportImage("resources/space.jpg")
		skyImg, err := assets.ImportImage("resources/cubemap.png")
		if err == nil {
			geom := renderer.CreateSkyBox()
			geom.Shader = shader
			geom.Material = renderer.NewMaterial(renderer.NewTexture("diffuseMap", skyImg))
			geom.CullBackface = false
			geom.Transform(mgl32.Scale3D(10000, 10000, 10000))
			gameEngine.AddSpatial(geom)
		}

		// load scene objs
		objs := []string{
			"resources/wellScene/floor.obj",
			"resources/wellScene/frame1.obj",
			"resources/wellScene/frame2.obj",
			"resources/wellScene/well.obj",
			"resources/wellScene/torches.obj",
		}
		sceneNode := renderer.CreateNode()
		for _, objFile := range objs {
			if geom, err := assets.ImportObjCached(objFile); err == nil {
				geom.Shader = shader
				sceneNode.Add(geom)
			}
		}
		gameEngine.AddSpatial(sceneNode)

		torchLocation := mgl32.Vec3{0.86, 1.75, 1.05}
		fire := fireParticles(shader)
		torchParticles := effects.NewParticleGroup(camera, fire)
		torchParticles.SetTranslation(torchLocation)
		gameEngine.AddSpatialTransparent(torchParticles)
		gameEngine.AddUpdatable(torchParticles)

		// light := renderer.CreateLight()
		// light.Ambient = [3]float32{0.0, 0.0, 0.0}
		// light.Diffuse = [3]float32{0.03, 0.02, 0.003}
		// light.Specular = [3]float32{0.0, 0.0, 0.0}
		// light.SetTranslation(torchLocation)
		// gameEngine.AddLight(light)

		// var x float64
		// gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
		// 	x += dt
		// 	mag := float32(math.Abs(0.6*math.Sin(3*x)+0.3*math.Sin(4*x)+0.15*math.Sin(7*x)+0.1*math.Sin(15*x))) + 0.5
		// 	light.Diffuse = [3]float32{0.03 * mag, 0.02 * mag, 0.003 * mag}
		// 	light.SetTranslation(torchLocation)
		// }))

		// input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		// camera + wasd controls
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.Location = mgl32.Vec3{-6, 2, -6}
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

func fireParticles(shader *renderer.Shader) *effects.ParticleSystem {
	img, _ := assets.ImportImageCached("resources/fire.png")
	material := renderer.NewMaterial(renderer.NewTexture("diffuseMap", img))
	material.Transparency = renderer.EMISSIVE
	material.DepthMask = false
	return effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     12,
		ParticleEmitRate: 2,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		Material:         material,
		Shader:           shader,
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
