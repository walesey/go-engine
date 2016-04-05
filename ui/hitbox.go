package ui

import vmath "github.com/walesey/go-engine/vectormath"

type Hitbox interface {
	AddOnClick(handler func(button int, release bool, position vmath.Vector2))
	AddOnHover(handler func())
	AddOnUnHover(handler func())
	AddOnMouseMove(handler func(position vmath.Vector2))
	MouseMove(position vmath.Vector2)
	MouseClick(button int, release bool, position vmath.Vector2)
	SetSize(size vmath.Vector2)
}

type HitboxImpl struct {
	size         vmath.Vector2
	eventHandler *EventHandler
	hoverState   bool
}

func (hb *HitboxImpl) AddOnClick(handler func(button int, release bool, position vmath.Vector2)) {
	hb.eventHandler.AddOnClick(handler)
}

func (hb *HitboxImpl) AddOnHover(handler func()) {
	hb.eventHandler.AddOnHover(handler)
}

func (hb *HitboxImpl) AddOnUnHover(handler func()) {
	hb.eventHandler.AddOnUnHover(handler)
}

func (hb *HitboxImpl) AddOnMouseMove(handler func(position vmath.Vector2)) {
	hb.eventHandler.AddOnMouseMove(handler)
}

func (hb *HitboxImpl) MouseMove(position vmath.Vector2) {
	if vmath.PointLiesInsideAABB(vmath.Vector2{}, hb.size, position) {
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

func (hb *HitboxImpl) MouseClick(button int, release bool, position vmath.Vector2) {
	if vmath.PointLiesInsideAABB(vmath.Vector2{}, hb.size, position) {
		hb.eventHandler.onClick(button, release, position)
	}
}

func (hb *HitboxImpl) SetSize(size vmath.Vector2) {
	hb.size = size
}

func NewHitbox() Hitbox {
	return &HitboxImpl{
		eventHandler: NewEventHandler(),
	}
}
