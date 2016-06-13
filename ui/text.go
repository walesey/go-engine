package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/walesey/go-engine/libs/freetype"
	"github.com/walesey/go-engine/libs/freetype/truetype"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
	"golang.org/x/image/font"
)

const (
	LEFT_ALIGN int = iota
	CENTER_ALIGN
	RIGHT_ALIGN
)

func LoadFontFromFile(fontfile string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Printf("Error Reading from font file: %v\n", err)
		return nil, err
	}
	return LoadFont(fontBytes)
}

func LoadFont(fontData []byte) (*truetype.Font, error) {
	f, err := truetype.Parse(fontData)
	if err != nil {
		log.Printf("Error Parsing font file: %v\n", err)
		return nil, err
	}
	return f, nil
}

type TextElement struct {
	id                 string
	node               *renderer.Node
	cursor             *renderer.Node
	img                *ImageElement
	width, height      float64
	text               string
	textColor          color.Color
	textSize           float64
	textFont           *truetype.Font
	textAlign          int
	size, offset       vmath.Vector2
	active             bool
	cursorPos          int
	onFocusHandlers    []func()
	onBlurHandlers     []func()
	onKeyPressHandlers []func(key string, release bool)
	dirty              bool
}

func (te *TextElement) getContext() *freetype.Context {
	c := freetype.NewContext()
	c.SetDPI(75)
	c.SetFont(te.textFont)
	c.SetFontSize(te.textSize)
	c.SetSrc(image.NewUniform(te.textColor))
	c.SetHinting(font.HintingNone)
	return c
}

func (te *TextElement) updateImage(size vmath.Vector2) {
	// Initialize the context.
	bg := image.Transparent
	c := te.getContext()

	// Establish image dimensions and do word wrap
	textHeight := c.PointToFixed(te.textSize)
	var width int
	var height int = int(textHeight >> 6)
	words := strings.Split(te.text, " ")
	lines := []string{""}
	lineNb := 0
	for _, word := range words {
		wordWithSpace := fmt.Sprintf("%v ", word)
		dimensions, _ := c.StringDimensions(wordWithSpace)
		width += int(dimensions.X >> 6)
		if width > int(size.X) {
			width = int(dimensions.X >> 6)
			height += int(dimensions.Y>>6) + 1
			lines = append(lines, "")
			lineNb += 1
		}
		lines[lineNb] = fmt.Sprintf("%v%v", lines[lineNb], wordWithSpace)
	}
	if te.height > 0 {
		height = int(te.height)
	}

	rgba := image.NewRGBA(image.Rect(0, 0, int(size.X), height+int(textHeight>>6)/3))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)

	// Draw the text.
	pt := freetype.Pt(0, int(textHeight>>6))
	for _, line := range lines {
		_, err := c.DrawString(line, pt)
		if err != nil {
			log.Printf("Error drawing string: %v\n", err)
			return
		}
		pt.Y += textHeight
	}

	te.img.SetImage(imaging.FlipV(rgba))
	te.img.SetWidth(float64(rgba.Bounds().Size().X))
	te.img.SetHeight(float64(rgba.Bounds().Size().Y))
}

func (te *TextElement) GetText() string {
	return te.text
}

func (te *TextElement) SetText(text string) {
	te.text = text
	te.cursorPos = len(text)
	te.dirty = true
}

func (te *TextElement) SetFont(textFont *truetype.Font) {
	te.textFont = textFont
	te.dirty = true
}

func (te *TextElement) SetTextSize(textSize float64) {
	te.textSize = textSize
	te.dirty = true
}

func (te *TextElement) SetTextColor(textColor color.Color) {
	te.textColor = textColor
	te.dirty = true
}

func (te *TextElement) SetWidth(width float64) {
	te.width = width
	te.dirty = true
}

func (te *TextElement) SetHeight(height float64) {
	te.height = height
	te.dirty = true
}

func (te *TextElement) Active() bool {
	return te.active
}

func (te *TextElement) Activate() {
	if !te.active {
		te.active = true
		for _, handler := range te.onFocusHandlers {
			handler()
		}
		te.ReRender()
	}
}

func (te *TextElement) Deactivate() {
	if te.active {
		te.active = false
		for _, handler := range te.onBlurHandlers {
			handler()
		}
		te.ReRender()
	}
}

func (te *TextElement) Render(size, offset vmath.Vector2) vmath.Vector2 {
	useWidth := size.X
	useHeight := size.Y
	if te.width > 0 {
		useWidth = te.width
	}
	if te.height > 0 {
		useHeight = te.height
	}
	useSize := vmath.Vector2{useWidth, useHeight}
	te.size = useSize
	te.offset = offset
	if te.dirty {
		te.updateImage(useSize)
	}
	te.dirty = false
	renderSize := te.img.Render(size, offset)
	if te.textAlign == CENTER_ALIGN {
		te.img.Render(size, offset.Add(vmath.Vector2{(size.X - renderSize.X) * 0.5, 0}))
	}
	if te.textAlign == RIGHT_ALIGN {
		te.img.Render(size, offset.Add(vmath.Vector2{size.X - renderSize.X, 0}))
	}
	te.RenderCursor()
	return renderSize
}

func (te *TextElement) ReRender() {
	te.Render(te.size, te.offset)
}

func (te *TextElement) RenderCursor() {
	c := te.getContext()
	cursorTranslation, _ := c.StringDimensions(string([]byte(te.text)[:te.cursorPos]))
	xPos := int(cursorTranslation.X >> 6)
	te.cursor.SetTranslation(vmath.Vector2{float64(xPos), 0}.ToVector3())
	if te.active {
		te.cursor.SetScale(vmath.Vector2{te.textSize, te.textSize}.ToVector3())
	} else {
		te.cursor.SetScale(vmath.Vector2{0, 0}.ToVector3())
	}
}

func (te *TextElement) Spatial() renderer.Spatial {
	return te.node
}

func (te *TextElement) GetId() string {
	return te.id
}

func (te *TextElement) mouseMove(position vmath.Vector2) {
	te.img.mouseMove(position)
}

func (te *TextElement) mouseClick(button int, release bool, position vmath.Vector2) {
	te.img.mouseClick(button, release, position)
}

func (te *TextElement) keyClick(key string, release bool) {
	if te.active && !release {
		textBytes := []byte(te.text)
		if key == "backspace" {
			if len(textBytes) > 0 && te.cursorPos > 0 {
				cursorPos := te.cursorPos
				newText := append(textBytes[:cursorPos-1], textBytes[cursorPos:]...)
				te.SetText(string(newText))
				te.cursorPos = cursorPos - 1
			}
		} else if key == "leftArrow" {
			if te.cursorPos > 0 {
				te.cursorPos--
			}
		} else if key == "rightArrow" {
			if te.cursorPos < len(te.text) {
				te.cursorPos++
			}
		} else {
			insertText := []rune(key)
			newText := []rune(te.text)
			newText = append(newText[:te.cursorPos], append(insertText, newText[te.cursorPos:]...)...)
			cursorPos := te.cursorPos
			te.SetText(string(newText))
			te.cursorPos = cursorPos + 1
		}
		for _, handler := range te.onKeyPressHandlers {
			handler(key, release)
		}
		te.ReRender()
	}
}

func (te *TextElement) SetAlign(align int) {
	te.textAlign = align
}

func (te *TextElement) AddOnFocus(handler func()) {
	te.onFocusHandlers = append(te.onFocusHandlers, handler)
}

func (te *TextElement) AddOnBlur(handler func()) {
	te.onBlurHandlers = append(te.onBlurHandlers, handler)
}

func (te *TextElement) AddOnKeyPress(handler func(key string, release bool)) {
	te.onKeyPressHandlers = append(te.onKeyPressHandlers, handler)
}

func NewTextElement(text string, textColor color.Color, textSize float64, textFont *truetype.Font) *TextElement {
	img := NewImageElement(image.NewAlpha(image.Rect(0, 0, 1, 1)))
	node := renderer.CreateNode()
	cursor := renderer.CreateBoxWithOffset(0.07, 1.1, 0, 0.1)
	mat := renderer.CreateMaterial()
	mat.LightingMode = renderer.MODE_UNLIT
	cursor.Material = mat
	cursor.SetColor(color.NRGBA{0, 0, 0, 255})
	cursorNode := renderer.CreateNode()
	cursorNode.Add(cursor)
	node.Add(cursorNode)
	node.Add(img.Spatial())
	textElem := &TextElement{
		img:       img,
		cursor:    cursorNode,
		node:      node,
		text:      text,
		textColor: textColor,
		textSize:  textSize,
		textFont:  textFont,
		dirty:     true,
	}
	return textElem
}
