package editorModels

import (
	vmath "github.com/walesey/go-engine/vectormath"
)

type NodeModel struct {
	Id          string           `json:"id"`
	Classes     []string         `json:"classes"`
	Scale       vmath.Vector3    `json:"scale"`
	Translation vmath.Vector3    `json:"translation"`
	Orientation vmath.Quaternion `json:"orientation"`
	Geometry    *string          `json:"geometry"`
	Children    []*NodeModel     `json:"children"`
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
