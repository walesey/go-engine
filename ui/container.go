package ui

import (
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Container struct {
	HorizontalAlign bool
	Hitbox          Hitbox
	node            *renderer.Node
	size, offset    vmath.Vector2
	children        []Element
}

func (c *Container) Render(offset vmath.Vector2) vmath.Vector2 {
	var width float64 = 0
	var height float64 = 0
	for _, child := range c.children {
		if c.HorizontalAlign {
			size := child.Render(vmath.Vector2{width, 0})
			width = width + size.X
			if size.Y > height {
				height = size.Y
			}
		} else {
			size := child.Render(vmath.Vector2{0, height})
			height = height + size.Y
			if size.X > width {
				width = size.X
			}
		}
	}
	size := vmath.Vector2{width, height}
	// c.background.SetScale(size.ToVector3()) //TODO: add background node
	c.node.SetTranslation(offset.ToVector3())
	c.Hitbox.SetSize(size)
	c.offset = offset
	return size
}

func (c *Container) Spatial() renderer.Spatial {
	return c.node
}

func (c *Container) AddChildren(children ...Element) {
	c.children = append(c.children, children...)
	for _, child := range children {
		c.node.Add(child.Spatial())
	}
}

func (c *Container) RemoveChildren(children ...Element) {
	for _, child := range children {
		c.node.Remove(child.Spatial())
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

func NewContainer() *Container {
	node := renderer.CreateNode()
	return &Container{
		node:     node,
		children: make([]Element, 0),
		Hitbox:   NewHitbox(),
	}
}
