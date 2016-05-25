package renderer

import (
	"fmt"

	vmath "github.com/walesey/go-engine/vectormath"
)

//A Spatial is something that can be Drawn by a Renderer
type Spatial interface {
	Draw(renderer Renderer)
	Destroy(renderer Renderer)
	Centre() vmath.Vector3
	Optimize(geometry *Geometry, transform Transform)
}

//An Entity is something that can be scaled, positioned and rotated (orientation)
type Entity interface {
	SetScale(scale vmath.Vector3)
	SetTranslation(translation vmath.Vector3)
	SetOrientation(orientation vmath.Quaternion)
}

//Node
type Node struct {
	children    []Spatial
	deleted     []Spatial
	Transform   Transform
	Scale       vmath.Vector3
	Translation vmath.Vector3
	Orientation vmath.Quaternion
}

func CreateNode() *Node {
	return &Node{
		children:    make([]Spatial, 0, 0),
		deleted:     make([]Spatial, 0, 0),
		Transform:   CreateTransform(),
		Scale:       vmath.Vector3{1, 1, 1},
		Translation: vmath.Vector3{0, 0, 0},
		Orientation: vmath.IdentityQuaternion(),
	}
}

func (node *Node) Draw(renderer Renderer) {
	renderer.PushTransform()
	renderer.ApplyTransform(node.Transform)
	for _, child := range node.children {
		child.Draw(renderer)
	}
	renderer.PopTransform()
	node.cleanupDeleted(renderer)
}

func (node *Node) Destroy(renderer Renderer) {
	for _, child := range node.children {
		child.Destroy(renderer)
	}
	node.cleanupDeleted(renderer)
}

func (node *Node) cleanupDeleted(renderer Renderer) {
	for _, child := range node.deleted {
		child.Destroy(renderer)
	}
	node.deleted = node.deleted[:0]
}

func (node *Node) Centre() vmath.Vector3 {
	return node.Translation
}

func (node *Node) Add(spatial Spatial) {
	node.children = append(node.children, spatial)
}

func (node *Node) Remove(spatial Spatial, destroy bool) {
	for index, child := range node.children {
		if child == spatial {
			node.children = append(node.children[:index], node.children[index+1:]...)
			if destroy {
				node.deleted = append(node.deleted, child)
			}
			break
		}
	}
}

func (node *Node) RemoveAll(destroy bool) {
	if destroy {
		node.deleted = append(node.deleted, node.children...)
	}
	node.children = node.children[:0]
}

func (node *Node) SetScale(scale vmath.Vector3) {
	node.Scale = scale
	node.Transform.From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) SetTranslation(translation vmath.Vector3) {
	node.Translation = translation
	node.Transform.From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) SetOrientation(orientation vmath.Quaternion) {
	node.Orientation = orientation
	node.Transform.From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) SetRotation(angle float64, axis vmath.Vector3) {
	node.Orientation = vmath.AngleAxis(angle, axis)
	node.Transform.From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) OptimizeNode() *Geometry {
	geometry := CreateGeometry(make([]uint32, 0, 0), make([]float32, 0, 0))
	node.Optimize(geometry, node.Transform)
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

func (node *Node) RelativePosition(n *Node) (vmath.Vector3, error) {
	if node == n {
		return vmath.Vector3{}, nil
	}
	for _, child := range node.children {
		if childNode, ok := child.(*Node); ok {
			if rPost, err := childNode.RelativePosition(n); err == nil {
				return rPost.Add(childNode.Translation), nil
			}
		}
	}
	return vmath.Vector3{}, fmt.Errorf("Node not found")
}
