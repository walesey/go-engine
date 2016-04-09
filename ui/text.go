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
	"github.com/walesey/freetype"
	"github.com/walesey/freetype/truetype"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
	"golang.org/x/image/font"
)

type TextElement struct {
	img           *ImageElement
	width, height float64
	text          string
	textColor     color.Color
	textSize      float64
	textFont      *truetype.Font
}

func LoadFont(fontfile string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Printf("Error Reading from font file: %v\n", err)
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Printf("Error Parsing font file: %v\n", err)
		return nil, err
	}
	return f, nil
}

func (te *TextElement) updateImage(size vmath.Vector2) {
	// Initialize the context.
	bg := image.White
	c := freetype.NewContext()
	c.SetDPI(75)
	c.SetFont(te.textFont)
	c.SetFontSize(te.textSize)
	c.SetSrc(image.NewUniform(te.textColor))
	c.SetHinting(font.HintingNone)

	// Establish image dimensions and do ward wrap
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

	rgba := image.NewRGBA(image.Rect(0, 0, int(size.X), height+1))
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
	te.img.SetSize(float64(rgba.Bounds().Size().X), float64(rgba.Bounds().Size().Y))
}

func (te *TextElement) SetText(text string) {
	te.text = text
}

func (te *TextElement) SetSize(width, height float64) {
	te.width, te.height = width, height
}

func (te *TextElement) Render(size, offset vmath.Vector2) vmath.Vector2 {
	if te.width > 0 {
		te.updateImage(vmath.Vector2{float64(te.width), size.Y})
	} else {
		te.updateImage(size)
	}
	return te.img.Render(size, offset)
}

func (te *TextElement) Spatial() renderer.Spatial {
	return te.img.Spatial()
}

func (te *TextElement) mouseMove(position vmath.Vector2) {
	te.img.mouseMove(position)
}

func (te *TextElement) mouseClick(button int, release bool, position vmath.Vector2) {
	te.img.mouseClick(button, release, position)
}

func NewTextElement(text string, textColor color.Color, textSize float64, textFont *truetype.Font) *TextElement {
	textElem := &TextElement{
		img:       NewImageElement(image.NewAlpha(image.Rect(0, 0, 1, 1))),
		text:      text,
		textColor: textColor,
		textSize:  textSize,
		textFont:  textFont,
	}
	return textElem
}
