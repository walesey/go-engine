package examples

import (
	"image/color"

    "github.com/Walesey/goEngine/vectorMath"
    "github.com/Walesey/goEngine/assets"
    "github.com/Walesey/goEngine/effects"
    "github.com/Walesey/goEngine/renderer"
    "github.com/Walesey/goEngine/controller"
    "github.com/Walesey/goEngine/actor"

    "github.com/codegangsta/cli"
    "github.com/go-gl/glfw/v3.1/glfw"
)


//
func Particles( c *cli.Context ){
    fps := renderer.CreateFPSMeter(1.0)
    fps.FpsCap = 60

    glRenderer := &renderer.OpenglRenderer{
        WindowTitle : "GoEngine",
        WindowWidth : 1900,
        WindowHeight : 1000,
    }

    assetLib,err := assets.LoadAssetLibrary("TestAssets/demo.asset")
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
    skyNode.BucketType = renderer.BUCKET_OPAQUE
    skyNode.Add(&geom)
    skyNode.SetRotation( 1.57, vectorMath.Vector3{0,1,0} )
    skyNode.SetScale( vectorMath.Vector3{5000, 5000, 5000} )

    geomsphere := assetLib.GetGeometry("sphere")
    sphereMat := assetLib.GetMaterial("sphereMat")
    geomsphere.Material = &sphereMat
    sphereNode := renderer.CreateNode()
    sphereNode.Add(&geomsphere)
    sphereNode.SetTranslation( vectorMath.Vector3{1,1,3} )

    fireMat := assets.CreateMaterial(assetLib.GetImage("fire"), nil, nil, nil)
    fireMat.LightingMode = renderer.MODE_UNLIT
    firesprite := effects.CreateSprite( 36, 6, 6, &fireMat )
    firespriteNode := renderer.CreateNode()
    firespriteNode.Add(&firesprite)

    smokeMat := assets.CreateMaterial(assetLib.GetImage("smoke"), nil, nil, nil)
    smokeMat.LightingMode = renderer.MODE_UNLIT
    smokesprite := effects.CreateSprite( 64, 8, 8, &smokeMat )
    //smoke particle effect
    smokeParticles := effects.CreateParticleSystem( effects.ParticleSettings{
    	MaxParticles: 200,
		ParticleEmitRate: 30,
		Sprite: smokesprite,
		FaceCamera: true,
		MaxLife: 5.0,
		MinLife: 7.0,
		StartSize: vectorMath.Vector3{0.4, 0.4, 0.4},
		EndSize: vectorMath.Vector3{2.4, 2.4, 2.4},
		StartColor: color.NRGBA{254, 254, 254, 254},
		EndColor: color.NRGBA{254, 254, 254, 0},
		MinTranslation: vectorMath.Vector3{-0.2, -0.2, -0.2},
		MaxTranslation: vectorMath.Vector3{0.2, 0.2, 0.2},
		MaxStartVelocity: vectorMath.Vector3{-0.2, 0.3, 0.2},
		MinStartVelocity: vectorMath.Vector3{-0.2, 0.5, 0.2},
		Acceleration: vectorMath.Vector3{0.0, 0.0, 0.0},
		MaxAngularVelocity: vectorMath.IdentityQuaternion(),
		MinAngularVelocity: vectorMath.IdentityQuaternion(),
		MaxRotationVelocity: 0.0,
		MinRotationVelocity: 0.0,
    })

    explosionMat := assets.CreateMaterial(assetLib.GetImage("explosion"), nil, nil, nil)
    explosionMat.LightingMode = renderer.MODE_UNLIT
    explosionsprite := effects.CreateSprite( 36, 6, 6, &explosionMat )
    explosionspriteNode := renderer.CreateNode()
    explosionspriteNode.Add(&explosionsprite)
    explosionspriteNode.SetTranslation( vectorMath.Vector3{2,0,0} )

    sceneGraph := renderer.CreateSceneGraph()
    sceneGraph.Add(&skyNode)
    sceneGraph.Add(&sphereNode)
    sceneGraph.Add(&firespriteNode)
    sceneGraph.Add(&smokeParticles.Node)
    sceneGraph.Add(&explosionspriteNode)

    //camera
    camera := renderer.CreateCamera(glRenderer)
    freeMoveActor := actor.CreateFreeMoveActor( camera )
    freeMoveActor.MoveSpeed = 3.0

    glRenderer.Init = func(){
        //setup reflection map
        cubeMap := renderer.CreateCubemap( assetLib.GetMaterial("skyboxMat").Diffuse )
        glRenderer.ReflectionMap( *cubeMap )

        //input/controller manager
        controllerManager := controller.NewControllerManager(glRenderer.Window)

        //camera free move actor
        mainController := controller.NewBasicMovementController(freeMoveActor)
        controllerManager.AddController( mainController )

        //test the portabitity of the actor / entity interfaces 
        mainController.BindAction(func(){ freeMoveActor.Entity = camera }, glfw.KeyQ, glfw.Press)
        mainController.BindAction(func(){ freeMoveActor.Entity = &sphereNode }, glfw.KeyW, glfw.Press)
        mainController.BindAction(func(){ freeMoveActor.Entity = &smokeParticles }, glfw.KeyE, glfw.Press)
    }

    glRenderer.Update = func(){
        fps.UpdateFPSMeter()

        //face the camera
        firespriteNode.SetFacing( 3.14, glRenderer.CameraLocation().Subtract(firespriteNode.Translation).Normalize(), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,-1} )
        explosionspriteNode.SetFacing( 3.14, glRenderer.CameraLocation().Subtract(explosionspriteNode.Translation).Normalize(), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,-1} )

        //update things that need updating
        firesprite.NextFrame()
        explosionsprite.NextFrame()
        smokeParticles.Update(0.018, glRenderer)

        freeMoveActor.Update(0.018)
    }

    glRenderer.Render = func(){
        sceneGraph.RenderScene(glRenderer)
    }

    glRenderer.Start()
}
