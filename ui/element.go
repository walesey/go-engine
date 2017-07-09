package ui

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

type Element interface {
	Render(size, offset mgl32.Vec2) mgl32.Vec2
	ReRender()
	Spatial() (spatial renderer.Spatial)
	GlobalOrthoOrder() int
	GetId() string
	SetId(id string)
	GetChildren() Children
	mouseMove(position mgl32.Vec2) bool
	mouseClick(button int, release bool, position mgl32.Vec2) bool
	keyClick(key string, release bool)
}

// Sort elements
type byGlobalOrthoOrder []Element

func (slice byGlobalOrthoOrder) Len() int {
	return len(slice)
}

func (slice byGlobalOrthoOrder) Less(i, j int) bool {
	return slice[i].GlobalOrthoOrder() > slice[j].GlobalOrthoOrder()
}

func (slice byGlobalOrthoOrder) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
