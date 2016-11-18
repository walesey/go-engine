package ui

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

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

func (w *Window) Centre() mgl32.Vec3 {
	return w.node.Centre()
}

func (w *Window) SetParent(parent *renderer.Node) {
	w.node.SetParent(parent)
}

func (w *Window) Optimize(geometry *renderer.Geometry, transform mgl32.Mat4) {
	w.node.Optimize(geometry, transform)
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
		container, ok := w.element.(*Container)
		if ok {
			return container.ElementById(id)
		}
	}
	return nil
}

func (w *Window) TextElementById(id string) *TextElement {
	container, ok := w.ElementById(id).(*Container)
	if ok && container.GetNbChildren() > 0 {
		if textElement, ok := container.GetChild(0).(*TextElement); ok {
			return textElement
		}
		if textField, ok := container.GetChild(0).(*TextField); ok {
			return textField.text
		}
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
		if !release {
			deactivateAllTextFields(w.element)
		}
		w.element.mouseClick(button, release, w.mousePos)
	}
}

func (w *Window) keyClick(key string, release bool) {
	if w.element != nil {
		w.element.keyClick(key, release)
	}
}

func deactivateAllTextFields(elem Element) {
	if container, ok := elem.(*Container); ok {
		for i := 0; i < container.GetNbChildren(); i++ {
			child := container.GetChild(i)
			deactivateAllTextFields(child)
			text, ok := child.(*TextField)
			if ok {
				text.Deactivate()
			}
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
	return &Window{
		node:          node,
		backgroundBox: box,
		background:    background,
		elementNode:   elementNode,
		size:          mgl32.Vec2{500, 1},
		Tabs:          []Activatable{},
	}
}
