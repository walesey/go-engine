package main

import (
	"image/color"
	"runtime"

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
		WindowTitle: "Particles",
		FullScreen:  true,
	}
	gameEngine := engine.NewEngine(glRenderer)
	gameEngine.InitFpsDial()

	gameEngine.Start(func() {

		// add particle effects
		fire := fireParticles()
		fire.Location = mgl32.Vec3{3, 0, 0}
		gameEngine.AddSpatialTransparent(fire)
		gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			fire.SetCameraLocation(gameEngine.Camera().GetTranslation())
			fire.Update(dt)
		}))

		// input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		// camera + wasd controls
		camera := gameEngine.Camera()
		freeMoveActor := actor.NewFreeMoveActor(camera)
		freeMoveActor.Location = mgl32.Vec3{}
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
	img, _ := assets.ImportImageCached("resources/fire.png")
	material := assets.CreateMaterial(img, nil, nil, nil)
	material.LightingMode = renderer.EMIT
	material.Transparency = renderer.EMISSIVE
	material.DepthMask = false
	return effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     12,
		ParticleEmitRate: 2,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		Material:         material,
		TotalFrames:      36,
		FramesX:          6,
		FramesY:          6,
		MaxLife:          1.5,
		MinLife:          2.5,
		StartColor:       color.NRGBA{254, 54, 0, 255},
		EndColor:         color.NRGBA{254, 100, 20, 200},
		StartSize:        mgl32.Vec3{0.3, 0.3, 0.3},
		EndSize:          mgl32.Vec3{0.9, 0.9, 0.9},
		MinTranslation:   mgl32.Vec3{-0.1, 0.1, -0.1},
		MaxTranslation:   mgl32.Vec3{0.1, 0.1, 0.1},
		MinStartVelocity: mgl32.Vec3{-0.02, -0.02, -0.02},
		MaxStartVelocity: mgl32.Vec3{0.02, 0.02, 0.02},
		Acceleration:     mgl32.Vec3{0.0, 0.0, 0.0},
		MinRotation:      3.0,
		MaxRotation:      3.6,
	})
}
