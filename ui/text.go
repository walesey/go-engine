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

const ( // TODO: implement text align
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

type textProps struct {
	width, height float32
	size, offset  mgl32.Vec2
	text          string
	placeholder   string
	textColor     color.Color
	textSize      float32
	textFont      *truetype.Font
	textAlign     int
	hidden        bool
}

type TextElement struct {
	id                   string
	node                 *renderer.Node
	img                  *ImageElement
	props, previousProps textProps
	onKeyPressHandlers   []func(key string, release bool)
}

func (te *TextElement) getContext() *freetype.Context {
	c := freetype.NewContext()
	c.SetDPI(75)
	c.SetFont(te.props.textFont)
	c.SetFontSize(float64(te.props.textSize))
	c.SetSrc(image.NewUniform(te.props.textColor))
	c.SetHinting(font.HintingNone)
	return c
}

func (te *TextElement) updateImage(size mgl32.Vec2) {
	// Initialize the context.
	bg := image.Transparent
	c := te.getContext()

	text := te.GetHiddenText()
	if len(text) == 0 {
		text = te.props.placeholder
		r, g, b, _ := te.props.textColor.RGBA()
		placeholderColor := color.RGBA{uint8(r), uint8(g), uint8(b), 80}
		c.SetSrc(image.NewUniform(placeholderColor))
	}

	// Establish image dimensions and do word wrap
	textHeight := c.PointToFixed(float64(te.props.textSize))
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
	if te.props.height > 0 {
		height = int(te.props.height)
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
	return te.props.text
}

func (te *TextElement) GetHiddenText() string {
	if te.props.hidden {
		text := ""
		for i := 0; i < len(te.props.text); i++ {
			text = fmt.Sprint(text, "*")
		}
		return text
	}
	return te.props.text
}

func (te *TextElement) SetText(text string) *TextElement {
	te.props.text = text
	return te
}

func (te *TextElement) SetPlaceholder(placeholder string) *TextElement {
	te.props.placeholder = placeholder
	return te
}

func (te *TextElement) SetFont(textFont *truetype.Font) *TextElement {
	te.props.textFont = textFont
	return te
}

func (te *TextElement) SetTextSize(textSize float32) *TextElement {
	te.props.textSize = textSize
	return te
}

func (te *TextElement) SetTextColor(textColor color.Color) *TextElement {
	te.props.textColor = textColor
	return te
}

func (te *TextElement) SetWidth(width float32) *TextElement {
	te.props.width = width
	return te
}

func (te *TextElement) SetHeight(height float32) *TextElement {
	te.props.height = height
	return te
}

func (te *TextElement) SetHidden(hidden bool) *TextElement {
	te.props.hidden = hidden
	return te
}

func (te *TextElement) Render(size, offset mgl32.Vec2) mgl32.Vec2 {
	te.props.size, te.props.offset = size, offset
	textWidth, textHeight := size.X(), size.Y()
	if te.props.width > 0 {
		textWidth = te.props.width
	}
	if te.props.height > 0 {
		textHeight = te.props.height
	}
	if te.previousProps != te.props {
		te.updateImage(mgl32.Vec2{textWidth, textHeight})
		te.previousProps = te.props
	}
	return te.img.Render(size, offset)
}

func (te *TextElement) ReRender() {
	te.Render(te.props.size, te.props.offset)
}

func (te *TextElement) Spatial() renderer.Spatial {
	return te.node
}

func (te *TextElement) GetId() string {
	return te.id
}

func (te *TextElement) SetId(id string) {
	te.id = id
}

func (te *TextElement) GetChildren() Children {
	return []Element{}
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
	te.props.textAlign = align
}

func (te *TextElement) AddOnKeyPress(handler func(key string, release bool)) {
	te.onKeyPressHandlers = append(te.onKeyPressHandlers, handler)
}

func NewTextElement(text string, textColor color.Color, textSize float32, textFont *truetype.Font) *TextElement {
	img := NewImageElement(image.NewAlpha(image.Rect(0, 0, 1, 1)))
	node := renderer.NewNode()
	node.Add(img.Spatial())
	textElem := &TextElement{
		img:  img,
		node: node,
		props: textProps{
			text:      text,
			textColor: textColor,
			textSize:  textSize,
			textFont:  textFont,
		},
	}
	if textFont == nil {
		defaultFont, _ := LoadFont(getDefaultFont())
		textElem.SetFont(defaultFont)
	}
	return textElem
}
