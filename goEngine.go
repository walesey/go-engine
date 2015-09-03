package main

import (
	"runtime"
    "image"
	"math"
    "fmt"
    "os"

	"github.com/Walesey/goEngine/vectorMath"
    "github.com/Walesey/goEngine/assets"
    "github.com/Walesey/goEngine/assetEncoding"
	"github.com/Walesey/goEngine/renderer"

    "github.com/codegangsta/cli"
	"github.com/go-gl/mathgl/mgl32"
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
            Usage: "import textures from file and save to a .asset file",
            Action: material,
        },
        {
            Name:  "geometry",
            Usage: "import obj geometry file and save to a .asset file",
            Action: geometry,
        },
    }

    app.Run(os.Args)
}

//CLI material creator
func material( c *cli.Context ){
    if len(c.Args()) != 6 {
        fmt.Println("Usage: goEngine material <assetFile> <name> <albedoFile> <normalFile> <specFile> <roughnessFile> ")
        return
    }
    diffuseMap := assets.ImportImage(c.Args()[2])
    normalMap := assets.ImportImage(c.Args()[3])
    specMap := assets.ImportImage(c.Args()[4])
    roughnessMap := assets.ImportImage(c.Args()[5])
    mat := createMaterial(diffuseMap, normalMap, specMap, roughnessMap)
    assetLib,_ := assetEncoding.LoadAssetLibrary(c.Args()[0])
    assetLib.AddMaterial( c.Args()[1], mat )
    assetLib.SaveToFile( c.Args()[0] )
}

//CLI geometry creator
func geometry( c *cli.Context ){
    if len(c.Args()) != 3 {
        fmt.Println("Usage: goEngine geometry <assetFile> <name> <objFile> ")
        return
    }
    geometry := createGeometry(c.Args()[2])
    assetLib,_ := assetEncoding.LoadAssetLibrary(c.Args()[0])
    assetLib.AddGeometry( c.Args()[1], geometry )
    assetLib.AddMaterial( fmt.Sprint(c.Args()[1], "Mat"), geometry.Material )
    assetLib.SaveToFile( c.Args()[0] )
}

// load obj file and return Geometry
func createGeometry( filePath string ) *renderer.Geometry{
    objData,_ := assets.ImportObj(filePath)
    mat := createMaterial(objData.Mtl.Map_Kd, objData.Mtl.Map_Disp, objData.Mtl.Map_Spec, objData.Mtl.Map_Roughness)
    geom := renderer.CreateGeometry( objData.Indicies, objData.Vertices )
    geom.Material = mat
    return geom
}

// Create material object from image files
func createMaterial( diffuse, normal, specular, roughness image.Image ) *renderer.Material{
    mat := renderer.CreateMaterial()
    mat.Diffuse = diffuse
    mat.Normal = normal
    mat.Specular = specular
    mat.Roughness = roughness
    return mat
}

//
func demo( c *cli.Context ){
    fps := renderer.CreateFPSMeter(1.0)
    fps.FpsCap = 60

    var sceneGraph renderer.SceneGraph
    var mainRenderer renderer.Renderer

    var boxNode *renderer.Node
    var boxNode2 *renderer.Node

    assetLib,err := assetEncoding.LoadAssetLibrary("TestAssets/demo.asset")
    if err != nil {
        panic(err)
    }

    i := (float32)(-45.0)

	mainRenderer = &renderer.OpenglRenderer{
        WindowTitle : "Go Engine",
        WindowWidth : 1700,
        WindowHeight : 950,

        Init : func(){
    		sceneGraph = renderer.CreateSceneGraph(mainRenderer)

            //setup reflection map
            cubeMap := renderer.CreateCubemap(assets.ImportImage("TestAssets/skybox/cubemap.png"));
            mainRenderer.ReflectionMap( *cubeMap )

            //setup scenegraph
            geom := assetLib.GetGeometry("skybox")
            geom.Material = assetLib.GetMaterial("skyboxMat")
            geom.Material.LightingMode = renderer.MODE_UNLIT
            geom.CullBackface = false
            skyNode := renderer.CreateNode()
            skyNode.Add(geom)
            sceneGraph.Add(skyNode)
            skyNode.Transform = &renderer.GlTransform{ mgl32.Scale3D(5000, 5000, 5000).Mul4(mgl32.HomogRotate3DY(1.57)) }

			geom = assetLib.GetGeometry("sphere")
            geom.Material = assetLib.GetMaterial("sphereMat")
			boxNode = renderer.CreateNode()
			boxNode.Add(geom)
        	sceneGraph.Add(boxNode)

            geom = assetLib.GetGeometry("sphere")
            geom.Material = assetLib.GetMaterial("sphereMat")
			boxNode2 = renderer.CreateNode()
			boxNode2.Add(geom)
        	sceneGraph.Add(boxNode2)
        },

        Update : func(){
            fps.UpdateFPSMeter()
        	i = i + 0.08
        	if i > 180 {
        		i = -45
        	}
        	sine := math.Sin((float64)(i/26))
        	cosine := math.Cos((float64)(i/26))

        	boxNode.Transform = &renderer.GlTransform{ mgl32.Translate3D(0 , 0, 0).Mul4(mgl32.HomogRotate3DY(1.57))  }
            boxNode2.Transform = &renderer.GlTransform{ mgl32.Translate3D(1, 2, i) }
        	//look at the box
        	mainRenderer.Camera( vectorMath.Vector3{5*cosine,1*sine,5*sine}, vectorMath.Vector3{0,0,0}, vectorMath.Vector3{0,1,0} )

            mainRenderer.CreateLight( 5,5,5, 100,100,100, 100,100,100, false, vectorMath.Vector3{1, 2, (float64)(i)}, 1 )
        },

        Render : func(){
        	sceneGraph.RenderScene()
        }}

    mainRenderer.Start()
}