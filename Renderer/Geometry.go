package Renderer

type Spatial interface {
    draw( renderer Renderer )
}

//Geometry
type Geometry struct {
    vaoId, textureId uint32
    Indicies []uint32
    Verticies []float32
}

func (geometry Geometry) draw( renderer Renderer ) {

}

//Node
type Node struct {
    children []Spatial
}

func CreateNode() Node{
    children := make([]Spatial, 0, 20)
    newNode := Node{ children }
    return newNode
}

func (node Node) draw( renderer Renderer ) {
    for _,child := range node.children {
        child.draw(renderer)
    }
}

// func (nPtr *Node) Add( spatial Spatial ) {
//     node := (*nPtr)
//     length := len(node.children)
//     if length >= cap(node.children) {
//         //expand the slice if necissary
//         node.children = append(node.children, make([]Spatial, 20)...)
//     }
//     //append to the slice
//     node.children = append(node.children, spatial);
// }
//
// func (nPtr *Node) Remove( spatial Spatial ) {
//     node := (*nPtr)
//     //find the address in the slice, remove it and return
//     for index,child := range node.children {
//         if child.getId() == id {
//             node.children = append(node.children[:index], node.children[:index+1]...)
//             break
//         }
//     }
// }
