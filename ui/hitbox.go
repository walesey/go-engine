package ui

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/emitter"
	"github.com/walesey/go-engine/util"
)

type Hitbox interface {
	AddOnClick(handler func(button int, release bool, position mgl32.Vec2))
	AddOnHover(handler func())
	AddOnUnHover(handler func())
	AddOnMouseMove(handler func(position mgl32.Vec2))
	MouseMove(position mgl32.Vec2) bool
	MouseClick(button int, release bool, position mgl32.Vec2) bool
	SetSize(size mgl32.Vec2)
}

type HitboxImpl struct {
	size       mgl32.Vec2
	events     emitter.EventEmitter
	hoverState bool
}

type clickEvent struct {
	button   int
	release  bool
	position mgl32.Vec2
}

func (hb *HitboxImpl) AddOnClick(handler func(button int, release bool, position mgl32.Vec2)) {
	hb.events.On("click", func(e emitter.Event) {
		if ce, ok := e.(clickEvent); ok {
			handler(ce.button, ce.release, ce.position)
		}
	})
}

func (hb *HitboxImpl) AddOnHover(handler func()) {
	hb.events.On("hover", func(e emitter.Event) {
		handler()
	})
}

func (hb *HitboxImpl) AddOnUnHover(handler func()) {
	hb.events.On("unHover", func(e emitter.Event) {
		handler()
	})
}

func (hb *HitboxImpl) AddOnMouseMove(handler func(position mgl32.Vec2)) {
	hb.events.On("mouseMove", func(e emitter.Event) {
		if pos, ok := e.(mgl32.Vec2); ok {
			handler(pos)
		}
	})
}

func (hb *HitboxImpl) MouseMove(position mgl32.Vec2) bool {
	if util.PointLiesInsideAABB(mgl32.Vec2{}, hb.size, position) {
		if !hb.hoverState {
			hb.hoverState = true
			hb.events.Do("hover", 1)
		}
		hb.events.Do("mouseMove", position)
		return true
	} else if hb.hoverState {
		hb.hoverState = false
		hb.events.Do("unHover", 1)
	}
	return false
}

func (hb *HitboxImpl) MouseClick(button int, release bool, position mgl32.Vec2) bool {
	if util.PointLiesInsideAABB(mgl32.Vec2{}, hb.size, position) {
		hb.events.Do("click", clickEvent{button, release, position})
		return true
	}
	return false
}

func (hb *HitboxImpl) SetSize(size mgl32.Vec2) {
	hb.size = size
}

func NewHitbox() Hitbox {
	return &HitboxImpl{
		events: emitter.New(1),
	}
}
