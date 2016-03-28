package ui

import (
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

func NewWindow() renderer.Spatial {
	window := renderer.CreateBox(100, 100)
	// window.Material = mat
	// window.Material.LightingMode = renderer.MODE_UNLIT
	node := renderer.CreateNode()
	node.SetTranslation(vmath.Vector3{600, 600})
	node.Add(window)
	return node
}
