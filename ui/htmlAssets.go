package ui

import (
	"image"

	"github.com/walesey/freetype/truetype"
)

type HtmlAssets struct {
	fontMap     map[string]*truetype.Font
	callbackMap map[string]func(element Element, args ...interface{})
	imageMap    map[string]image.Image
}

func (assets HtmlAssets) AddFont(key string, font *truetype.Font) {
	assets.fontMap[key] = font
}

func (assets HtmlAssets) AddCallback(key string, callback func(element Element, args ...interface{})) {
	assets.callbackMap[key] = callback
}

func (assets HtmlAssets) AddImage(key string, img image.Image) {
	assets.imageMap[key] = img
}

func NewHtmlAssets() HtmlAssets {
	return HtmlAssets{
		fontMap:     make(map[string]*truetype.Font),
		callbackMap: make(map[string]func(element Element, args ...interface{})),
		imageMap:    make(map[string]image.Image),
	}
}
