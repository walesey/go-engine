package ui

import "github.com/go-gl/mathgl/mgl32"

type EventHandler struct {
	onClickHandlers     []func(button int, release bool, position mgl32.Vec2)
	onHoverHandlers     []func()
	onUnHoverHandlers   []func()
	onMouseMoveHandlers []func(position mgl32.Vec2)
}

func (eh *EventHandler) AddOnClick(handler func(button int, release bool, position mgl32.Vec2)) {
	eh.onClickHandlers = append(eh.onClickHandlers, handler)
}

func (eh *EventHandler) AddOnHover(handler func()) {
	eh.onHoverHandlers = append(eh.onHoverHandlers, handler)
}

func (eh *EventHandler) AddOnUnHover(handler func()) {
	eh.onUnHoverHandlers = append(eh.onUnHoverHandlers, handler)
}

func (eh *EventHandler) AddOnMouseMove(handler func(position mgl32.Vec2)) {
	eh.onMouseMoveHandlers = append(eh.onMouseMoveHandlers, handler)
}

func (eh *EventHandler) onClick(button int, release bool, position mgl32.Vec2) {
	for _, handler := range eh.onClickHandlers {
		handler(button, release, position)
	}
}

func (eh *EventHandler) onHover() {
	for _, handler := range eh.onHoverHandlers {
		handler()
	}
}

func (eh *EventHandler) onUnHover() {
	for _, handler := range eh.onUnHoverHandlers {
		handler()
	}
}

func (eh *EventHandler) onMouseMove(position mgl32.Vec2) {
	for _, handler := range eh.onMouseMoveHandlers {
		handler(position)
	}
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		onClickHandlers:     make([]func(button int, release bool, position mgl32.Vec2), 0),
		onHoverHandlers:     make([]func(), 0),
		onUnHoverHandlers:   make([]func(), 0),
		onMouseMoveHandlers: make([]func(position mgl32.Vec2), 0),
	}
}
