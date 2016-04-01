package ui

import (
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Window struct {
	node    *renderer.Node
	element Element
	size    vmath.Vector2
}

func (w *Window) Draw(renderer renderer.Renderer) {
	w.node.Draw(renderer)
}

func (w *Window) Centre() vmath.Vector3 {
	return w.node.Centre()
}

func (w *Window) Optimize(geometry *renderer.Geometry, transform renderer.Transform) {
	w.node.Optimize(geometry, transform)
}

func (w *Window) SetScale(scale vmath.Vector3) {
	w.node.SetScale(scale)
}

func (w *Window) SetTranslation(translation vmath.Vector3) {
	w.node.SetTranslation(translation)
}

func (w *Window) SetOrientation(orientation vmath.Quaternion) {
	w.node.SetOrientation(orientation)
}

func (w *Window) SetElement(element Element) {
	if w.element != nil {
		w.node.Remove(w.element.Spatial())
	}
	w.element = element
	w.node.Add(element.Spatial())
	w.Render()
}

func (w *Window) Render() {
	w.element.Render(vmath.Vector2{0, 0})
}

func NewWindow() *Window {
	node := renderer.CreateNode()
	window := &Window{
		node: node,
	}
	return window
}
