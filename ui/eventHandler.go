package ui

import vmath "github.com/walesey/go-engine/vectormath"

type EventHandler struct {
	onClickHandlers     []func(button int, release bool, position vmath.Vector2)
	onHoverHandlers     []func()
	onUnHoverHandlers   []func()
	onMouseMoveHandlers []func(position vmath.Vector2)
}

func (eh *EventHandler) AddOnClick(handler func(button int, release bool, position vmath.Vector2)) {
	eh.onClickHandlers = append(eh.onClickHandlers, handler)
}

func (eh *EventHandler) AddOnHover(handler func()) {
	eh.onHoverHandlers = append(eh.onHoverHandlers, handler)
}

func (eh *EventHandler) AddOnUnHover(handler func()) {
	eh.onUnHoverHandlers = append(eh.onUnHoverHandlers, handler)
}

func (eh *EventHandler) AddOnMouseMove(handler func(position vmath.Vector2)) {
	eh.onMouseMoveHandlers = append(eh.onMouseMoveHandlers, handler)
}

func (eh *EventHandler) onClick(button int, release bool, position vmath.Vector2) {
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

func (eh *EventHandler) onMouseMove(position vmath.Vector2) {
	for _, handler := range eh.onMouseMoveHandlers {
		handler(position)
	}
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		onClickHandlers:     make([]func(button int, release bool, position vmath.Vector2), 0),
		onHoverHandlers:     make([]func(), 0),
		onUnHoverHandlers:   make([]func(), 0),
		onMouseMoveHandlers: make([]func(position vmath.Vector2), 0),
	}
}
