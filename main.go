package main

import (
	"runtime"
	"math"

	"github.com/Walesey/goEngine/vectorMath"
    "github.com/Walesey/goEngine/assets"
	"github.com/Walesey/goEngine/renderer"

	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
    //Use all cpu cores
    runtime.GOMAXPROCS(runtime.NumCPU())
}

func main(){

    fps := renderer.CreateFPSMeter(1.0)
    fps.FpsCap = 60

    var sceneGraph renderer.SceneGraph
    var mainRenderer renderer.Renderer

    var boxNode *renderer.Node
    var boxNode2 *renderer.Node

    i := (float32)(-45.0)

	mainRenderer = &renderer.OpenglRenderer{
        WindowTitle : "Go Engine",
        WindowWidth : 1700,
        WindowHeight : 950,

        Init : func(){
    		sceneGraph = renderer.CreateSceneGraph(mainRenderer)

            //setup reflection map
            cubeMap := assets.ImportCubemap("TestAssets/skybox/cubemap.jpg")
            mainRenderer.ReflectionMap( cubeMap.Right, cubeMap.Left, cubeMap.Top, cubeMap.Bottom, cubeMap.Back, cubeMap.Front )

            //setup skybox
            skycube,_ := assets.ImportObj("TestAssets/skybox/skybox.obj")
            skyMat := renderer.CreateMaterial()
            skyMat.Diffuse = skycube.Mtl.Map_Kd
            skyMat.LightingMode = renderer.MODE_UNLIT

            hulk,_ := assets.ImportObj("TestAssets/sphere1.obj")
            // hulk,_ := assets.ImportObj("TestAssets/gun/rifle.obj")
            hulkMat := renderer.CreateMaterial()
            hulkMat.Diffuse = hulk.Mtl.Map_Kd
            hulkMat.Normal = hulk.Mtl.Map_Disp
            hulkMat.Specular = hulk.Mtl.Map_Spec
            hulkMat.Roughness = hulk.Mtl.Map_Roughness

            ares,_ := assets.ImportObj("TestAssets/sphere1.obj")
            aresMat := renderer.CreateMaterial()
            aresMat.Diffuse = ares.Mtl.Map_Kd
            aresMat.Normal = ares.Mtl.Map_Disp
            aresMat.Specular = ares.Mtl.Map_Spec
            aresMat.Roughness = ares.Mtl.Map_Roughness

            //setup scenegraph
            geom := renderer.CreateGeometry( skycube.Indicies, skycube.Vertices )
            geom.Material = skyMat
            skyNode := renderer.CreateNode()
            skyNode.Add(geom)
            sceneGraph.Add(skyNode)
            skyNode.Transform = &renderer.GlTransform{ mgl32.Scale3D(5000, 5000, 5000) }

			geom = renderer.CreateGeometry( hulk.Indicies, hulk.Vertices )
            geom.Material = hulkMat
			boxNode = renderer.CreateNode()
			boxNode.Add(geom)
        	sceneGraph.Add(boxNode)

			geom = renderer.CreateGeometry( ares.Indicies, ares.Vertices )
            geom.Material = aresMat
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
