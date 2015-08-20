package renderer

import (
    "goEngine/vectorMath"
)

type Spatial interface {
    load( renderer Renderer )
    draw( renderer Renderer )
}

//Geometry
type Geometry struct {
    vaoId, textureId uint32
    loaded bool
    Indicies []uint32
    Verticies []float32
}

func CreateGeometry( indicies []uint32, verticies []float32 ) Geometry {
    return Geometry{ Indicies : indicies, Verticies : verticies, loaded : false }
}

func (geometry *Geometry) draw( renderer Renderer ) {
    renderer.DrawGeometry( geometry )
}

func (geometry *Geometry) load( renderer Renderer ) {
    if !geometry.loaded {
        renderer.CreateGeometry( geometry )
        geometry.loaded = true
    }
}

//Node
type Node struct {
    children []Spatial
    location vectorMath.Vector
    orientation vectorMath.Quaternion
}

func CreateNode() Node{
    //create slice to store children
    children := make([]Spatial, 0, 20)
    newNode := Node{ children: children }
    return newNode
}

func (node *Node) draw( renderer Renderer ) {
    for _,child := range node.children {
        child.draw(renderer)
    }
}

func (node *Node) load( renderer Renderer ) bool {
    for _,child := range node.children {
        child.load(renderer)
    }
    return false
}

func (node *Node) Add( spatial Spatial ) {
    length := len(node.children)
    if length >= cap(node.children) {
        //expand the slice if necissary
        node.children = append(node.children, make([]Spatial, 20)...)
    }
    //append to the slice
    node.children = append(node.children, spatial);
}

func (node *Node) Remove( spatial Spatial ) {
    //find the address in the slice, remove it and return
    for index,child := range node.children {
        if child == spatial {
            node.children = append(node.children[:index], node.children[:index+1]...)
            break
        }
    }
}