package renderer

import(
    "github.com/Walesey/goEngine/vectorMath"
)

const (
    MODE_UNLIT int32 = iota
    MODE_LIT
)

//A Spatial is something that can be Drawn by a Renderer
type Spatial interface {
    Draw( renderer Renderer )
    Optimize( geometry *Geometry, transform Transform )
}

//An Entity is something that can be scaled, positioned and rotated (orientation)
type Entity interface {
    SetScale( scale vectorMath.Vector3 )
    SetTranslation( translation vectorMath.Vector3 ) 
    SetOrientation( orientation vectorMath.Quaternion )
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

func (node *Node) OptimizeNode( geometry *Geometry ) {
    geometry.Verticies = make([]float32, 0)
    geometry.Indicies = make([]uint32, 0)
    node.Optimize(geometry, node.Transform)
    geometry.VboDirty = true
}

func (node *Node) Optimize( geometry *Geometry, transform Transform ) {
    newTransform := CreateTransform()
    newTransform.Set(transform)
    newTransform.ApplyTransform(node.Transform)
    for _,child := range node.children {
        child.Optimize(geometry, newTransform)
    }
}
