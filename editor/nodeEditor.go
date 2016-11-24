package editor

import (
	"fmt"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/ui"
)

func (e *Editor) closeNodeEditor() {
	if e.fileBrowserOpen {
		e.gameEngine.RemoveSpatial(e.fileBrowser.window, false)
		e.fileBrowserOpen = false
	}
}

func (e *Editor) openNodeEditor(node *editorModels.NodeModel, callback func()) {
	window, container, _ := e.defaultWindow()
	window.SetTranslation(mgl32.Vec3{100, 100, 1})
	window.SetScale(mgl32.Vec3{500, 0, 1})
	uiController := ui.NewUiController(window).(glfwController.Controller)

	e.uiAssets.AddCallback("nodeEditorOk", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			node.Id = window.TextElementById("name").GetText()
			node.Classes = []string{}
			for i := 1; i <= 3; i++ {
				className := window.TextElementById(fmt.Sprintf("class%v", i)).GetText()
				if len(className) > 0 {
					node.Classes = append(node.Classes, className)
				}
			}
			e.gameEngine.RemoveSpatial(window, true)
			e.controllerManager.RemoveController(uiController)
			callback()
		}
	})

	e.uiAssets.AddCallback("nodeEditorCancel", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			e.gameEngine.RemoveSpatial(window, true)
			e.controllerManager.RemoveController(uiController)
		}
	})

	e.controllerManager.AddController(uiController)
	window.Tabs, _ = ui.LoadHTML(container, strings.NewReader(nodeEditMenuHtml), strings.NewReader(globalCss), e.uiAssets)

	e.gameEngine.AddOrtho(window)

	window.TextElementById("name").SetText(node.Id)
	for i, class := range node.Classes {
		if i < 3 {
			window.TextElementById(fmt.Sprintf("class%v", i+1)).SetText(class)
		}
	}
	window.Render()
}
