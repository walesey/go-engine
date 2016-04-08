package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"

	"github.com/disintegration/imaging"
	"github.com/walesey/freetype"
	"github.com/walesey/freetype/truetype"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
	"golang.org/x/image/font"
)

type TextElement struct {
	img       *ImageElement
	text      string
	textColor color.Color
	textSize  float64
	useFont   *truetype.Font
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

func (te *TextElement) UpdateImage(dimensions vmath.Vector2) {
	// Initialize the context.
	bg := image.Transparent
	c := freetype.NewContext()
	c.SetDPI(4 * te.textSize)
	c.SetFont(te.useFont)
	c.SetFontSize(te.textSize)
	c.SetSrc(image.NewUniform(te.textColor))
	c.SetHinting(font.HintingNone)

	// Establish image dimensions
	textHeight := c.PointToFixed(te.textSize)
	// var width fixed.Int26_6
	// var height fixed.Int26_6
	// for _, s := range text {
	// dimensions, _ := c.StringDimensions(te.text)
	// height = height + dimensions.Y
	// if dimensions.X > width {
	// 	width = dimensions.X
	// }
	// }
	// imgWidth := int(width >> 6)
	// imgHeight := int(height >> 6)
	fmt.Println(dimensions)
	rgba := image.NewRGBA(image.Rect(0, 0, int(dimensions.X), int(dimensions.Y)))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)

	// Draw the text.
	pt := freetype.Pt(0, int(textHeight>>6))
	// for _, s := range text {
	_, err := c.DrawString(te.text, pt)
	if err != nil {
		log.Printf("Error drawing string: %v\n", err)
		return
	}
	pt.Y += textHeight
	// }

	te.img.SetImage(imaging.FlipV(rgba))
}

func (te *TextElement) SetSize(width, height int) {
	te.img.SetSize(width, height)
}

func (te *TextElement) Render(offset vmath.Vector2) vmath.Vector2 {
	dimensions := te.img.Render(offset)
	te.UpdateImage(dimensions)
	return dimensions
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

func NewTextElement(text string, textColor color.Color, size float64) *TextElement {
	useFont, _ := LoadFont("TestAssets/luximr.ttf")
	textElem := &TextElement{
		img:       NewImageElement(image.NewAlpha(image.Rect(0, 0, 1, 1))),
		text:      text,
		textColor: textColor,
		textSize:  size,
		useFont:   useFont,
	}
	textElem.SetSize(300, 100)
	return textElem
}
