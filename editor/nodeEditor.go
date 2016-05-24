package editor

import (
	"fmt"
	"strings"

	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

func (e *Editor) closeNodeEditor() {
	if e.fileBrowserOpen {
		e.gameEngine.RemoveOrtho(e.fileBrowser.window, false)
		e.fileBrowserOpen = false
	}
}

func (e *Editor) openNodeEditor(node *editorModels.NodeModel, callback func()) {
	window, container, _ := e.defaultWindow()
	window.SetTranslation(vmath.Vector3{100, 100, 1})
	window.SetScale(vmath.Vector3{500, 0, 1})
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
			e.gameEngine.RemoveOrtho(window, true)
			e.controllerManager.RemoveController(uiController)
			callback()
		}
	})

	e.uiAssets.AddCallback("nodeEditorCancel", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			e.gameEngine.RemoveOrtho(window, true)
			e.controllerManager.RemoveController(uiController)
		}
	})

	e.controllerManager.AddController(uiController)
	ui.LoadHTML(container, strings.NewReader(nodeEditMenuHtml), strings.NewReader(globalCss), e.uiAssets)

	e.gameEngine.AddOrtho(window)

	window.TextElementById("name").SetText(node.Id)
	for i, class := range node.Classes {
		if i < 3 {
			window.TextElementById(fmt.Sprintf("class%v", i+1)).SetText(class)
		}
	}
	window.Render()
}
