package main

import (
	"runtime"
	"math"

	"goEngine/vectorMath"
	
	"goEngine/renderer"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main(){

    var sceneGraph renderer.SceneGraph
    var mainRenderer renderer.Renderer

    var boxNode *renderer.Node
    var boxNode2 *renderer.Node

    i := (float32)(-30.0)

	mainRenderer = &renderer.OpenglRenderer{
        WindowTitle : "Go Engine",
        WindowWidth : 800,
        WindowHeight : 800,

        Init : func(){
    		sceneGraph = renderer.CreateSceneGraph(mainRenderer)

			geom := renderer.CreateGeometry( cubeIndicies, cubeVertices )
			boxNode = renderer.CreateNode()
			boxNode.Add(geom)
        	sceneGraph.Add(boxNode)

			geom = renderer.CreateGeometry( cubeIndicies, cubeVertices )
			boxNode2 = renderer.CreateNode()
			boxNode2.Add(geom)
        	boxNode.Add(boxNode2)
        },

        Update : func(){
        	i = i + 0.02
        	if i > 130 {
        		i = -30
        	}
        	sine := 5*(float32)(math.Sin((float64)(i/5)))
        	cosine := 5*(float32)(math.Cos((float64)(i/5)))
        	//move the boxes
        	boxNode.Transform = &renderer.GlTransform{ mgl32.Translate3D(i ,0, 0) }
        	boxNode2.Transform = &renderer.GlTransform{ mgl32.Translate3D(0 , sine, cosine) }
        	//look at the box
        	mainRenderer.Camera( vectorMath.Vector3{9,9,9}, vectorMath.Vector3{(float64)(i) ,0, 0}, vectorMath.Vector3{0,1,0} )
        },

        Render : func(){
        	sceneGraph.RenderScene()
        }}

     mainRenderer.Start()
}


var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}

var cubeIndicies = []uint32{
//top
	1, 2, 3,
	3, 2, 4,

//bottom
	5, 6, 7,
	7, 6, 8,
}

// //TEST
// var cubeVertices = []float32{
// 	//  X, Y, Z, U, V
//
// 	-1.0, -1.0, -1.0, 0.0, 0.0,
// 	1.0, -1.0, -1.0, 1.0, 0.0,
// 	-1.0, 1.0, -1.0, 0.0, 1.0,
// 	1.0, 1.0, -1.0, 1.0, 1.0,
//
//     -1.0, -1.0, 1.0, 0.0, 0.0,
//     1.0, -1.0, 1.0, 1.0, 0.0,
//     -1.0, 1.0, 1.0, 0.0, 1.0,
//     1.0, 1.0, 1.0, 1.0, 1.0,
// }
