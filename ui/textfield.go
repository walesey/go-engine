package ui

import (
	"image"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/libs/freetype/truetype"
	"github.com/walesey/go-engine/renderer"
)

type TextField struct {
	container        *Container
	text             *TextElement
	cursor           *renderer.Node
	cursorPos        int
	active           bool
	onFocusHandlers  []func()
	onBlurHandlers   []func()
	onChangeHandlers []func(string)
}

func (tf *TextField) Render(size, offset mgl32.Vec2) mgl32.Vec2 {
	renderSize := tf.container.Render(size, offset)
	tf.RenderCursor()
	return renderSize
}

func (tf *TextField) ReRender() {
	tf.container.ReRender()
	tf.RenderCursor()
}

func (tf *TextField) RenderCursor() {
	c := tf.text.getContext()
	text := tf.GetHiddenText()
	cursorTranslation, _ := c.StringDimensions(string([]byte(text)[:tf.cursorPos]))
	xPos := int(cursorTranslation.X >> 6)
	tf.cursor.SetTranslation(mgl32.Vec2{float32(xPos), 0}.Vec3(0))
	if tf.active {
		tf.cursor.SetScale(mgl32.Vec2{tf.text.props.textSize, tf.text.props.textSize}.Vec3(0))
	} else {
		tf.cursor.SetScale(mgl32.Vec2{0, 0}.Vec3(0))
	}
}

func (tf *TextField) Spatial() renderer.Spatial {
	return tf.container.Spatial()
}

func (tf *TextField) GlobalOrthoOrder() int {
	return tf.container.GlobalOrthoOrder()
}

func (tf *TextField) GetId() string {
	return tf.container.GetId()
}

func (tf *TextField) SetId(id string) {
	tf.container.SetId(id)
}

func (tf *TextField) GetChildren() Children {
	return []Element{}
}

func (tf *TextField) mouseMove(position mgl32.Vec2) bool {
	return tf.container.mouseMove(position)
}

func (tf *TextField) mouseClick(button int, release bool, position mgl32.Vec2) bool {
	return tf.container.mouseClick(button, release, position)
}

func (tf *TextField) keyClick(key string, release bool) {
	if tf.active && !release {
		textBytes := []byte(tf.text.GetText())
		if key == "backspace" {
			if len(textBytes) > 0 && tf.cursorPos > 0 {
				cursorPos := tf.cursorPos
				newText := append(textBytes[:cursorPos-1], textBytes[cursorPos:]...)
				tf.SetText(string(newText))
				tf.cursorPos = cursorPos - 1
			}
		} else if key == "leftArrow" {
			if tf.cursorPos > 0 {
				tf.cursorPos--
			}
		} else if key == "rightArrow" {
			if tf.cursorPos < len(tf.text.GetText()) {
				tf.cursorPos++
			}
		} else {
			insertText := []rune(key)
			newText := []rune(tf.text.GetText())
			newText = append(newText[:tf.cursorPos], append(insertText, newText[tf.cursorPos:]...)...)
			cursorPos := tf.cursorPos
			tf.SetText(string(newText))
			tf.cursorPos = cursorPos + 1
		}
		for _, handler := range tf.text.onKeyPressHandlers {
			handler(key, release)
		}
		tf.ReRender()
	}
}

func (tf *TextField) SetBackgroundImage(img image.Image) {
	tf.container.SetBackgroundImage(img)
}

func (tf *TextField) SetBackgroundColor(r, g, b, a uint8) {
	tf.container.SetBackgroundColor(r, g, b, a)
}

func (tf *TextField) GetText() string {
	return tf.text.GetText()
}

func (tf *TextField) GetHiddenText() string {
	return tf.text.GetHiddenText()
}

func (tf *TextField) GetHitbox() Hitbox {
	return tf.container.Hitbox
}

func (tf *TextField) SetText(text string) *TextField {
	tf.text.SetText(text)
	tf.cursorPos = len(text)
	for _, handler := range tf.onChangeHandlers {
		handler(text)
	}
	return tf
}

func (tf *TextField) SetPlaceholder(placeholder string) *TextField {
	tf.text.SetPlaceholder(placeholder)
	return tf
}

func (tf *TextField) SetFont(textFont *truetype.Font) *TextField {
	tf.text.SetFont(textFont)
	return tf
}

func (tf *TextField) SetTextSize(textSize float32) *TextField {
	tf.text.SetTextSize(textSize)
	return tf
}

func (tf *TextField) SetTextColor(textColor color.Color) *TextField {
	tf.text.SetTextColor(textColor)
	return tf
}

func (tf *TextField) SetWidth(width float32) {
	tf.container.SetWidth(width)
}

func (tf *TextField) UsePercentWidth(usePercent bool) {
	tf.container.UsePercentWidth(usePercent)
}

func (tf *TextField) SetHeight(height float32) {
	tf.container.SetHeight(height)
}

func (tf *TextField) UsePercentHeight(usePercent bool) {
	tf.container.UsePercentHeight(usePercent)
}

func (tf *TextField) SetMargin(margin Margin) {
	tf.container.SetMargin(margin)
}

func (tf *TextField) SetMarginPercent(marginPercent MarginPercentages) {
	tf.container.SetMarginPercent(marginPercent)
}

func (tf *TextField) SetPadding(padding Margin) {
	tf.container.SetPadding(padding)
}

func (tf *TextField) SetPaddingPercent(paddingPercent MarginPercentages) {
	tf.container.SetPaddingPercent(paddingPercent)
}

func (tf *TextField) SetHidden(hidden bool) *TextField {
	tf.text.SetHidden(hidden)
	return tf
}

func (tf *TextField) Active() bool {
	return tf.active
}

func (tf *TextField) Activate() {
	if !tf.active {
		tf.active = true
		for _, handler := range tf.onFocusHandlers {
			handler()
		}
		tf.ReRender()
	}
}

func (tf *TextField) Deactivate() {
	if tf.active {
		tf.active = false
		for _, handler := range tf.onBlurHandlers {
			handler()
		}
		tf.ReRender()
	}
}

func (tf *TextField) AddOnFocus(handler func()) {
	tf.onFocusHandlers = append(tf.onFocusHandlers, handler)
}

func (tf *TextField) AddOnBlur(handler func()) {
	tf.onBlurHandlers = append(tf.onBlurHandlers, handler)
}

func (tf *TextField) AddOnKeyPress(handler func(key string, release bool)) {
	tf.AddOnKeyPress(handler)
}

func (tf *TextField) AddOnChange(handler func(string)) {
	tf.onChangeHandlers = append(tf.onChangeHandlers, handler)
}

func NewTextField(text string, textColor color.Color, textSize float32, textFont *truetype.Font) *TextField {
	cursor := renderer.CreateBoxWithOffset(0.07, 1.1, 0, 0.1)
	cursor.SetColor(color.NRGBA{0, 0, 0, 255})
	cursorNode := renderer.NewNode()
	cursorNode.Material = renderer.NewMaterial()
	cursorNode.Add(cursor)
	tf := &TextField{
		container: NewContainer(),
		text:      NewTextElement(text, textColor, textSize, textFont),
		cursor:    cursorNode,
		cursorPos: len(text),
	}
	tf.text.node.Add(cursorNode)
	tf.container.AddChildren(tf.text)
	tf.GetHitbox().AddOnClick(func(button int, release bool, position mgl32.Vec2) {
		if !release {
			tf.Activate()
		}
	})
	tf.SetBackgroundColor(0, 0, 0, 0)
	tf.SetHeight(textSize * 1.5)
	tf.SetPadding(NewMargin(2))
	return tf
}
