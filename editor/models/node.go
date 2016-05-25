package editorModels

import (
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type NodeModel struct {
	Id          string           `json:"id"`
	Classes     []string         `json:"classes"`
	Scale       vmath.Vector3    `json:"scale"`
	Translation vmath.Vector3    `json:"translation"`
	Orientation vmath.Quaternion `json:"orientation"`
	Geometry    *string          `json:"geometry"`
	Reference   *string          `json:"reference"`
	Children    []*NodeModel     `json:"children"`
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
		geom := *n.Geometry
		copiedNode.Geometry = &geom
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
		Scale:       vmath.Vector3{1, 1, 1},
		Translation: vmath.Vector3{},
		Orientation: vmath.IdentityQuaternion(),
		Children:    make([]*NodeModel, 0),
	}
}
