package renderer

import (
	"github.com/walesey/go-engine/vectormath"
)

//A Spatial is something that can be Drawn by a Renderer
type Spatial interface {
	Draw(renderer Renderer)
	Centre() vectormath.Vector3
	Optimize(geometry *Geometry, transform Transform)
}

//An Entity is something that can be scaled, positioned and rotated (orientation)
type Entity interface {
	SetScale(scale vectormath.Vector3)
	SetTranslation(translation vectormath.Vector3)
	SetOrientation(orientation vectormath.Quaternion)
}

//Node
type Node struct {
	children    []Spatial
	Transform   Transform
	Scale       vectormath.Vector3
	Translation vectormath.Vector3
	Orientation vectormath.Quaternion
}

func CreateNode() Node {
	//create slice to store children
	children := make([]Spatial, 0, 0)
	return Node{
		children:    children,
		Transform:   CreateTransform(),
		Scale:       vectormath.Vector3{1, 1, 1},
		Translation: vectormath.Vector3{0, 0, 0},
		Orientation: vectormath.IdentityQuaternion(),
	}
}

func (node *Node) Draw(renderer Renderer) {
	renderer.PushTransform()
	renderer.ApplyTransform(node.Transform)
	for _, child := range node.children {
		child.Draw(renderer)
	}
	renderer.PopTransform()
}

func (node *Node) Centre() vectormath.Vector3 {
	return node.Translation
}

func (node *Node) Add(spatial Spatial) {
	//append to the slice
	node.children = append(node.children, spatial)
}

func (node *Node) Remove(spatial Spatial) {
	//find the address in the slice, remove it and return
	for index, child := range node.children {
		if child == spatial {
			if index+1 == len(node.children) {
				node.children = node.children[:index]
			} else {
				node.children = append(node.children[:index], node.children[index+1:]...)
			}
			break
		}
	}
}

func (node *Node) SetScale(scale vectormath.Vector3) {
	node.Scale = scale
	node.Transform.From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) SetTranslation(translation vectormath.Vector3) {
	node.Translation = translation
	node.Transform.From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) SetOrientation(orientation vectormath.Quaternion) {
	node.Orientation = orientation
	node.Transform.From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) SetRotation(angle float64, axis vectormath.Vector3) {
	node.Orientation = vectormath.AngleAxis(angle, axis)
	node.Transform.From(node.Scale, node.Translation, node.Orientation)
}

//used for eg. sprites facing the direction of the camera - all vectors need to be normalized
func (node *Node) SetFacing(rotation float64, newNormal, normal, tangent vectormath.Vector3) {
	FacingTransform(node.Transform, rotation, newNormal, normal, tangent)
}

func (node *Node) OptimizeNode() Geometry {
	geometry := CreateGeometry(make([]uint32, 0, 0), make([]float32, 0, 0))
	node.Optimize(&geometry, node.Transform)
	geometry.VboDirty = true
	return geometry
}

func (node *Node) Optimize(geometry *Geometry, transform Transform) {
	newTransform := CreateTransform()
	newTransform.Set(transform)
	newTransform.ApplyTransform(node.Transform)
	for _, child := range node.children {
		child.Optimize(geometry, newTransform)
	}
}
