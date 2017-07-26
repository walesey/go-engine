package ui

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

var activeWindow *Window

type Activatable interface {
	Active() bool
	Activate()
	Deactivate()
}

type Window struct {
	node, elementNode, background *renderer.Node
	backgroundBox                 *renderer.Geometry
	element                       Element
	size, position                mgl32.Vec2
	mousePos                      mgl32.Vec2
	Tabs                          []Activatable
}

func (w *Window) Draw(renderer renderer.Renderer, transform mgl32.Mat4) {
	w.node.Draw(renderer, transform)
}

func (w *Window) Destroy(renderer renderer.Renderer) {
	w.node.Destroy(renderer)
}

func (w *Window) Center() mgl32.Vec3 {
	return w.node.Center()
}

func (w *Window) SetParent(parent *renderer.Node) {
	w.node.SetParent(parent)
}

func (w *Window) Optimize(geometry *renderer.Geometry, transform mgl32.Mat4) {
	w.node.Optimize(geometry, transform)
}

func (w *Window) BoundingRadius() float32 {
	return w.node.BoundingRadius()
}

func (w *Window) OrthoOrder() int {
	return w.node.OrthoOrder()
}

func (w *Window) SetScale(scale mgl32.Vec3) {
	w.background.SetScale(scale)
	w.size = scale.Vec2()
	w.Render()
}

func (w *Window) SetTranslation(translation mgl32.Vec3) {
	w.node.SetTranslation(translation)
	w.position = translation.Vec2()
}

func (w *Window) SetOrientation(orientation mgl32.Quat) {
	w.node.SetOrientation(orientation)
}

func (w *Window) SetBackgroundColor(r, g, b, a uint8) {
	w.backgroundBox.SetColor(color.NRGBA{r, g, b, a})
}

func (w *Window) SetElement(element Element) {
	if w.element != nil {
		w.elementNode.Remove(w.element.Spatial(), true)
	}
	w.element = element
	w.elementNode.Add(element.Spatial())
	w.Render()
}

func (w *Window) Render() {
	if w.element != nil {
		size := w.element.Render(w.size, mgl32.Vec2{0, 0})
		width, height := w.size.X(), w.size.Y()
		if size.X() > width {
			width = size.X()
		}
		if size.Y() > height {
			height = size.Y()
		}
		w.background.SetScale(mgl32.Vec2{width, height}.Vec3(0))
	}
}

func (w *Window) ElementById(id string) Element {
	if w.element != nil {
		if w.element.GetId() == id {
			return w.element
		}
		return w.element.GetChildren().GetChildById(id)
	}
	return nil
}

func (w *Window) TextElementById(id string) *TextElement {
	if w.element != nil {
		return w.element.GetChildren().TextElementById(id)
	}
	return nil
}

func (w *Window) mouseMove(position mgl32.Vec2) {
	w.mousePos = position.Sub(w.position)
	if w.element != nil {
		w.element.mouseMove(w.mousePos)
	}
}

func (w *Window) mouseClick(button int, release bool) {
	if w.element != nil {
		// Deselect all fields
		if !release {
			deactivateAllFields(w.element)
		}
		// Process the click
		if w.element.mouseClick(button, release, w.mousePos) {
			// set this to the active window
			if !release {
				if activeWindow != nil {
					activeWindow.node.OrthoOrderValue = 10
				}
				w.node.OrthoOrderValue = 100
				activeWindow = w
			}
		}
	}
}

func (w *Window) keyClick(key string, release bool) {
	if w.element != nil {
		w.element.keyClick(key, release)
	}
}

func deactivateAllFields(elem Element) {
	for _, child := range elem.GetChildren() {
		deactivateAllFields(child)
		switch t := child.(type) {
		case *TextField:
			t.Deactivate()
		case *Dropdown:
			t.deactivate_setFlag()
		}
	}
}

func NewWindow() *Window {
	node := renderer.NewNode()
	elementNode := renderer.NewNode()
	background := renderer.NewNode()
	background.Material = renderer.NewMaterial()
	box := renderer.CreateBoxWithOffset(1, 1, 0, 0)
	box.SetColor(color.NRGBA{255, 255, 255, 255})
	background.Add(box)
	node.Add(background)
	node.Add(elementNode)
	node.OrthoOrderValue = 10
	return &Window{
		node:          node,
		backgroundBox: box,
		background:    background,
		elementNode:   elementNode,
		size:          mgl32.Vec2{500, 1},
		Tabs:          []Activatable{},
	}
}
