package main

import (
    "runtime"
    "math"
    "fmt"
    "os"

    "github.com/Walesey/goEngine/vectorMath"
    "github.com/Walesey/goEngine/assets"
    "github.com/Walesey/goEngine/effects"
    "github.com/Walesey/goEngine/renderer"

    "github.com/codegangsta/cli"
)

func init() {
    // GLFW event handling must run on the main OS thread
    runtime.LockOSThread()
    // Use all cpu cores
    runtime.GOMAXPROCS(runtime.NumCPU())
}

func main(){
    app := cli.NewApp()
    app.Name = "goEngine"
    app.Usage = ""
    app.EnableBashCompletion = true
    app.Commands = []cli.Command{
        {
            Name:  "demo",
            Usage: "run a demo program",
            Action: demo,
        },
        {
            Name:  "material",
            Usage: "Import material from file and save to a .asset file",
            Action: materialImport,
        },
        {
            Name:  "image",
            Usage: "Import texture from file and save to a .asset file",
            Action: imageImport,
        },
        {
            Name:  "geometry",
            Usage: "Import obj geometry file and save to a .asset file",
            Action: geometryImport,
        },
        {
            Name:  "list",
            Usage: "List all assets in a .asset file",
            Action: list,
        },
        {
            Name:  "remove",
            Usage: "Remove an asset from a .asset file",
            Action: remove,
        },
    }

    app.Run(os.Args)
}

//CLI remove asset from file
func remove( c *cli.Context ){
    if len(c.Args()) != 2 {
        fmt.Println("Usage: goEngine material <assetFile> <name> ")
        return
    }
    assetLib,_ := assets.LoadAssetLibrary(c.Args()[0])
    delete(assetLib.Assets, c.Args()[1])
    assetLib.SaveToFile( c.Args()[0] )
}

//CLI list assets from file
func list( c *cli.Context ){
    if len(c.Args()) != 1 {
        fmt.Println("Usage: goEngine material <assetFile> ")
        return
    }
    assetLib,_ := assets.LoadAssetLibrary(c.Args()[0])
    for name := range assetLib.Assets {
        fmt.Println( name, ": ", assetLib.Assets[name].Type )
    }
}

//CLI material creator
func materialImport( c *cli.Context ){
    if len(c.Args()) != 6 {
        fmt.Println("Usage: goEngine material <assetFile> <name> <albedoFile> <normalFile> <specFile> <roughnessFile> ")
        return
    }
    diffuseMap := assets.ImportImage(c.Args()[2])
    normalMap := assets.ImportImage(c.Args()[3])
    specMap := assets.ImportImage(c.Args()[4])
    roughnessMap := assets.ImportImage(c.Args()[5])
    mat := assets.CreateMaterial(diffuseMap, normalMap, specMap, roughnessMap)
    assetLib,_ := assets.LoadAssetLibrary(c.Args()[0])
    assetLib.AddMaterial( c.Args()[1], mat )
    assetLib.SaveToFile( c.Args()[0] )
}

//CLI geometry creator
func geometryImport( c *cli.Context ){
    if len(c.Args()) != 3 {
        fmt.Println("Usage: goEngine geometry <assetFile> <name> <objFile> ")
        return
    }
    geometry := assets.ImportObj(c.Args()[2])
    assetLib,_ := assets.LoadAssetLibrary(c.Args()[0])
    assetLib.AddGeometry( c.Args()[1], geometry )
    assetLib.AddMaterial( fmt.Sprint(c.Args()[1], "Mat"), geometry.Material )
    assetLib.SaveToFile( c.Args()[0] )
}

//CLI image creator
func imageImport( c *cli.Context ){
    if len(c.Args()) != 3 {
        fmt.Println("Usage: goEngine image <assetFile> <name> <imageFile>")
        return
    }
    imageAsset := assets.ImportImage(c.Args()[2])
    assetLib,_ := assets.LoadAssetLibrary(c.Args()[0])
    assetLib.AddImage( c.Args()[1], imageAsset )
    assetLib.SaveToFile( c.Args()[0] )
}

//
func demo( c *cli.Context ){
    fps := renderer.CreateFPSMeter(1.0)
    fps.FpsCap = 60

    glRenderer := &renderer.OpenglRenderer{
        WindowTitle : "GoEngine",
        WindowWidth : 900,
        WindowHeight : 700,
    }

    assetLib,err := assets.LoadAssetLibrary("TestAssets/demo.asset")
    if err != nil {
        panic(err)
    }

    //setup scenegraph
    geom := assetLib.GetGeometry("skybox")
    geom.Material = assetLib.GetMaterial("skyboxMat")
    geom.Material.LightingMode = renderer.MODE_UNLIT
    geom.CullBackface = false
    skyNode := renderer.CreateNode()
    skyNode.BucketType = renderer.BUCKET_OPAQUE
    skyNode.Add(&geom)
    skyNode.SetRotation( 1.57, vectorMath.Vector3{0,1,0} )
    skyNode.SetScale( vectorMath.Vector3{5000, 5000, 5000} )

    geomsphere := assetLib.GetGeometry("sphere")
    geomsphere.Material = assetLib.GetMaterial("sphereMat")
    boxNode2 := renderer.CreateNode()
    boxNode2.Add(&geomsphere)

    material := assets.CreateMaterial(assetLib.GetImage("fire"), nil, nil, nil)
    material.LightingMode = renderer.MODE_UNLIT
    firesprite := effects.CreateSprite( 36, 6, 6, material )
    firespriteNode := renderer.CreateNode()
    firespriteNode.Add(&firesprite)

    material = assets.CreateMaterial(assetLib.GetImage("smoke"), nil, nil, nil)
    material.LightingMode = renderer.MODE_UNLIT
    smokesprite := effects.CreateSprite( 64, 8, 8, material )
    smokespriteNode := renderer.CreateNode()
    smokespriteNode.Add(&smokesprite)
    smokespriteNode.SetTranslation( vectorMath.Vector3{-2,0,0} )

    material = assets.CreateMaterial(assetLib.GetImage("explosion"), nil, nil, nil)
    material.LightingMode = renderer.MODE_UNLIT
    explosionsprite := effects.CreateSprite( 36, 6, 6, material )
    explosionspriteNode := renderer.CreateNode()
    explosionspriteNode.Add(&explosionsprite)
    explosionspriteNode.SetTranslation( vectorMath.Vector3{2,0,0} )

    sceneGraph := renderer.CreateSceneGraph()
    sceneGraph.Add(&skyNode)
    sceneGraph.Add(&boxNode2)
    sceneGraph.Add(&firespriteNode)
    sceneGraph.Add(&smokespriteNode)
    sceneGraph.Add(&explosionspriteNode)

    i := -45.0

    glRenderer.Init = func(){
        //setup reflection map
        cubeMap := renderer.CreateCubemap( assetLib.GetMaterial("skyboxMat").Diffuse );
        glRenderer.ReflectionMap( *cubeMap )
    }

    glRenderer.Update = func(){
        fps.UpdateFPSMeter()
        i = i + 0.11
        if i > 180 {
            i = -45
        }
        sine := math.Sin((float64)(i/26))
        cosine := math.Cos((float64)(i/26))

        boxNode2.SetTranslation( vectorMath.Vector3{1, 2, i} )
        //look at the box
        cameraLocation := vectorMath.Vector3{5*cosine,3*sine,5*sine}
        glRenderer.Camera( cameraLocation, vectorMath.Vector3{0,0,0}, vectorMath.Vector3{0,1,0} )

        glRenderer.CreateLight( 5,5,5, 100,100,100, 100,100,100, false, vectorMath.Vector3{1, 2, (float64)(i)}, 1 )

        //face the camera
        firespriteNode.SetFacing( 3.14, glRenderer.CameraLocation().Subtract(smokespriteNode.Translation).Normalize(), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,-1} )
        smokespriteNode.SetFacing( 3.14, glRenderer.CameraLocation().Subtract(smokespriteNode.Translation).Normalize(), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,-1} )
        explosionspriteNode.SetFacing( 3.14, glRenderer.CameraLocation().Subtract(smokespriteNode.Translation).Normalize(), vectorMath.Vector3{0,1,0}, vectorMath.Vector3{0,0,-1} )
    }

    glRenderer.Render = func(){
        sceneGraph.RenderScene(glRenderer)
    }

    glRenderer.Start()
}
