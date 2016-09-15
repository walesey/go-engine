package glfwController

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/controller"
)

type Controller interface {
	controller.Controller
	KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey)
	CursorPosCallback(window *glfw.Window, xpos, ypos float32)
	ScrollCallback(window *glfw.Window, xoffset, yoffset float32)
}

type KeyAction struct {
	key    controller.Key
	action controller.Action
}

type MouseButtonAction struct {
	button controller.MouseButton
	action controller.Action
}

type ActionMap struct {
	keyAction            func(key controller.Key, action controller.Action)
	mouseAction          func(button controller.MouseButton, action controller.Action)
	keyActionMap         map[KeyAction][]func()
	mouseButtonActionMap map[MouseButtonAction][]func()
	axisActions          []func(xpos, ypos float32)
	scrollActions        []func(xoffset, yoffset float32)
}

func NewActionMap() controller.Controller {
	am := &ActionMap{
		keyAction:            func(key controller.Key, action controller.Action) {},
		mouseAction:          func(button controller.MouseButton, action controller.Action) {},
		keyActionMap:         make(map[KeyAction][]func()),
		mouseButtonActionMap: make(map[MouseButtonAction][]func()),
		axisActions:          make([]func(xpos, ypos float32), 0, 0),
		scrollActions:        make([]func(xoffset, yoffset float32), 0, 0),
	}
	return am
}

func (am *ActionMap) SetKeyAction(function func(key controller.Key, action controller.Action)) {
	am.keyAction = function
}

func (am *ActionMap) SetMouseAction(function func(button controller.MouseButton, action controller.Action)) {
	am.mouseAction = function
}

//Bindings
func (am *ActionMap) BindKeyAction(function func(), key controller.Key, action controller.Action) {
	ka := KeyAction{key, action}
	if m, ok := am.keyActionMap[ka]; ok {
		am.keyActionMap[ka] = append(m, function)
	} else {
		am.keyActionMap[ka] = []func(){function}
	}
}

func (am *ActionMap) BindMouseAction(function func(), button controller.MouseButton, action controller.Action) {
	mba := MouseButtonAction{button, action}
	if m, ok := am.mouseButtonActionMap[mba]; ok {
		am.mouseButtonActionMap[mba] = append(m, function)
	} else {
		am.mouseButtonActionMap[mba] = []func(){function}
	}
}

func (am *ActionMap) BindAxisAction(function func(xpos, ypos float32)) {
	am.axisActions = append(am.axisActions, function)
}

func (am *ActionMap) BindScrollAction(function func(xoffset, yoffset float32)) {
	am.scrollActions = append(am.scrollActions, function)
}

//Callbacks
func (am *ActionMap) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	am.keyAction(getKey(key), getAction(action))
	ka := KeyAction{getKey(key), getAction(action)}
	if m, ok := am.keyActionMap[ka]; ok {
		for _, function := range m {
			function()
		}
	}
}

func (am *ActionMap) MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	am.mouseAction(getMouseButton(button), getAction(action))
	mba := MouseButtonAction{getMouseButton(button), getAction(action)}
	if m, ok := am.mouseButtonActionMap[mba]; ok {
		for _, function := range m {
			function()
		}
	}
}

func (am *ActionMap) CursorPosCallback(window *glfw.Window, xpos, ypos float32) {
	for _, action := range am.axisActions {
		action(xpos, ypos)
	}
}

func (am *ActionMap) ScrollCallback(window *glfw.Window, xoffset, yoffset float32) {
	for _, action := range am.scrollActions {
		action(xoffset, yoffset)
	}
}
