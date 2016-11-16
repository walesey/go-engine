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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/libs/freetype"
	"github.com/walesey/go-engine/libs/freetype/truetype"
	"github.com/walesey/go-engine/renderer"

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
	img                *ImageElement
	width, height      float32
	text               string
	placeholder        string
	textColor          color.Color
	textSize           float32
	textFont           *truetype.Font
	textAlign          int
	size, offset       mgl32.Vec2
	hidden             bool
	onKeyPressHandlers []func(key string, release bool)
	dirty              bool
}

func (te *TextElement) getContext() *freetype.Context {
	c := freetype.NewContext()
	c.SetDPI(75)
	c.SetFont(te.textFont)
	c.SetFontSize(float64(te.textSize))
	c.SetSrc(image.NewUniform(te.textColor))
	c.SetHinting(font.HintingNone)
	return c
}

func (te *TextElement) updateImage(size mgl32.Vec2) {
	// Initialize the context.
	bg := image.Transparent
	c := te.getContext()

	text := te.GetHiddenText()
	if len(text) == 0 {
		text = te.placeholder
		r, g, b, _ := te.textColor.RGBA()
		placeholderColor := color.RGBA{uint8(r), uint8(g), uint8(b), 80}
		c.SetSrc(image.NewUniform(placeholderColor))
	}

	// Establish image dimensions and do word wrap
	textHeight := c.PointToFixed(float64(te.textSize))
	var width int
	var height int = int(textHeight >> 6)
	words := strings.Split(text, " ")
	lines := []string{""}
	lineNb := 0
	for _, word := range words {
		wordWithSpace := fmt.Sprintf("%v ", word)
		dimensions, _ := c.StringDimensions(wordWithSpace)
		width += int(dimensions.X >> 6)
		if width > int(size.X()) {
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

	rgba := image.NewRGBA(image.Rect(0, 0, int(size.X()), height+int(textHeight>>6)/3))
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
	te.img.SetWidth(float32(rgba.Bounds().Size().X))
	te.img.SetHeight(float32(rgba.Bounds().Size().Y))
}

func (te *TextElement) GetText() string {
	return te.text
}

func (te *TextElement) GetHiddenText() string {
	if te.hidden {
		text := ""
		for i := 0; i < len(te.text); i++ {
			text = fmt.Sprint(text, "*")
		}
		return text
	}
	return te.text
}

func (te *TextElement) SetText(text string) *TextElement {
	te.text = text
	te.dirty = true
	return te
}

func (te *TextElement) SetPlaceholder(placeholder string) *TextElement {
	te.placeholder = placeholder
	te.dirty = true
	return te
}

func (te *TextElement) SetFont(textFont *truetype.Font) *TextElement {
	te.textFont = textFont
	te.dirty = true
	return te
}

func (te *TextElement) SetTextSize(textSize float32) *TextElement {
	te.textSize = textSize
	te.dirty = true
	return te
}

func (te *TextElement) SetTextColor(textColor color.Color) *TextElement {
	te.textColor = textColor
	te.dirty = true
	return te
}

func (te *TextElement) SetWidth(width float32) *TextElement {
	te.width = width
	te.dirty = true
	return te
}

func (te *TextElement) SetHeight(height float32) *TextElement {
	te.height = height
	te.dirty = true
	return te
}

func (te *TextElement) SetHidden(hidden bool) *TextElement {
	te.hidden = hidden
	return te
}

func (te *TextElement) Render(size, offset mgl32.Vec2) mgl32.Vec2 {
	useWidth := size.X()
	useHeight := size.Y()
	if te.width > 0 {
		useWidth = te.width
	}
	if te.height > 0 {
		useHeight = te.height
	}
	useSize := mgl32.Vec2{useWidth, useHeight}
	te.size = useSize
	te.offset = offset
	if te.dirty {
		te.updateImage(useSize)
	}
	te.dirty = false
	renderSize := te.img.Render(size, offset)
	if te.textAlign == CENTER_ALIGN {
		te.img.Render(size, offset.Add(mgl32.Vec2{(size.X() - renderSize.X()) * 0.5, 0}))
	}
	if te.textAlign == RIGHT_ALIGN {
		te.img.Render(size, offset.Add(mgl32.Vec2{size.X() - renderSize.X(), 0}))
	}
	return renderSize
}

func (te *TextElement) ReRender() {
	te.Render(te.size, te.offset)
}

func (te *TextElement) Spatial() renderer.Spatial {
	return te.node
}

func (te *TextElement) GetId() string {
	return te.id
}

func (te *TextElement) mouseMove(position mgl32.Vec2) {
	te.img.mouseMove(position)
}

func (te *TextElement) mouseClick(button int, release bool, position mgl32.Vec2) {
	te.img.mouseClick(button, release, position)
}

func (te *TextElement) keyClick(key string, release bool) {
	te.img.keyClick(key, release)
}

func (te *TextElement) SetAlign(align int) {
	te.textAlign = align
}

func (te *TextElement) AddOnKeyPress(handler func(key string, release bool)) {
	te.onKeyPressHandlers = append(te.onKeyPressHandlers, handler)
}

func NewTextElement(text string, textColor color.Color, textSize float32, textFont *truetype.Font) *TextElement {
	img := NewImageElement(image.NewAlpha(image.Rect(0, 0, 1, 1)))
	node := renderer.NewNode()
	node.Add(img.Spatial())
	textElem := &TextElement{
		img:       img,
		node:      node,
		text:      text,
		textColor: textColor,
		textSize:  textSize,
		textFont:  textFont,
		dirty:     true,
	}
	if textFont == nil {
		defaultFont, _ := LoadFont(getDefaultFont())
		textElem.SetFont(defaultFont)
	}
	return textElem
}
