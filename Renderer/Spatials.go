package renderer

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

//vericies format : x,y,z,   nx,ny,nz,   u,v
//indicies format : f1,f2,f3 (triangles)
func CreateGeometry( indicies []uint32, verticies []float32 ) *Geometry {
    return &Geometry{ Indicies : indicies, Verticies : verticies, loaded : false }
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
            node.children = append(node.children[:index], node.children[:index+1]...)
            break
        }
    }
}
