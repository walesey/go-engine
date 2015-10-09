package controller

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Controller interface {
	KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	CursorPosCallback(window *glfw.Window, xpos, ypos float64)
	ScrollCallback(window *glfw.Window, xoffset, yoffset float64)
}

type KeyAction struct {
	key    glfw.Key
	action glfw.Action
}

type ActionMap struct {
	keyActionMap  map[KeyAction]func()
	axisActions   []func(xpos, ypos float64)
	scrollActions []func(xoffset, yoffset float64)
}

func NewActionMap() *ActionMap {
	am := &ActionMap{
		keyActionMap:  make(map[KeyAction]func()),
		axisActions:   make([]func(xpos, ypos float64), 0, 0),
		scrollActions: make([]func(xoffset, yoffset float64), 0, 0),
	}
	return am
}

//Bindings
func (am *ActionMap) BindAction(function func(), key glfw.Key, action glfw.Action) {
	ka := KeyAction{key, action}
	am.keyActionMap[ka] = function
}

func (am *ActionMap) BindAxisAction(function func(xpos, ypos float64)) {
	am.axisActions = append(am.axisActions, function)
}

func (am *ActionMap) BindScrollAction(function func(xoffset, yoffset float64)) {
	am.scrollActions = append(am.scrollActions, function)
}

//Callbacks
func (am *ActionMap) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	ka := KeyAction{key, action}
	if am.keyActionMap[ka] != nil {
		am.keyActionMap[ka]()
	}
}

func (am *ActionMap) CursorPosCallback(window *glfw.Window, xpos, ypos float64) {
	for _, action := range am.axisActions {
		action(xpos, ypos)
	}
}

func (am *ActionMap) ScrollCallback(window *glfw.Window, xoffset, yoffset float64) {
	for _, action := range am.scrollActions {
		action(xoffset, yoffset)
	}
}
