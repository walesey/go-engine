package renderer

import(
    "image"
    "image/color"

    "github.com/Walesey/goEngine/vectorMath"
)

const (
    MODE_UNLIT int32 = iota
    MODE_LIT
)

//A Spatial is something that can be Drawn by a Renderer
type Spatial interface {
    Draw( renderer Renderer )
}

//An Entity is something that can be scaled, positioned and rotated (orientation)
type Entity interface {
    SetScale( scale vectorMath.Vector3 )
    SetTranslation( translation vectorMath.Vector3 ) 
    SetOrientation( orientation vectorMath.Quaternion )
}

type Material struct {
    diffuseId, normalId, specularId, glossId, roughnessId uint32
    loaded bool
    LightingMode int32
    Diffuse, Normal, Specular, Roughness image.Image
}

func CreateMaterial() Material {
    return Material{ loaded : false, LightingMode : MODE_LIT }
}

type Flipbook struct {
    IndexX, IndexY int
    FrameSizeX, FrameSizeY float32
}

//Geometry
type Geometry struct {
    vboId, iboId uint32
    loaded bool
    Indicies []uint32
    Verticies []float32
    Material *Material
    CullBackface bool
    Flipbook Flipbook
    Color color.NRGBA
}

//vericies format : x,y,z,   nx,ny,nz,tx,ty,tz,btx,bty,btz,   u,v
//indicies format : f1,f2,f3 (triangles)
func CreateGeometry( indicies []uint32, verticies []float32 ) Geometry {
    mat := CreateMaterial()
    return Geometry{ 
        Indicies : indicies, 
        Verticies : verticies, 
        Material: &mat, 
        loaded : false, 
        CullBackface : true, 
        Flipbook: Flipbook{0, 0, 1.0, 1.0},
        Color: color.NRGBA{255,255,255,255},
    }
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
    Scale vectorMath.Vector3
    Translation vectorMath.Vector3
    Orientation vectorMath.Quaternion
}

func CreateNode() Node{
    //create slice to store children
    children := make([]Spatial, 0, 0)
    return Node{
        children: children,
        Transform: CreateTransform(),
        Scale: vectorMath.Vector3{1,1,1},
        Translation: vectorMath.Vector3{0,0,0},
        Orientation: vectorMath.IdentityQuaternion(),
    }
}

func (node *Node) Draw( renderer Renderer ) {
    renderer.PushTransform()
    renderer.ApplyTransform( node.Transform )
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

func (node *Node) SetScale( scale vectorMath.Vector3 ) {
    node.Scale = scale
    node.Transform.From( node.Scale, node.Translation, node.Orientation )
}

func (node *Node) SetTranslation( translation vectorMath.Vector3 ) {
    node.Translation = translation
    node.Transform.From( node.Scale, node.Translation, node.Orientation )
}

func (node *Node) SetOrientation( orientation vectorMath.Quaternion  ) {
    node.Orientation = orientation
    node.Transform.From( node.Scale, node.Translation, node.Orientation )
}

func (node *Node) SetRotation( angle float64, axis vectorMath.Vector3 ) {
    node.Orientation = vectorMath.AngleAxis( angle, axis )
    node.Transform.From( node.Scale, node.Translation, node.Orientation )
}

//used for eg. sprites facing the direction of the camera - all vectors need to be normalized
func (node *Node) SetFacing( rotation float64, newNormal, normal, tangent vectorMath.Vector3 ) {
    angleCorrection := -tangent.AngleBetween( newNormal.Subtract(newNormal.Project(normal)) )
    if normal.Cross(tangent).Dot(newNormal) < 0 {
        angleCorrection = -angleCorrection
    }
    angleQ := vectorMath.AngleAxis( rotation + angleCorrection, normal )
    betweenVectorsQ := vectorMath.BetweenVectors( normal, newNormal ) 
    node.Orientation = betweenVectorsQ.Multiply(angleQ)
    node.Transform.From( node.Scale, node.Translation, node.Orientation )
}

//Primitives
func CreateBox( height, width float32 ) Geometry {
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
