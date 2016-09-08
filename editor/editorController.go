package editor

import (
	"github.com/walesey/go-engine/controller"
	vmath "github.com/walesey/go-engine/vectormath"
)

const mouseSpeed = 0.01

func NewEditorController(e *Editor) controller.Controller {
	c := controller.CreateController()
	var mouseDown bool
	x, y := 0.0, 0.0
	axisLock := vmath.Vector3{1, 1, 1}

	doMouseMove := func(xpos, ypos float64) {
		if mouseDown {
			switch {
			case e.mouseMode == "scale":
				e.ScaleSelectedNodeModel(xpos-x, ypos-y, axisLock)
			case e.mouseMode == "translate":
				e.MoveSelectedNodeModel(xpos-x, ypos-y, axisLock)
			case e.mouseMode == "rotate":
				e.RotateSelectedNodeModel(xpos-x, ypos-y, axisLock)
			}
		}
		x, y = xpos, ypos
	}
	c.BindAxisAction(doMouseMove)
	c.BindMouseAction(func() { mouseDown = true }, controller.MouseButtonRight, controller.Press)
	c.BindMouseAction(func() { mouseDown = false }, controller.MouseButtonRight, controller.Release)
	c.BindKeyAction(func() { axisLock = vmath.Vector3{1, 0, 0} }, controller.KeyX, controller.Press)
	c.BindKeyAction(func() { axisLock = vmath.Vector3{0, 1, 0} }, controller.KeyY, controller.Press)
	c.BindKeyAction(func() { axisLock = vmath.Vector3{0, 0, 1} }, controller.KeyZ, controller.Press)

	c.BindKeyAction(func() {
		e.mouseMode = "scale"
		axisLock = vmath.Vector3{1, 1, 1}
	}, controller.KeyT, controller.Press)
	c.BindKeyAction(func() {
		e.mouseMode = "translate"
		axisLock = vmath.Vector3{1, 0, 0}
	}, controller.KeyG, controller.Press)
	c.BindKeyAction(func() {
		e.mouseMode = "rotate"
		axisLock = vmath.Vector3{0, 1, 0}
	}, controller.KeyR, controller.Press)
	return c
}

func (e *Editor) ScaleSelectedNodeModel(x, y float64, axisLock vmath.Vector3) {
	selectedModel, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root)
	if selectedModel != nil {
		selectedModel.Scale = selectedModel.Scale.Add(axisLock.MultiplyScalar(x * mouseSpeed))
	}
	updateMap(e.currentMap.Root)
}

func (e *Editor) MoveSelectedNodeModel(x, y float64, axisLock vmath.Vector3) {
	selectedModel, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root)
	if selectedModel != nil {
		selectedModel.Translation = selectedModel.Translation.Add(axisLock.MultiplyScalar(x * mouseSpeed))
	}
	updateMap(e.currentMap.Root)
}

func (e *Editor) RotateSelectedNodeModel(x, y float64, axisLock vmath.Vector3) {
	selectedModel, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root)
	if selectedModel != nil {
		selectedModel.Orientation = vmath.AngleAxis(x*mouseSpeed, axisLock).Multiply(selectedModel.Orientation).Normalize()
	}
	updateMap(e.currentMap.Root)
}
