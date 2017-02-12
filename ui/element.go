package ui

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/renderer"
)

type Element interface {
	Render(size, offset mgl32.Vec2) mgl32.Vec2
	ReRender()
	Spatial() renderer.Spatial
	GetId() string
	SetId(id string)
	GetChildren() Children
	mouseMove(position mgl32.Vec2) bool
	mouseClick(button int, release bool, position mgl32.Vec2) bool
	keyClick(key string, release bool)
}
