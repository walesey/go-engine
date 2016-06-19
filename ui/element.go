package ui

import (
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

type Element interface {
	Render(size, offset vmath.Vector2) vmath.Vector2
	ReRender()
	Spatial() renderer.Spatial
	GetId() string
	mouseMove(position vmath.Vector2)
	mouseClick(button int, release bool, position vmath.Vector2)
	keyClick(key string, release bool)
}
