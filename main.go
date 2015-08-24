package main

import (
	"runtime"
	"math"

	"goEngine/vectorMath"
    "goEngine/assets"

	"goEngine/renderer"
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

            hulk,_ := assets.ImportObj("TestAssets/gun/rifle.obj")
            hulkMat := renderer.CreateMaterial()
            hulkMat.Diffuse = hulk.Mtl.Map_Kd
            hulkMat.Normal = hulk.Mtl.Map_Disp
            hulkMat.Specular = hulk.Mtl.Map_Spec
            hulkMat.Roughness = hulk.Mtl.Map_Roughness

            ares,_ := assets.ImportObj("TestAssets/alarm/alarm.obj")
            aresMat := renderer.CreateMaterial()
            aresMat.Diffuse = ares.Mtl.Map_Kd
            aresMat.Normal = ares.Mtl.Map_Disp
            aresMat.Specular = ares.Mtl.Map_Spec
            aresMat.Roughness = ares.Mtl.Map_Roughness

			geom := renderer.CreateGeometry( hulk.Indicies, hulk.Vertices )
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
        	i = i + 0.04
        	if i > 80 {
        		i = -45
        	}
        	sine := math.Sin((float64)(i/36))
        	cosine := math.Cos((float64)(i/36))

        	boxNode.Transform = &renderer.GlTransform{ mgl32.Translate3D(0 , 0, 5).Mul4(mgl32.HomogRotate3DY(1.57))  }
            boxNode2.Transform = &renderer.GlTransform{ mgl32.Translate3D(1, 3, i) }
        	//look at the box
        	mainRenderer.Camera( vectorMath.Vector3{10*cosine,2,10*sine}, vectorMath.Vector3{0 , 0, (float64)(5-5)}, vectorMath.Vector3{0,1,0} )

            mainRenderer.CreateLight( 5,5,5,   160,0,0,   120,0,0,   vectorMath.Vector3{1, 3, (float64)(i)}, 1 )
        },

        Render : func(){
        	sceneGraph.RenderScene()
        }}

    mainRenderer.Start()
}
