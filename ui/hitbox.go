package ui

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/util"
)

type Hitbox interface {
	AddOnClick(handler func(button int, release bool, position mgl32.Vec2))
	AddOnHover(handler func())
	AddOnUnHover(handler func())
	AddOnMouseMove(handler func(position mgl32.Vec2))
	MouseMove(position mgl32.Vec2)
	MouseClick(button int, release bool, position mgl32.Vec2)
	SetSize(size mgl32.Vec2)
}

type HitboxImpl struct {
	size         mgl32.Vec2
	eventHandler *EventHandler
	hoverState   bool
}

func (hb *HitboxImpl) AddOnClick(handler func(button int, release bool, position mgl32.Vec2)) {
	hb.eventHandler.AddOnClick(handler)
}

func (hb *HitboxImpl) AddOnHover(handler func()) {
	hb.eventHandler.AddOnHover(handler)
}

func (hb *HitboxImpl) AddOnUnHover(handler func()) {
	hb.eventHandler.AddOnUnHover(handler)
}

func (hb *HitboxImpl) AddOnMouseMove(handler func(position mgl32.Vec2)) {
	hb.eventHandler.AddOnMouseMove(handler)
}

func (hb *HitboxImpl) MouseMove(position mgl32.Vec2) {
	if util.PointLiesInsideAABB(mgl32.Vec2{}, hb.size, position) {
		if !hb.hoverState {
			hb.hoverState = true
			hb.eventHandler.onHover()
		}
		hb.eventHandler.onMouseMove(position)
	} else if hb.hoverState {
		hb.hoverState = false
		hb.eventHandler.onUnHover()
	}
}

func (hb *HitboxImpl) MouseClick(button int, release bool, position mgl32.Vec2) {
	if util.PointLiesInsideAABB(mgl32.Vec2{}, hb.size, position) {
		hb.eventHandler.onClick(button, release, position)
	}
}

func (hb *HitboxImpl) SetSize(size mgl32.Vec2) {
	hb.size = size
}

func NewHitbox() Hitbox {
	return &HitboxImpl{
		eventHandler: NewEventHandler(),
	}
}
