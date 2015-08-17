package main

import (
    "runtime"

    "goEngine/Renderer"
)

// var sceneGraph Renderer.SceneGraph

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main(){
    // sceneGraph = Renderer.CreateSceneGraph(renderer)
    // sceneGraph.GetRootNode().Add(geom)
    var geom Renderer.Geometry

    renderer := Renderer.OpenglRenderer{
        WindowTitle : "Go Engine",
        WindowWidth : 800,
        WindowHeight : 800,
        Init : func(renderer Renderer.Renderer){
            geom = Renderer.Geometry{ Indicies : cubeIndicies, Verticies : cubeVertices  }
            renderer.CreateGeometry(&geom)
        },
        Update : func(renderer Renderer.Renderer){

        },
        Render : func(renderer Renderer.Renderer){
            renderer.DrawGeometry(&geom)
        }}
     renderer.Start();
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
