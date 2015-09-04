package renderer

import(
    "github.com/Walesey/goEngine/vectorMath"
    "image"
) 

const (
    MODE_UNLIT int32 = iota
    MODE_LIT
)

type Spatial interface {
    Draw( renderer Renderer )
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

func (geometry *Geometry) Draw( renderer Renderer ) {
    geometry.load( renderer )
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
    scale *vectorMath.Vector3
    translation *vectorMath.Vector3
    orientation *vectorMath.Quaternion 
}

func CreateNode() *Node{
    //create slice to store children
    children := make([]Spatial, 0, 0)
    return &Node{ 
        children: children, 
        Transform: CreateTransform(), 
        scale: &vectorMath.Vector3{1,1,1},
        translation: &vectorMath.Vector3{0,0,0},
        orientation: vectorMath.IdentityQuaternion(),
    }
}

func (node *Node) Draw( renderer Renderer ) {
    renderer.PushTransform()
    if node.Transform != nil{
        renderer.ApplyTransform( node.Transform )
    }
    for _,child := range node.children {
        child.Draw(renderer)
    }
    renderer.PopTransform()
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

func (node *Node) Scale( scale vectorMath.Vector3 ) {
    node.scale.Set( scale )
    node.Transform.From( *node.scale, *node.translation, *node.orientation )
}

func (node *Node) Translation( translation vectorMath.Vector3 ) {
    node.translation.Set( translation )
    node.Transform.From( *node.scale, *node.translation, *node.orientation )
}

func (node *Node) Orientation( orientation vectorMath.Quaternion  ) {
    node.orientation.Set( orientation )
    node.Transform.From( *node.scale, *node.translation, *node.orientation )
}

func (node *Node) Rotation( angle float64, axis vectorMath.Vector3 ) {
    node.orientation.AngleAxis( angle, axis )
    node.Transform.From( *node.scale, *node.translation, *node.orientation )
}


//Primitives
func CreateBox( height, width float32 ) *Geometry {
    verticies := []float32{
        -width/2,0,height/2,  0,1,0, 1,0,-1, -1,0,-1, 0,0, 
        width/2,0,height/2,   0,1,0, 1,0,-1, -1,0,-1, 1,0,
        width/2,0,-height/2,  0,1,0, 1,0,-1, -1,0,-1, 1,1,
        width/2,0,-height/2,  0,1,0, 1,0,-1, -1,0,-1, 1,1, 
        -width/2,0,-height/2, 0,1,0, 1,0,-1, -1,0,-1, 0,1,
        -width/2,0,height/2,  0,1,0, 1,0,-1, -1,0,-1, 0,0,
    }
    indicies := []uint32{0,1,2,3,4,5}
    return CreateGeometry(indicies, verticies)
}
