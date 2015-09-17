package examples

import (
	"image/color"

	"github.com/walesey/go-engine/actor"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/effects"
	"github.com/walesey/go-engine/renderer"
	"github.com/walesey/go-engine/vectormath"

	"github.com/codegangsta/cli"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

//
func Particles(c *cli.Context) {
	fps := renderer.CreateFPSMeter(1.0)
	fps.FpsCap = 6000

	glRenderer := &renderer.OpenglRenderer{
		WindowTitle:  "GoEngine",
		WindowWidth:  1900,
		WindowHeight: 1000,
	}

	assetLib, err := assets.LoadAssetLibrary("TestAssets/particles.asset")
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
	skyNode.SetRotation(1.57, vectormath.Vector3{0, 1, 0})
	skyNode.SetScale(vectormath.Vector3{5000, 5000, 5000})

	geomsphere := assetLib.GetGeometry("sphere")
	sphereMat := assetLib.GetMaterial("sphereMat")
	geomsphere.Material = &sphereMat
	sphereNode := renderer.CreateNode()
	sphereNode.Add(&geomsphere)
	sphereNode.SetTranslation(vectormath.Vector3{1, 1, 3})

	//particle effects
	explosionMat := assets.CreateMaterial(assetLib.GetImage("explosion"), nil, nil, nil)
	explosionMat.LightingMode = renderer.MODE_UNLIT
	explosionParticles := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:        4,
		ParticleEmitRate:    2,
		BaseGeometry:        renderer.CreateBox(float32(1), float32(1)),
		Material:            explosionMat,
		TotalFrames:         36,
		FramesX:             6,
		FramesY:             6,
		FaceCamera:          true,
		MaxLife:             1.0,
		MinLife:             2.0,
		StartSize:           vectormath.Vector3{0.4, 0.4, 0.4},
		EndSize:             vectormath.Vector3{2.4, 2.4, 2.4},
		StartColor:          color.NRGBA{254, 254, 254, 254},
		EndColor:            color.NRGBA{254, 254, 254, 254},
		MinTranslation:      vectormath.Vector3{-0.1, -0.1, -0.1},
		MaxTranslation:      vectormath.Vector3{0.1, 0.1, 0.1},
		MaxStartVelocity:    vectormath.Vector3{0.2, 1.8, 0.2},
		MinStartVelocity:    vectormath.Vector3{-0.2, 2.5, -0.2},
		Acceleration:        vectormath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity:  vectormath.IdentityQuaternion(),
		MinAngularVelocity:  vectormath.IdentityQuaternion(),
		MaxRotationVelocity: 0.0,
		MinRotationVelocity: 0.0,
	})

	fireMat := assets.CreateMaterial(assetLib.GetImage("fire"), nil, nil, nil)
	fireMat.LightingMode = renderer.MODE_UNLIT
	fireParticles := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:        10,
		ParticleEmitRate:    2,
		BaseGeometry:        renderer.CreateBox(float32(1), float32(1)),
		Material:            fireMat,
		TotalFrames:         36,
		FramesX:             6,
		FramesY:             6,
		FaceCamera:          true,
		MaxLife:             1.0,
		MinLife:             1.3,
		StartSize:           vectormath.Vector3{1.0, 1.0, 1.0},
		EndSize:             vectormath.Vector3{1.7, 1.7, 1.7},
		StartColor:          color.NRGBA{254, 54, 0, 200},
		EndColor:            color.NRGBA{254, 100, 20, 50},
		MinTranslation:      vectormath.Vector3{-0.1, 0.1, -0.1},
		MaxTranslation:      vectormath.Vector3{0.1, 0.3, 0.1},
		MaxStartVelocity:    vectormath.Vector3{0.02, 0.02, 0.02},
		MinStartVelocity:    vectormath.Vector3{-0.02, -0.02, -0.02},
		Acceleration:        vectormath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity:  vectormath.IdentityQuaternion(),
		MinAngularVelocity:  vectormath.IdentityQuaternion(),
		MaxRotationVelocity: 0.3,
		MinRotationVelocity: -0.3,
	})

	smokeMat := assets.CreateMaterial(assetLib.GetImage("smoke"), nil, nil, nil)
	smokeMat.LightingMode = renderer.MODE_UNLIT
	smokeParticles := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:        100,
		ParticleEmitRate:    35,
		BaseGeometry:        renderer.CreateBox(float32(1), float32(1)),
		Material:            smokeMat,
		TotalFrames:         64,
		FramesX:             8,
		FramesY:             8,
		FaceCamera:          true,
		MaxLife:             2.5,
		MinLife:             2.3,
		StartSize:           vectormath.Vector3{0.4, 0.4, 0.4},
		EndSize:             vectormath.Vector3{2.4, 2.4, 2.4},
		StartColor:          color.NRGBA{254, 254, 254, 50},
		EndColor:            color.NRGBA{254, 254, 254, 0},
		MinTranslation:      vectormath.Vector3{-0.2, -0.2, -0.2},
		MaxTranslation:      vectormath.Vector3{0.2, 0.2, 0.2},
		MaxStartVelocity:    vectormath.Vector3{0.2, 3.3, 0.2},
		MinStartVelocity:    vectormath.Vector3{-0.2, 3.5, -0.2},
		Acceleration:        vectormath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity:  vectormath.IdentityQuaternion(),
		MinAngularVelocity:  vectormath.IdentityQuaternion(),
		MaxRotationVelocity: 0.0,
		MinRotationVelocity: 0.0,
	})

	birdMat := assets.CreateMaterial(assetLib.GetImage("bird"), nil, nil, nil)
	birdMat.LightingMode = renderer.MODE_UNLIT
	birdSprite := effects.CreateSprite(22, 5, 5, &birdMat)
	birdSprite.SetTranslation(vectormath.Vector3{-2, 0, -1})

	sceneGraph := renderer.CreateSceneGraph()
	sceneGraph.AddBackGround(&skyNode)
	sceneGraph.Add(&sphereNode)
	sceneGraph.Add(&fireParticles)
	sceneGraph.Add(&smokeParticles)
	sceneGraph.Add(&explosionParticles)
	sceneGraph.Add(&birdSprite)

	//camera
	camera := renderer.CreateCamera(glRenderer)
	freeMoveActor := actor.CreateFreeMoveActor(camera)
	freeMoveActor.MoveSpeed = 3.0

	glRenderer.Init = func() {
		//setup reflection map
		cubeMap := renderer.CreateCubemap(assetLib.GetMaterial("skyboxMat").Diffuse)
		glRenderer.ReflectionMap(*cubeMap)

		//post effects
		cell := renderer.Shader{
			Name: "shaders/cell/cellCoarse",
		}
		bloomHorizontal := renderer.Shader{
			Name: "shaders/bloom/bloomHorizontal",
			Uniforms: []renderer.Uniform{
				renderer.Uniform{"size", mgl32.Vec2{1900, 1000}},
				renderer.Uniform{"quality", 2.0},
				renderer.Uniform{"samples", 15},
			},
		}
		bloomVertical := renderer.Shader{
			Name: "shaders/bloom/bloomVertical",
			Uniforms: []renderer.Uniform{
				renderer.Uniform{"size", mgl32.Vec2{1900, 1000}},
				renderer.Uniform{"quality", 2.0},
				renderer.Uniform{"samples", 15},
			},
		}

		//input/controller manager
		controllerManager := controller.NewControllerManager(glRenderer.Window)

		//camera free move actor
		mainController := controller.NewBasicMovementController(freeMoveActor)
		controllerManager.AddController(mainController)

		//test the portabitity of the actor / entity interfaces
		mainController.BindAction(func() { freeMoveActor.Entity = camera }, glfw.KeyQ, glfw.Press)
		mainController.BindAction(func() { freeMoveActor.Entity = &sphereNode }, glfw.KeyW, glfw.Press)
		mainController.BindAction(func() { freeMoveActor.Entity = &explosionParticles }, glfw.KeyE, glfw.Press)
		mainController.BindAction(func() { freeMoveActor.Entity = &birdSprite }, glfw.KeyR, glfw.Press)

		mainController.BindAction(func() { //no post effects
			glRenderer.DestroyPostEffects(bloomVertical)
			glRenderer.DestroyPostEffects(bloomHorizontal)
			glRenderer.DestroyPostEffects(cell)
		}, glfw.KeyA, glfw.Press)

		mainController.BindAction(func() { //bloom effect
			glRenderer.CreatePostEffect(bloomVertical)
			glRenderer.CreatePostEffect(bloomHorizontal)
			glRenderer.DestroyPostEffects(cell)
		}, glfw.KeyS, glfw.Press)

		mainController.BindAction(func() { //cell effect
			glRenderer.DestroyPostEffects(bloomVertical)
			glRenderer.DestroyPostEffects(bloomHorizontal)
			glRenderer.CreatePostEffect(cell)
		}, glfw.KeyD, glfw.Press)
	}

	glRenderer.Update = func() {
		fps.UpdateFPSMeter()

		//update things that need updating
		explosionParticles.Update(0.018, glRenderer)
		fireParticles.Update(0.018, glRenderer)
		smokeParticles.Update(0.018, glRenderer)

		birdSprite.NextFrame()

		freeMoveActor.Update(0.018)
	}

	glRenderer.Render = func() {
		sceneGraph.RenderScene(glRenderer)
	}

	glRenderer.Start()
}
