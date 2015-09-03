package renderer

import "image"

const (
    MODE_UNLIT int32 = 0 + iota
    MODE_LIT
)

type Spatial interface {
    load( renderer Renderer )
    draw( renderer Renderer )
}

type Material struct {
    diffuseId, normalId, specularId, glossId, roughnessId uint32
    loaded bool
    LightingMode int32
    Diffuse, Normal, Specular, Roughness image.Image
}

func CreateMaterial() *Material {
    return &Material{ loaded : false, LightingMode : MODE_LIT }
}

//Geometry
type Geometry struct {
    vboId, iboId uint32
    loaded bool
    Indicies []uint32
    Verticies []float32
    Material *Material
    CullBackface bool
}


//vericies format : x,y,z,   nx,ny,nz,tx,ty,tz,btx,bty,btz,   u,v
//indicies format : f1,f2,f3 (triangles)
func CreateGeometry( indicies []uint32, verticies []float32 ) *Geometry {

    return &Geometry{ Indicies : indicies, Verticies : verticies, Material: CreateMaterial(), loaded : false, CullBackface : true }
}

func (geometry *Geometry) draw( renderer Renderer ) {
    renderer.DrawGeometry( geometry )
}

func (geometry *Geometry) load( renderer Renderer ) {
    if !geometry.loaded {
        renderer.CreateGeometry( geometry )
        geometry.loaded = true
    }
    if !geometry.Material.loaded {
        renderer.CreateMaterial( geometry.Material )
        geometry.Material.loaded = true
    }
}


//Node
type Node struct {
    children []Spatial
    Transform Transform
}

func CreateNode() *Node{
    //create slice to store children
    children := make([]Spatial, 0, 0)
    return &Node{ children: children }
}

func (node *Node) draw( renderer Renderer ) {
    renderer.PushTransform()
    if node.Transform != nil{
        renderer.ApplyTransform( node.Transform )
    }
    for _,child := range node.children {
        child.draw(renderer)
    }
    renderer.PopTransform()
}

func (node *Node) load( renderer Renderer ) {
    for _,child := range node.children {
        child.load(renderer)
    }
}

func (node *Node) Add( spatial Spatial ) {
    //append to the slice
    node.children = append(node.children, spatial)
}

func (node *Node) Remove( spatial Spatial ) {
    //find the address in the slice, remove it and return
    for index,child := range node.children {
        if child == spatial {
            node.children = append(node.children[:index], node.children[index+1:]...)
            break
        }
    }
}
