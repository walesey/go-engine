package editor

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/controller"
)

const mouseSpeed = 0.01

func NewEditorController(e *Editor) controller.Controller {
	c := controller.CreateController()
	var mouseDown bool
	var x, y float32 = 0.0, 0.0
	axisLock := mgl32.Vec3{1, 1, 1}

	doMouseMove := func(xpos, ypos float32) {
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
	c.BindKeyAction(func() { axisLock = mgl32.Vec3{1, 0, 0} }, controller.KeyX, controller.Press)
	c.BindKeyAction(func() { axisLock = mgl32.Vec3{0, 1, 0} }, controller.KeyY, controller.Press)
	c.BindKeyAction(func() { axisLock = mgl32.Vec3{0, 0, 1} }, controller.KeyZ, controller.Press)

	c.BindKeyAction(func() {
		e.mouseMode = "scale"
		axisLock = mgl32.Vec3{1, 1, 1}
	}, controller.KeyT, controller.Press)
	c.BindKeyAction(func() {
		e.mouseMode = "translate"
		axisLock = mgl32.Vec3{1, 0, 0}
	}, controller.KeyG, controller.Press)
	c.BindKeyAction(func() {
		e.mouseMode = "rotate"
		axisLock = mgl32.Vec3{0, 1, 0}
	}, controller.KeyR, controller.Press)
	return c
}

func (e *Editor) ScaleSelectedNodeModel(x, y float32, axisLock mgl32.Vec3) {
	selectedModel, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root)
	if selectedModel != nil {
		selectedModel.Scale = selectedModel.Scale.Add(axisLock.Mul(x * mouseSpeed))
	}
	updateMap(e.currentMap.Root)
}

func (e *Editor) MoveSelectedNodeModel(x, y float32, axisLock mgl32.Vec3) {
	selectedModel, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root)
	if selectedModel != nil {
		selectedModel.Translation = selectedModel.Translation.Add(axisLock.Mul(x * mouseSpeed))
	}
	updateMap(e.currentMap.Root)
}

func (e *Editor) RotateSelectedNodeModel(x, y float32, axisLock mgl32.Vec3) {
	selectedModel, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root)
	if selectedModel != nil {
		selectedModel.Orientation = mgl32.QuatRotate(x*mouseSpeed, axisLock).Mul(selectedModel.Orientation).Normalize()
	}
	updateMap(e.currentMap.Root)
}
