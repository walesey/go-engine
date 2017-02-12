package main

import (
	"fmt"
	"go/build"
	"image/color"
	"math"
	"os"
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
	"github.com/walesey/go-engine/ui"
	"github.com/walesey/go-fileserver/client"
)

func init() {
	// Use all cpu cores
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Set default glfw controller
	controller.SetDefaultConstructor(glfwController.NewActionMap)
	// set working dir to access assets
	p, _ := build.Import("github.com/walesey/go-engine", "", build.FindOnly)
	os.Chdir(p.Dir)
}

//
func main() {
	glRenderer := opengl.NewOpenglRenderer("Lighting", 1920, 1080, true)
	gameEngine := engine.NewEngine(glRenderer)
	gameEngine.InitFpsDial()

	gameEngine.Start(func() {
		// Import the main shader
		shader, err := assets.ImportShader("shaders/build/pbr.vert", "shaders/build/pbr.frag")
		if err != nil {
			panic("error importing shader")
		}
		shader.AddTexture("diffuseMap")
		shader.AddTexture("normalMap")
		shader.AddTexture("aoMap")
		shader.AddTexture("metalnessMap")
		shader.AddTexture("roughnessMap")
		shader.AddTexture("glowMap")
		shader.FragDataLocations = []string{"outputColor", "brightColor"}
		gameEngine.DefaultShader(shader)

		// Import a glow post effect
		glowShader, err := assets.ImportShader("shaders/build/postEffects/glow.vert", "shaders/build/postEffects/glow.frag")
		if err != nil {
			fmt.Println("error importing glow shader", err)
		} else {
			glowShader.InputBuffers = 2
			glRenderer.CreatePostEffect(glowShader)
		}

		// enable anti aliasing with 4 samples
		glRenderer.AntiAliasing(4)

		// fetch assets from the server and then load the scene
		fetchAssets(gameEngine, func() {
			setupScene(gameEngine, shader)
		})

		// input/controller manager
		controllerManager := glfwController.NewControllerManager(glRenderer.Window)

		// camera + wasd controls
		freeMoveActor := actor.NewFreeMoveActor(gameEngine.Camera())
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

func fetchAssets(gameEngine engine.Engine, complete func()) {
	os.Mkdir("TestAssets", 0777)
	patcher := client.NewClient("TestAssets", "http://walesey.net")
	done := make(chan bool)
	go func() {
		patcher.SyncFiles("TestFiles")
		done <- true
	}()

	progressBar := ui.NewProgressBar("Downloading assets...")
	ui.SetProgressBar(progressBar, 0)
	gameEngine.AddOrtho(progressBar)
	progress := 0
	var loader engine.Updatable
	loader = engine.UpdatableFunc(func(dt float64) {
		select {
		case <-patcher.Complete:
			progress++
			ui.SetProgressBar(progressBar, 1+int((progress*20)/patcher.TotalFiles))
		case <-done:
			gameEngine.RemoveSpatial(progressBar, true)
			gameEngine.RemoveUpdatable(loader)
			gameEngine.RequestAnimationFrame(complete)
		default:
		}
	})
	gameEngine.AddUpdatable(loader)
}

func setupScene(gameEngine engine.Engine, shader *renderer.Shader) {
	camera := gameEngine.Camera()

	transparentNode := renderer.NewNode()
	gameEngine.AddSpatialTransparent(transparentNode)
	transparentNode.RendererParams = &renderer.RendererParams{
		DepthTest:    true,
		Unlit:        true,
		Transparency: renderer.EMISSIVE,
	}

	// Sky cubemap
	skyImg, err := assets.ImportImage("TestAssets/cloudSky.jpg")
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

	for i := 0; i < 2; i++ {
		torchLocation := mgl32.Vec3{0.86, 1.76, 1.05}
		if i == 1 {
			torchLocation = mgl32.Vec3{0.86, 1.76, -1.05}
		}

		fire := fireParticles()
		spark := sparkParticles()
		torchParticles := effects.NewParticleGroup(camera, fire, spark)
		torchParticles.SetTranslation(torchLocation)
		transparentNode.Add(torchParticles)
		gameEngine.AddUpdatable(torchParticles)

		light := renderer.NewLight(renderer.POINT)
		light.SetTranslation(torchLocation.Add(mgl32.Vec3{0, 0.05, 0}))
		gameEngine.AddLight(light)

		var x float64
		gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			x += dt
			mag := float32(math.Abs(0.6*math.Sin(3*x)+0.3*math.Sin(4*x)+0.15*math.Sin(7*x)+0.1*math.Sin(15*x))) + 0.5
			mag *= 0.05
			light.Color = [3]float32{1 * mag, 0.6 * mag, 0.4 * mag}
		}))
	}
}

func fireParticles() *effects.ParticleSystem {
	img, _ := assets.ImportImageCached("TestAssets/Fire.png")
	fire := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     12,
		ParticleEmitRate: 3,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		TotalFrames:      36,
		FramesX:          6,
		FramesY:          6,
		MaxLife:          0.8,
		MinLife:          1.5,
		StartColor:       color.NRGBA{255, 180, 80, 255},
		EndColor:         color.NRGBA{255, 60, 20, 255},
		StartSize:        mgl32.Vec3{1, 1, 1}.Mul(0.16),
		EndSize:          mgl32.Vec3{1, 1, 1}.Mul(0.23),
		MinTranslation:   mgl32.Vec3{1, 1, 1}.Mul(-0.01),
		MaxTranslation:   mgl32.Vec3{1, 1, 1}.Mul(0.01),
		MinStartVelocity: mgl32.Vec3{-0.02, 0, -0.02},
		MaxStartVelocity: mgl32.Vec3{0.02, 0.08, 0.02},
		MinRotation:      3.0,
		MaxRotation:      3.6,
	})
	fire.Node.Material = renderer.NewMaterial(renderer.NewTexture("glowMap", img, false))
	return fire
}

func sparkParticles() *effects.ParticleSystem {
	img, _ := assets.ImportImageCached("TestAssets/Spark.png")
	sparks := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:     7,
		ParticleEmitRate: 7,
		BaseGeometry:     renderer.CreateBox(float32(1), float32(1)),
		TotalFrames:      1,
		FramesX:          1,
		FramesY:          1,
		MaxLife:          1.0,
		MinLife:          0.7,
		StartColor:       color.NRGBA{255, 220, 180, 255},
		EndColor:         color.NRGBA{255, 155, 55, 255},
		StartSize:        mgl32.Vec3{1, 1, 1}.Mul(0.008),
		EndSize:          mgl32.Vec3{1, 1, 1}.Mul(0.005),
		MinTranslation:   mgl32.Vec3{1, 1, 1}.Mul(-0.04),
		MaxTranslation:   mgl32.Vec3{1, 1, 1}.Mul(0.04),
		MinStartVelocity: mgl32.Vec3{-0.5, 0.0, -0.5},
		MaxStartVelocity: mgl32.Vec3{0.5, 1.0, 0.5},
		Acceleration:     mgl32.Vec3{0.0, -1.0, 0.0},
		OnParticleUpdate: func(p *effects.Particle) {
			p.Scale[0] = p.Scale[0] * (1 + p.Velocity.Len()*3.5)
			p.Orientation = mgl32.QuatBetweenVectors(mgl32.Vec3{1, 0, 0}, p.Velocity)
		},
	})
	sparks.Node.Material = renderer.NewMaterial(renderer.NewTexture("glowMap", img, false))
	return sparks
}
