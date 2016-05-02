package editorModels

type MapModel struct {
	Name string     `json:"name"`
	Root *NodeModel `json:"root"`
}
