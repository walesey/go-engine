package glfwController

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/controller"
)

type Controller interface {
	controller.Controller
	KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)
	MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey)
	CursorPosCallback(window *glfw.Window, xpos, ypos float64)
	ScrollCallback(window *glfw.Window, xoffset, yoffset float64)
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
	keyActionList        []func(key controller.Key, action controller.Action)
	keyActionMap         map[KeyAction]func()
	mouseButtonActionMap map[MouseButtonAction]func()
	axisActions          []func(xpos, ypos float64)
	scrollActions        []func(xoffset, yoffset float64)
}

func NewActionMap() controller.Controller {
	am := &ActionMap{
		keyActionList:        make([]func(key controller.Key, action controller.Action), 0, 0),
		keyActionMap:         make(map[KeyAction]func()),
		mouseButtonActionMap: make(map[MouseButtonAction]func()),
		axisActions:          make([]func(xpos, ypos float64), 0, 0),
		scrollActions:        make([]func(xoffset, yoffset float64), 0, 0),
	}
	return am
}

//Bindings
func (am *ActionMap) BindAction(function func(), key controller.Key, action controller.Action) {
	ka := KeyAction{key, action}
	am.keyActionMap[ka] = function
}

func (am *ActionMap) BindKeyAction(function func(key controller.Key, action controller.Action)) {
	am.keyActionList = append(am.keyActionList, function)
}

func (am *ActionMap) BindMouseAction(function func(), button controller.MouseButton, action controller.Action) {
	mba := MouseButtonAction{button, action}
	am.mouseButtonActionMap[mba] = function
}

func (am *ActionMap) BindAxisAction(function func(xpos, ypos float64)) {
	am.axisActions = append(am.axisActions, function)
}

func (am *ActionMap) BindScrollAction(function func(xoffset, yoffset float64)) {
	am.scrollActions = append(am.scrollActions, function)
}

//Callbacks
func (am *ActionMap) KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	ka := KeyAction{getKey(key), getAction(action)}
	if am.keyActionMap[ka] != nil {
		am.keyActionMap[ka]()
	}
	for _, function := range am.keyActionList {
		function(getKey(key), getAction(action))
	}
}

func (am *ActionMap) MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	mba := MouseButtonAction{getMouseButton(button), getAction(action)}
	if am.mouseButtonActionMap[mba] != nil {
		am.mouseButtonActionMap[mba]()
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
