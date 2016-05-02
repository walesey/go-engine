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
