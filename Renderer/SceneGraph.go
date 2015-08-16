package Renderer

import (
    "fmt"
)

//Interfaces
type SceneGraph interface {
    RenderScene()
    Add( spatial Spatial )
}

//DefaultSceneGraph
type DefaultSceneGraph struct {
    SceneRenderer Renderer
    root Node
}

func CreateSceneGraph( renderer Renderer ) SceneGraph{
    return DefaultSceneGraph{ renderer, CreateNode() }
}

func (sceneGraph DefaultSceneGraph) RenderScene() {
    // sceneGraph.root.load( sceneGraph.SceneRenderer )
    sceneGraph.root.draw( sceneGraph.SceneRenderer )
}

func (sceneGraph DefaultSceneGraph) GetRootNode() Node {
    return sceneGraph.root
}

func (sceneGraph DefaultSceneGraph) Add( spatial Spatial) {
    // sceneGraph.root.Add(spatial)
    fmt.Println(sceneGraph.root.children)
}
