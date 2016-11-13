package renderer

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/util"
)

//A Spatial is something that can be Drawn by a Renderer
type Spatial interface {
	Draw(renderer Renderer, transform mgl32.Mat4)
	Optimize(geometry *Geometry, transform mgl32.Mat4)
	Destroy(renderer Renderer)
	Centre() mgl32.Vec3
}

//An Entity is something that can be scaled, positioned and rotated (orientation)
type Entity interface {
	SetScale(scale mgl32.Vec3)
	SetTranslation(translation mgl32.Vec3)
	SetOrientation(orientation mgl32.Quat)
}

//Node
type Node struct {
	children    []Spatial
	deleted     []Spatial
	Transform   mgl32.Mat4
	Scale       mgl32.Vec3
	Translation mgl32.Vec3
	Orientation mgl32.Quat
}

func CreateNode() *Node {
	node := &Node{
		children:    make([]Spatial, 0, 0),
		deleted:     make([]Spatial, 0, 0),
		Translation: mgl32.Vec3{0, 0, 0},
		Orientation: mgl32.QuatIdent(),
	}
	node.SetScale(mgl32.Vec3{1, 1, 1})
	return node
}

func (node *Node) Draw(renderer Renderer, transform mgl32.Mat4) {
	tx := transform.Mul4(node.Transform)
	for _, child := range node.children {
		child.Draw(renderer, tx)
	}
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

func (node *Node) Centre() mgl32.Vec3 {
	return node.Translation
}

func (node *Node) Add(spatial Spatial) {
	node.children = append(node.children, spatial)
}

func (node *Node) Remove(spatial Spatial, destroy bool) {
	for i, child := range node.children {
		if child == spatial {
			node.children[i] = node.children[len(node.children)-1]
			node.children[len(node.children)-1] = nil
			node.children = node.children[:len(node.children)-1]
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

func (node *Node) SetScale(scale mgl32.Vec3) {
	node.Scale = scale
	node.Transform = util.Mat4From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) SetTranslation(translation mgl32.Vec3) {
	node.Translation = translation
	node.Transform = util.Mat4From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) SetOrientation(orientation mgl32.Quat) {
	node.Orientation = orientation
	node.Transform = util.Mat4From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) SetRotation(angle float32, axis mgl32.Vec3) {
	node.Orientation = mgl32.QuatRotate(angle, axis)
	node.Transform = util.Mat4From(node.Scale, node.Translation, node.Orientation)
}

func (node *Node) OptimizeNode() *Geometry {
	geometry := CreateGeometry(make([]uint32, 0, 0), make([]float32, 0, 0))
	node.Optimize(geometry, node.Transform)
	geometry.VboDirty = true
	return geometry
}

func (node *Node) Optimize(geometry *Geometry, transform mgl32.Mat4) {
	newTransform := transform.Mul4(node.Transform)
	for _, child := range node.children {
		child.Optimize(geometry, newTransform)
	}
}

func (node *Node) RelativePosition(n *Node) (mgl32.Vec3, error) {
	if node == n {
		return mgl32.Vec3{}, nil
	}
	for _, child := range node.children {
		if childNode, ok := child.(*Node); ok {
			if rPost, err := childNode.RelativePosition(n); err == nil {
				return mgl32.TransformCoordinate(rPost, childNode.Transform), nil
			}
		}
	}
	return mgl32.Vec3{}, fmt.Errorf("Node not found")
}
