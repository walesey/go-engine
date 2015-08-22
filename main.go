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

    i := (float32)(-5.0)

	mainRenderer = &renderer.OpenglRenderer{
        WindowTitle : "Go Engine",
        WindowWidth : 2400,
        WindowHeight : 1200,

        Init : func(){
    		sceneGraph = renderer.CreateSceneGraph(mainRenderer)

            hulk := assets.ImportObj("TestAssets/faces.obj")
            ares := assets.ImportObj("TestAssets/faces.obj")

			geom := renderer.CreateGeometry( hulk.Indicies, hulk.Vertices )
			boxNode = renderer.CreateNode()
			boxNode.Add(geom)
        	sceneGraph.Add(boxNode)

			geom = renderer.CreateGeometry( ares.Indicies, ares.Vertices )
			boxNode2 = renderer.CreateNode()
			boxNode2.Add(geom)
        	boxNode.Add(boxNode2)
        },

        Update : func(){
            fps.UpdateFPSMeter()
        	i = i + 0.01
        	if i > 15 {
        		i = -5
        	}
        	sine := 2*(float32)(math.Sin((float64)(i/6)))
        	cosine := 2*(float32)(math.Cos((float64)(i/6)))
        	//move the boxes
        	boxNode.Transform = &renderer.GlTransform{ mgl32.Translate3D(0 , 0, i) }
        	boxNode2.Transform = &renderer.GlTransform{ mgl32.Translate3D(cosine, sine, 0) }
        	//look at the box
        	mainRenderer.Camera( vectorMath.Vector3{2,2,2}, vectorMath.Vector3{0 , 0, (float64)(i)}, vectorMath.Vector3{0,1,0} )
        },

        Render : func(){
        	sceneGraph.RenderScene()
        }}

    mainRenderer.Start()
}