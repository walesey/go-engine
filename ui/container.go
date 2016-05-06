package ui

import (
	"image/color"

	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Margin struct {
	Top, Right, Bottom, Left float64
}

type MarginPercentages struct {
	Top, Right, Bottom, Left bool
}

type Container struct {
	id               string
	Hitbox           Hitbox
	width, height    float64
	percentWidth     bool
	percentHeight    bool
	margin, padding  Margin
	marginPercent    MarginPercentages
	paddingPercent   MarginPercentages
	node             *renderer.Node
	elementsNode     *renderer.Node
	background       *renderer.Node
	backgroundBox    *renderer.Geometry
	size, offset     vmath.Vector2
	backgroundOffset vmath.Vector2
	elementsOffset   vmath.Vector2
	children         []Element
}

func (c *Container) Render(size, offset vmath.Vector2) vmath.Vector2 {
	padding := convertMargin(c.padding, c.paddingPercent, size.X)
	margin := convertMargin(c.margin, c.marginPercent, size.X)
	sizeMinusMargins := size.Subtract(vmath.Vector2{
		margin.Left + margin.Right + padding.Left + padding.Right,
		margin.Top + margin.Bottom + padding.Top + padding.Bottom,
	})
	containerSize := &sizeMinusMargins
	if c.width > 0 {
		containerSize.X = c.getWidth(size.X) - padding.Left - padding.Right
	}
	if c.height > 0 {
		containerSize.Y = c.getHeight(size.X) - padding.Top - padding.Bottom
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
	c.offset = offset
	c.backgroundOffset = vmath.Vector2{margin.Left, margin.Top}
	c.elementsOffset = vmath.Vector2{margin.Left + padding.Left, margin.Top + padding.Top}
	backgroundSize := containerSize.Add(vmath.Vector2{padding.Left + padding.Right, padding.Top + padding.Bottom})
	totalSize := backgroundSize.Add(vmath.Vector2{margin.Left + margin.Right, margin.Top + margin.Bottom})

	c.background.SetScale(backgroundSize.ToVector3())
	c.background.SetTranslation(c.backgroundOffset.ToVector3())
	c.elementsNode.SetTranslation(c.elementsOffset.ToVector3())
	c.node.SetTranslation(c.offset.ToVector3())
	c.Hitbox.SetSize(backgroundSize)
	return totalSize
}

func convertMargin(margin Margin, percentages MarginPercentages, parentWidth float64) Margin {
	result := &Margin{margin.Top, margin.Right, margin.Bottom, margin.Left}
	if percentages.Bottom {
		result.Bottom = parentWidth * margin.Bottom / 100.0
	}
	if percentages.Left {
		result.Left = parentWidth * margin.Left / 100.0
	}
	if percentages.Right {
		result.Right = parentWidth * margin.Right / 100.0
	}
	if percentages.Top {
		result.Top = parentWidth * margin.Top / 100.0
	}
	return *result
}

func (c *Container) Spatial() renderer.Spatial {
	return c.node
}

func (c *Container) SetWidth(width float64) {
	c.width = width
}

func (c *Container) UsePercentWidth(usePercent bool) {
	c.percentWidth = usePercent
}

func (c *Container) getWidth(parentWidth float64) float64 {
	if c.percentWidth {
		return parentWidth * c.width / 100.0
	}
	return c.width
}

func (c *Container) SetHeight(height float64) {
	c.height = height
}

func (c *Container) UsePercentHeight(usePercent bool) {
	c.percentHeight = usePercent
}

func (c *Container) getHeight(parentWidth float64) float64 {
	if c.percentHeight {
		return parentWidth * c.height / 100.0
	}
	return c.height
}

func (c *Container) SetMargin(margin Margin) {
	c.margin = margin
}

func (c *Container) SetMarginPercent(marginPercent MarginPercentages) {
	c.marginPercent = marginPercent
}

func (c *Container) SetPadding(padding Margin) {
	c.padding = padding
}

func (c *Container) SetPaddingPercent(paddingPercent MarginPercentages) {
	c.paddingPercent = paddingPercent
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

func (c *Container) RemoveAllChildren() {
	for _, child := range c.children {
		c.elementsNode.Remove(child.Spatial(), true)
	}
	c.children = c.children[:0]
}

func (c *Container) GetChild(index int) Element {
	if index >= c.GetNbChildren() {
		return nil
	}
	return c.children[index]
}

func (c *Container) GetNbChildren() int {
	return len(c.children)
}

func (c *Container) GetId() string {
	return c.id
}

func (c *Container) ElementById(id string) Element {
	for _, child := range c.children {
		if child.GetId() == id {
			return child
		}
		container, ok := child.(*Container)
		if ok {
			if elem := container.ElementById(id); elem != nil {
				return elem
			}
		}
	}
	return nil
}

func (c *Container) mouseMove(position vmath.Vector2) {
	offsetPos := position.Subtract(c.offset)
	c.Hitbox.MouseMove(offsetPos.Subtract(c.backgroundOffset))
	for _, child := range c.children {
		child.mouseMove(offsetPos.Subtract(c.elementsOffset))
	}
}

func (c *Container) mouseClick(button int, release bool, position vmath.Vector2) {
	offsetPos := position.Subtract(c.offset)
	c.Hitbox.MouseClick(button, release, offsetPos.Subtract(c.backgroundOffset))
	for _, child := range c.children {
		child.mouseClick(button, release, offsetPos.Subtract(c.elementsOffset))
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

func NewMarginPercentages(percentage bool) MarginPercentages {
	return MarginPercentages{percentage, percentage, percentage, percentage}
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
