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
	fps.FpsCap = 60

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

	//geom := assetLib.GetGeometry("nightskybox")
	//skyboxMat := assetLib.GetMaterial("nightskyboxMat")
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
	sphereNode.SetTranslation(vectormath.Vector3{0, 0, 0})

	//particle effects
	explosionMat := assets.CreateMaterial(assetLib.GetImage("explosion"), nil, nil, nil)
	explosionMat.LightingMode = renderer.MODE_UNLIT
	explosionMat.Transparency = renderer.TRANSPARENCY_EMISSIVE
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
	fireMat.Transparency = renderer.TRANSPARENCY_EMISSIVE
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
		MaxParticles:        38,
		ParticleEmitRate:    15,
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
		MaxStartVelocity:    vectormath.Vector3{0.2, 0.8, 0.2},
		MinStartVelocity:    vectormath.Vector3{-0.2, 0.6, -0.2},
		Acceleration:        vectormath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity:  vectormath.IdentityQuaternion(),
		MinAngularVelocity:  vectormath.IdentityQuaternion(),
		MaxRotationVelocity: 0.0,
		MinRotationVelocity: 0.0,
	})

	sparkMat := assets.CreateMaterial(assetLib.GetImage("spark"), nil, nil, nil)
	sparkMat.LightingMode = renderer.MODE_EMIT
	sparkMat.Transparency = renderer.TRANSPARENCY_EMISSIVE
	sparkParticles := effects.CreateParticleSystem(effects.ParticleSettings{
		MaxParticles:        100,
		ParticleEmitRate:    110,
		BaseGeometry:        renderer.CreateBox(float32(1), float32(1)),
		Material:            sparkMat,
		TotalFrames:         1,
		FramesX:             1,
		FramesY:             1,
		FaceCamera:          true,
		MaxLife:             0.9,
		MinLife:             0.7,
		StartSize:           vectormath.Vector3{0.02, 0.02, 0.02},
		EndSize:             vectormath.Vector3{0.02, 0.02, 0.02},
		StartColor:          color.NRGBA{255, 5, 5, 255},
		EndColor:            color.NRGBA{255, 5, 5, 255},
		MinTranslation:      vectormath.Vector3{0, -0, 0},
		MaxTranslation:      vectormath.Vector3{0, -0, 0},
		MaxStartVelocity:    vectormath.Vector3{0.6, 0.3, 0.6},
		MinStartVelocity:    vectormath.Vector3{-0.6, 0.5, -0.6},
		Acceleration:        vectormath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity:  vectormath.IdentityQuaternion(),
		MinAngularVelocity:  vectormath.IdentityQuaternion(),
		MaxRotationVelocity: 0.0,
		MinRotationVelocity: 0.0,
	})
	fireParticles.Location = vectormath.Vector3{2, 0, -2}
	smokeParticles.Location = vectormath.Vector3{-2, 0, 2}
	explosionParticles.Location = vectormath.Vector3{-2, 0, -2}
	sparkParticles.Location = vectormath.Vector3{2, 0, 2}

	birdMat := assets.CreateMaterial(assetLib.GetImage("bird"), nil, nil, nil)
	birdMat.LightingMode = renderer.MODE_UNLIT
	birdSprite := effects.CreateSprite(22, 5, 5, &birdMat)
	birdSprite.SetTranslation(vectormath.Vector3{-2, 0, -1})

	sceneGraph := renderer.CreateSceneGraph()
	sceneGraph.AddBackGround(skyNode)
	sceneGraph.Add(sphereNode)
	sceneGraph.Add(&fireParticles)
	sceneGraph.Add(&smokeParticles)
	sceneGraph.Add(&explosionParticles)
	sceneGraph.Add(&birdSprite)
	sceneGraph.Add(&sparkParticles)

	//camera
	camera := renderer.CreateCamera(glRenderer)
	freeMoveActor := actor.CreateFreeMoveActor(camera)
	freeMoveActor.MoveSpeed = 3.0
	freeMoveActor.Location = vectormath.Vector3{-2, 0, 0}

	glRenderer.Init = func() {
		//Lighting
		glRenderer.CreateLight(0.1, 0.1, 0.1, 1, 1, 1, 1, 1, 1, true, vectormath.Vector3{0, -1, 0}, 0)

		//setup reflection map
		//cubeMap := renderer.CreateCubemap(assetLib.GetMaterial("nightskyboxMat").Diffuse)
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

		//input/controller manager
		controllerManager := controller.NewControllerManager(glRenderer.Window)

		//lock the cursor
		glRenderer.LockCursor(true)

		//camera free move actor
		mainController := controller.NewBasicMovementController(freeMoveActor)
		controllerManager.AddController(mainController)

		customController := controller.NewActionMap()
		controllerManager.AddController(customController)

		//close window and exit on escape
		customController.BindAction(func() {
			glRenderer.Window.SetShouldClose(true)
		}, glfw.KeyEscape, glfw.Press)

		//test the portabitity of the actor / entity interfaces
		customController.BindAction(func() { freeMoveActor.Entity = camera }, glfw.KeyQ, glfw.Press)
		customController.BindAction(func() { freeMoveActor.Entity = sphereNode }, glfw.KeyW, glfw.Press)
		customController.BindAction(func() { freeMoveActor.Entity = &explosionParticles }, glfw.KeyE, glfw.Press)
		customController.BindAction(func() { freeMoveActor.Entity = &birdSprite }, glfw.KeyR, glfw.Press)

		customController.BindAction(func() { //no post effects
			glRenderer.DestroyPostEffects(bloomVertical)
			glRenderer.DestroyPostEffects(bloomHorizontal)
			glRenderer.DestroyPostEffects(cell)
		}, glfw.KeyA, glfw.Press)

		customController.BindAction(func() { //bloom effect
			glRenderer.CreatePostEffect(bloomVertical)
			glRenderer.CreatePostEffect(bloomHorizontal)
			glRenderer.DestroyPostEffects(cell)
		}, glfw.KeyS, glfw.Press)

		customController.BindAction(func() { //cell effect
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
		sparkParticles.Update(0.018, glRenderer)

		birdSprite.NextFrame()

		freeMoveActor.Update(0.018)
	}

	glRenderer.Render = func() {
		sceneGraph.RenderScene(glRenderer)
	}

	glRenderer.Start()
}
