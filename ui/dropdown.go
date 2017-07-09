package ui

import (
	"image"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/libs/freetype/truetype"
	"github.com/walesey/go-engine/renderer"
)

type Dropdown struct {
	Parent               Element
	node                 *renderer.Node
	container            *Container
	dropdown             *Container
	text                 *TextElement
	active               bool
	dropdownFlag         bool
	dropdownFlagState    bool
	dropdownVisible      bool
	options              []string
	dropdownHeightOffset mgl32.Vec2
	dropdownMarginOffset mgl32.Vec2
	onFocusHandlers      []func()
	onBlurHandlers       []func()
	onChangeHandlers     []func(string)
}

func (dd *Dropdown) Render(size, offset mgl32.Vec2) mgl32.Vec2 {
	dropdownSize := mgl32.Vec2{size[0], 99999}
	dd.dropdown.Render(dropdownSize, offset.Add(dd.dropdownMarginOffset).Add(dd.dropdownHeightOffset))
	renderSize := dd.container.Render(size, offset)
	return renderSize
}

func (dd *Dropdown) ReRender() {
	dd.container.ReRender()
}

func (dd *Dropdown) Spatial() renderer.Spatial {
	return dd.node
}

func (dd *Dropdown) GlobalOrthoOrder() int {
	return dd.node.OrthoOrder()
}

func (dd *Dropdown) GetId() string {
	return dd.container.GetId()
}

func (dd *Dropdown) SetId(id string) {
	dd.container.SetId(id)
}

func (dd *Dropdown) GetChildren() Children {
	return []Element{}
}

func (dd *Dropdown) mouseMove(position mgl32.Vec2) bool {
	containerMoved := dd.container.mouseMove(position)
	dropdownMoved := dd.isDropdownVisible() && dd.dropdown.mouseMove(position)
	return containerMoved || dropdownMoved
}

func (dd *Dropdown) mouseClick(button int, release bool, position mgl32.Vec2) bool {
	containerClicked := dd.container.mouseClick(button, release, position)
	dropdownClicked := dd.isDropdownVisible() && dd.dropdown.mouseClick(button, release, position)
	return containerClicked || dropdownClicked
}

func (dd *Dropdown) onClickOption(option string) func(button int, release bool, position mgl32.Vec2) {
	return func(button int, release bool, position mgl32.Vec2) {
		if release {
			dd.SetText(option)
			dd.ReRender()
			dd.dropdownFlag = false
		}
	}
}

func (dd *Dropdown) keyClick(key string, release bool) {
	if dd.active && !release {
		if key == "upArrow" {
			//TODO
		} else if key == "downArrow" {
			//TODO
		}
		for _, handler := range dd.text.onKeyPressHandlers {
			handler(key, release)
		}
		dd.ReRender()
	}
}

func (dd *Dropdown) SetBackgroundImage(img image.Image) {
	dd.container.SetBackgroundImage(img)
}

func (dd *Dropdown) SetBackgroundColor(r, g, b, a uint8) {
	dd.container.SetBackgroundColor(r, g, b, a)
}

func (dd *Dropdown) GetText() string {
	return dd.text.GetText()
}

func (dd *Dropdown) GetHiddenText() string {
	return dd.text.GetHiddenText()
}

func (dd *Dropdown) GetHitbox() Hitbox {
	return dd.container.Hitbox
}

func (dd *Dropdown) SetText(text string) *Dropdown {
	for _, option := range dd.options {
		if option == text {
			dd.text.SetText(text)
			for _, handler := range dd.onChangeHandlers {
				handler(text)
			}
		}
	}
	return dd
}

func (dd *Dropdown) SetPlaceholder(placeholder string) *Dropdown {
	dd.text.SetPlaceholder(placeholder)
	return dd
}

func (dd *Dropdown) SetFont(textFont *truetype.Font) *Dropdown {
	dd.text.SetFont(textFont)
	return dd
}

func (dd *Dropdown) SetTextSize(textSize float32) *Dropdown {
	dd.text.SetTextSize(textSize)
	return dd
}

func (dd *Dropdown) SetTextColor(textColor color.Color) *Dropdown {
	dd.text.SetTextColor(textColor)
	return dd
}

func (dd *Dropdown) SetWidth(width float32) {
	dd.container.SetWidth(width)
	dd.dropdown.SetWidth(width)
}

func (dd *Dropdown) UsePercentWidth(usePercent bool) {
	dd.container.UsePercentWidth(usePercent)
	dd.dropdown.UsePercentWidth(usePercent)
}

func (dd *Dropdown) SetHeight(height float32) {
	dd.container.SetHeight(height)
	dd.dropdown.SetHeight(height * float32(len(dd.options)))
	dd.dropdownHeightOffset = mgl32.Vec2{0, height}
}

func (dd *Dropdown) UsePercentHeight(usePercent bool) {
	dd.container.UsePercentHeight(usePercent)
	dd.dropdown.UsePercentHeight(usePercent)
}

func (dd *Dropdown) SetMargin(margin Margin) {
	dd.container.SetMargin(margin)
	dd.dropdownMarginOffset = mgl32.Vec2{margin.Left, margin.Top}
}

func (dd *Dropdown) SetMarginPercent(marginPercent MarginPercentages) {
	dd.container.SetMarginPercent(marginPercent)
}

func (dd *Dropdown) SetPadding(padding Margin) {
	dd.container.SetPadding(padding)
}

func (dd *Dropdown) SetPaddingPercent(paddingPercent MarginPercentages) {
	dd.container.SetPaddingPercent(paddingPercent)
}

func (dd *Dropdown) SetHidden(hidden bool) *Dropdown {
	dd.text.SetHidden(hidden)
	return dd
}

func (dd *Dropdown) Active() bool {
	return dd.active
}

func (dd *Dropdown) Activate() {
	if !dd.active {
		dd.active = true
		for _, handler := range dd.onFocusHandlers {
			handler()
		}
		dd.ReRender()
	}
	dd.ToggleDropdown()
	dd.dropdownFlag = false
}

func (dd *Dropdown) Deactivate() {
	if dd.active {
		dd.active = false
		for _, handler := range dd.onBlurHandlers {
			handler()
		}
		dd.ReRender()
	}
	dd.CloseDropdown()
	dd.dropdownFlag = false
}

func (dd *Dropdown) deactivate_setFlag() {
	dd.dropdownFlagState = dd.dropdownVisible
	dd.Deactivate()
	dd.dropdownFlag = true
}

func (dd *Dropdown) isDropdownVisible() bool {
	if dd.dropdownFlag {
		return dd.dropdownFlagState
	}
	return dd.dropdownVisible
}

func (dd *Dropdown) ToggleDropdown() {
	if dd.isDropdownVisible() {
		dd.CloseDropdown()
	} else {
		dd.OpenDropdown()
	}
}

func (dd *Dropdown) OpenDropdown() {
	dd.node.Add(dd.dropdown.Spatial())
	dd.node.OrthoOrderValue = 9999
	dd.reRenderParent()
	dd.dropdownVisible = true
}

func (dd *Dropdown) CloseDropdown() {
	dd.node.Remove(dd.dropdown.Spatial(), false)
	dd.node.OrthoOrderValue = 0
	dd.reRenderParent()
	dd.dropdownVisible = false
}

func (dd *Dropdown) reRenderParent() {
	if dd.Parent != nil {
		dd.Parent.ReRender()
	}
}

func (dd *Dropdown) AddOnFocus(handler func()) {
	dd.onFocusHandlers = append(dd.onFocusHandlers, handler)
}

func (dd *Dropdown) AddOnBlur(handler func()) {
	dd.onBlurHandlers = append(dd.onBlurHandlers, handler)
}

func (dd *Dropdown) AddOnKeyPress(handler func(key string, release bool)) {
	dd.AddOnKeyPress(handler)
}

func (dd *Dropdown) AddOnChange(handler func(string)) {
	dd.onChangeHandlers = append(dd.onChangeHandlers, handler)
}

func NewDropdown(options []string, textColor color.Color, textSize float32, textFont *truetype.Font, parent Element) *Dropdown {
	dd := &Dropdown{
		Parent:    parent,
		node:      renderer.NewNode(),
		container: NewContainer(),
		dropdown:  NewContainer(),
		text:      NewTextElement("", textColor, textSize, textFont),
		options:   options,
	}
	dd.node.Add(dd.container.Spatial())
	dd.container.AddChildren(dd.text)
	dd.GetHitbox().AddOnClick(func(button int, release bool, position mgl32.Vec2) {
		if !release {
			dd.Activate()
		}
	})
	dd.SetBackgroundColor(0, 0, 0, 0)
	dd.SetHeight(textSize * 1.5)
	dd.SetPadding(NewMargin(2))
	// setup dropdown container
	for _, option := range options {
		txtC := NewContainer()
		txtC.SetPadding(Margin{textSize * 0.5, 2, textSize * 0.5, 2})
		txtC.SetBackgroundColor(255, 255, 255, 255)
		txtC.Hitbox.AddOnHover(func() { txtC.SetBackgroundColor(240, 240, 240, 255) })
		txtC.Hitbox.AddOnUnHover(func() { txtC.SetBackgroundColor(255, 255, 255, 255) })
		txtC.Hitbox.AddOnClick(dd.onClickOption(option))
		txt := NewTextElement(option, color.Black, textSize, textFont)
		txtC.AddChildren(txt)
		dd.dropdown.AddChildren(txtC)
	}
	return dd
}
