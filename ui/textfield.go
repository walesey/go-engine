package ui

import (
	"image"
	"image/color"

	"github.com/walesey/go-engine/libs/freetype/truetype"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type TextField struct {
	container       *Container
	text            *TextElement
	cursor          *renderer.Node
	cursorPos       int
	active          bool
	onFocusHandlers []func()
	onBlurHandlers  []func()
}

func (tf *TextField) Render(size, offset vmath.Vector2) vmath.Vector2 {
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
	tf.cursor.SetTranslation(vmath.Vector2{float64(xPos), 0}.ToVector3())
	if tf.active {
		tf.cursor.SetScale(vmath.Vector2{tf.text.textSize, tf.text.textSize}.ToVector3())
	} else {
		tf.cursor.SetScale(vmath.Vector2{0, 0}.ToVector3())
	}
}

func (tf *TextField) Spatial() renderer.Spatial {
	return tf.container.Spatial()
}

func (tf *TextField) GetId() string {
	return tf.container.GetId()
}

func (tf *TextField) mouseMove(position vmath.Vector2) {
	tf.container.mouseMove(position)
}

func (tf *TextField) mouseClick(button int, release bool, position vmath.Vector2) {
	tf.container.mouseClick(button, release, position)
}

func (tf *TextField) keyClick(key string, release bool) {
	if tf.active && !release {
		textBytes := []byte(tf.text.text)
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
			if tf.cursorPos < len(tf.text.text) {
				tf.cursorPos++
			}
		} else {
			insertText := []rune(key)
			newText := []rune(tf.text.text)
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

func (tf *TextField) SetTextSize(textSize float64) *TextField {
	tf.text.SetTextSize(textSize)
	return tf
}

func (tf *TextField) SetTextColor(textColor color.Color) *TextField {
	tf.text.SetTextColor(textColor)
	return tf
}

func (tf *TextField) SetWidth(width float64) {
	tf.container.SetWidth(width)
}

func (tf *TextField) UsePercentWidth(usePercent bool) {
	tf.container.UsePercentHeight(usePercent)
}

func (tf *TextField) SetHeight(height float64) {
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

func NewTextField(text string, textColor color.Color, textSize float64, textFont *truetype.Font) *TextField {
	cursor := renderer.CreateBoxWithOffset(0.07, 1.1, 0, 0.1)
	mat := renderer.CreateMaterial()
	mat.LightingMode = renderer.MODE_UNLIT
	cursor.Material = mat
	cursor.SetColor(color.NRGBA{0, 0, 0, 255})
	cursorNode := renderer.CreateNode()
	cursorNode.Add(cursor)
	tf := &TextField{
		container: NewContainer(),
		text:      NewTextElement(text, textColor, textSize, textFont),
		cursor:    cursorNode,
	}
	tf.text.node.Add(cursorNode)
	tf.container.AddChildren(tf.text)
	tf.GetHitbox().AddOnClick(func(button int, release bool, position vmath.Vector2) {
		if !release {
			tf.Activate()
		}
	})
	tf.SetBackgroundColor(0, 0, 0, 0)
	tf.SetHeight(textSize * 1.5)
	tf.SetPadding(NewMargin(2))
	return tf
}
