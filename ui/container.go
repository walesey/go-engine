package ui

import (
	"image/color"

	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Margin struct {
	Top, Right, Bottom, Left float64
}

type Container struct {
	Hitbox          Hitbox
	width, height   float64
	margin, padding Margin
	node            *renderer.Node
	elementsNode    *renderer.Node
	background      *renderer.Node
	backgroundBox   *renderer.Geometry
	size, offset    vmath.Vector2
	children        []Element
}

func (c *Container) Render(size, offset vmath.Vector2) vmath.Vector2 {
	sizeMinusMargins := size.Subtract(vmath.Vector2{
		c.margin.Left + c.margin.Right + c.padding.Left + c.padding.Right,
		c.margin.Top + c.margin.Bottom + c.padding.Top + c.padding.Bottom,
	})
	containerSize := &sizeMinusMargins
	if c.width > 0 {
		containerSize.X = c.width
	}
	if c.height > 0 {
		containerSize.Y = c.height
	}
	var width, height, highest float64 = 0, 0, 0
	for _, child := range c.children {
		childSize := child.Render(*containerSize, vmath.Vector2{width, height})
		width += childSize.X
		if width > containerSize.X {
			height += highest
			highest = 0
			childSize = child.Render(*containerSize, vmath.Vector2{0, height})
			width = childSize.X
		}
		if childSize.Y > highest {
			highest = childSize.Y
		}
	}
	height += highest
	if vmath.ApproxEqual(c.height, 0, 0.001) {
		containerSize.Y = height
	}
	//offsets and sizes
	c.offset = offset.Add(vmath.Vector2{c.margin.Left, c.margin.Top})
	elementsOffset := vmath.Vector2{c.padding.Left, c.padding.Top}
	backgroundSize := containerSize.Add(vmath.Vector2{c.padding.Left + c.padding.Right, c.padding.Top + c.padding.Bottom})
	totalSize := backgroundSize.Add(vmath.Vector2{c.margin.Left + c.margin.Right, c.margin.Top + c.margin.Bottom})

	c.background.SetScale(backgroundSize.ToVector3())
	c.node.SetTranslation(c.offset.ToVector3())
	c.elementsNode.SetTranslation(elementsOffset.ToVector3())
	c.Hitbox.SetSize(backgroundSize)
	return totalSize
}

func (c *Container) Spatial() renderer.Spatial {
	return c.node
}

func (c *Container) SetSize(width, height float64) {
	c.width, c.height = width, height
}

func (c *Container) SetMargin(margin Margin) {
	c.margin = margin
}

func (c *Container) SetPadding(padding Margin) {
	c.padding = padding
}

func (c *Container) SetBackgroundColor(r, g, b, a uint8) {
	c.backgroundBox.SetColor(color.NRGBA{r, g, b, a})
}

func (c *Container) AddChildren(children ...Element) {
	c.children = append(c.children, children...)
	for _, child := range children {
		c.elementsNode.Add(child.Spatial())
	}
}

func (c *Container) RemoveChildren(children ...Element) {
	for _, child := range children {
		c.elementsNode.Remove(child.Spatial(), true)
		for index, containerChild := range c.children {
			if containerChild == child {
				c.children = append(c.children[:index], c.children[index+1:]...)
			}
		}
	}
}

func (c *Container) mouseMove(position vmath.Vector2) {
	offsetPos := position.Subtract(c.offset)
	c.Hitbox.MouseMove(offsetPos)
	for _, child := range c.children {
		child.mouseMove(offsetPos)
	}
}

func (c *Container) mouseClick(button int, release bool, position vmath.Vector2) {
	offsetPos := position.Subtract(c.offset)
	c.Hitbox.MouseClick(button, release, offsetPos)
	for _, child := range c.children {
		child.mouseClick(button, release, offsetPos)
	}
}

func (c *Container) keyClick(key string, release bool) {
	for _, child := range c.children {
		child.keyClick(key, release)
	}
}

func NewMargin(margin float64) Margin {
	return Margin{margin, margin, margin, margin}
}

func NewContainer() *Container {
	node := renderer.CreateNode()
	elementsNode := renderer.CreateNode()
	background := renderer.CreateNode()
	box := renderer.CreateBoxWithOffset(1, 1, 0, 0)
	box.SetColor(color.NRGBA{0, 0, 0, 0})
	mat := renderer.CreateMaterial()
	mat.LightingMode = renderer.MODE_UNLIT
	box.Material = mat
	background.Add(box)
	node.Add(background)
	node.Add(elementsNode)
	return &Container{
		node:          node,
		elementsNode:  elementsNode,
		background:    background,
		backgroundBox: box,
		children:      make([]Element, 0),
		Hitbox:        NewHitbox(),
		padding:       NewMargin(0),
		margin:        NewMargin(0),
	}
}
