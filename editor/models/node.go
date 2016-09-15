package editorModels

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

type NodeModel struct {
	Id          string       `json:"id"`
	Classes     []string     `json:"classes"`
	Scale       mgl32.Vec3   `json:"scale"`
	Translation mgl32.Vec3   `json:"translation"`
	Orientation mgl32.Quat   `json:"orientation"`
	Geometry    *string      `json:"geometry"`
	Reference   *string      `json:"reference"`
	Children    []*NodeModel `json:"children"`
	node        *renderer.Node
}

func (n *NodeModel) GetNode() *renderer.Node {
	return n.node
}

func (n *NodeModel) SetNode(node *renderer.Node) {
	n.node = node
}

func (n *NodeModel) CopyShallow(nameGenerator func(name string) string) *NodeModel {
	copiedNode := &NodeModel{
		Id:          nameGenerator(n.Id),
		Classes:     n.Classes,
		Scale:       n.Scale,
		Translation: n.Translation,
		Orientation: n.Orientation,
		Children:    make([]*NodeModel, len(n.Children)),
	}
	if n.Geometry != nil {
		copiedNode.Geometry = n.Geometry
	}
	if n.Reference != nil {
		ref := *n.Reference
		copiedNode.Reference = &ref
	}
	return copiedNode
}

func (n *NodeModel) Copy(nameGenerator func(name string) string) *NodeModel {
	copiedNode := n.CopyShallow(nameGenerator)
	for i, child := range n.Children {
		copiedNode.Children[i] = child.Copy(nameGenerator)
	}
	return copiedNode
}

func NewNodeModel(id string) *NodeModel {
	return &NodeModel{
		Id:          id,
		Scale:       mgl32.Vec3{1, 1, 1},
		Translation: mgl32.Vec3{},
		Orientation: mgl32.QuatIdent(),
		Children:    make([]*NodeModel, 0),
	}
}
