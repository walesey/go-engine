package editor

import (
	"fmt"
	"strings"

	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

var uniqueIdCounter int

type Overview struct {
	window         *ui.Window
	assets         ui.HtmlAssets
	selectedNodeId string
	openNodes      map[string]bool
}

func (e *Editor) initOverviewMenu() {

	e.uiAssets.AddCallback("import", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			e.openFileBrowser("import", func(filePath string) {
				e.setGeametry(filePath)
				e.overviewMenu.updateTree(e.currentMap)
				e.closeFileBrowser()
			}, ".obj")
		}
	})

	e.uiAssets.AddCallback("newGroup", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			uniqueIdCounter++
			e.createNewGroup(fmt.Sprintf("group_%v", uniqueIdCounter))
			e.overviewMenu.updateTree(e.currentMap)
		}
	})

	e.uiAssets.AddCallback("copyGroup", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			e.copyNewGroup()
			e.overviewMenu.updateTree(e.currentMap)
			e.refreshMap()
		}
	})

	e.uiAssets.AddCallback("referenceGroup", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			e.referenceGroup()
			e.overviewMenu.updateTree(e.currentMap)
			e.refreshMap()
		}
	})

	e.uiAssets.AddCallback("deleteGroup", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			e.deleteGroup()
			e.overviewMenu.updateTree(e.currentMap)
			e.refreshMap()
		}
	})

	e.uiAssets.AddCallback("editGroup", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			e.editGroup(func() {
				e.overviewMenu.updateTree(e.currentMap)
			})
		}
	})

	e.uiAssets.AddCallback("scale", func(element ui.Element, args ...interface{}) {
		e.mouseMode = "scale"
	})

	e.uiAssets.AddCallback("translate", func(element ui.Element, args ...interface{}) {
		e.mouseMode = "translate"
	})

	e.uiAssets.AddCallback("rotate", func(element ui.Element, args ...interface{}) {
		e.mouseMode = "rotate"
	})

	e.uiAssets.AddCallback("reset", func(element ui.Element, args ...interface{}) {
		e.resetGroup()
	})

	window, container, _ := e.defaultWindow()
	window.SetTranslation(vmath.Vector3{10, 10, 1})
	window.SetScale(vmath.Vector3{500, 0, 1})

	e.controllerManager.AddController(ui.NewUiController(window).(glfwController.Controller))
	ui.LoadHTML(container, window, strings.NewReader(overviewMenuHtml), strings.NewReader(globalCss), e.uiAssets)

	e.gameEngine.AddOrtho(window)
	e.overviewMenu = &Overview{
		window:    window,
		assets:    e.uiAssets,
		openNodes: make(map[string]bool),
	}
	e.overviewMenu.updateTree(e.currentMap)
}

func (e *Editor) setGeametry(filePath string) {
	if node, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root); node != nil {
		node.Geometry = &filePath
		e.refreshMap()
	}
}

func (e *Editor) createNewGroup(name string) {
	if node, _ := findNodeById(e.overviewMenu.selectedNodeId, e.currentMap.Root); node != nil {
		if check, _ := findNodeById(name, e.currentMap.Root); check != nil {
			return // id already taken
		}
		node.Children = append(node.Children, editorModels.NewNodeModel(name))
	}
}

func (e *Editor) copyNewGroup() {
	if node, parent := findNodeById(e.overviewMenu.selectedNodeId, e.currentMap.Root); node != nil && parent != nil {
		uniqueIdCounter++
		parent.Children = append(parent.Children, node.Copy(func(name string) string {
			return fmt.Sprintf("%v_copy_%v", name, uniqueIdCounter)
		}))
	}
}

func (e *Editor) referenceGroup() {
	if node, parent := findNodeById(e.overviewMenu.selectedNodeId, e.currentMap.Root); node != nil && parent != nil {
		uniqueIdCounter++
		name := fmt.Sprintf("%v_link_%v", node.Id, uniqueIdCounter)
		for check, _ := findNodeById(name, e.currentMap.Root); check != nil; uniqueIdCounter++ {
			name = fmt.Sprintf("%v_link_%v", node.Id, uniqueIdCounter)
		}
		newModel := editorModels.NewNodeModel(name)
		newModel.Reference = &node.Id
		parent.Children = append(parent.Children, newModel)
	}
}

func (e *Editor) deleteGroup() {
	if node, parent := findNodeById(e.overviewMenu.selectedNodeId, e.currentMap.Root); node != nil && parent != nil {
		for i, childNode := range parent.Children {
			if childNode == node {
				parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
				break
			}
		}
	}
}

func (e *Editor) editGroup(cb func()) {
	if node, _ := findNodeById(e.overviewMenu.selectedNodeId, e.currentMap.Root); node != nil {
		e.openNodeEditor(node, cb)
	}
}

func (e *Editor) resetGroup() {
	if node, _ := findNodeById(e.overviewMenu.selectedNodeId, e.currentMap.Root); node != nil {
		selectedNode := node.GetNode()
		if selectedNode != nil {
			selectedNode.SetScale(node.Scale)
			selectedNode.SetTranslation(node.Translation)
			selectedNode.SetOrientation(node.Orientation)
		}
		node.Scale = vmath.Vector3{1, 1, 1}
		node.Translation = vmath.Vector3{}
		node.Orientation = vmath.IdentityQuaternion()
	}
}

func searchNodesByReference(referenceId string, model *editorModels.NodeModel) []*editorModels.NodeModel {
	results := []*editorModels.NodeModel{}
	if *model.Reference == referenceId {
		results = append(results, model)
	}
	for _, childModel := range model.Children {
		results = append(results, searchNodesByReference(referenceId, childModel)...)
	}
	return results
}

func findNodeById(nodeId string, model *editorModels.NodeModel) (node, parent *editorModels.NodeModel) {
	if model.Id == nodeId {
		node = model
		return
	}
	for _, childModel := range model.Children {
		node, parent = findNodeById(nodeId, childModel)
		if node != nil {
			if parent == nil {
				parent = model
			}
			return
		}
	}
	return
}

func (o *Overview) getSelectedNode(model *editorModels.NodeModel) (node, parent *editorModels.NodeModel) {
	return findNodeById(o.selectedNodeId, model)
}

func (o *Overview) updateTree(mapModel *editorModels.MapModel) {
	var updateNode func(model *editorModels.NodeModel, container *ui.Container)
	updateNode = func(model *editorModels.NodeModel, container *ui.Container) {

		//TODO: add native html support for selecting deselecting from a list
		onclickName := fmt.Sprintf("onClickMapNode_%v", model.Id)
		o.assets.AddCallback(onclickName, func(element ui.Element, args ...interface{}) {
			if len(args) >= 2 && !args[1].(bool) { // not on release
				if o.selectedNodeId == model.Id {
					o.openNodes[model.Id] = !o.openNodes[model.Id]
				}
				o.selectedNodeId = model.Id
				o.updateTree(mapModel)
			}
		})

		isOpen := o.openNodes[model.Id]

		iconImg := "planetClosed"
		if isOpen {
			iconImg = "planetOpen"
		}
		if model.Reference != nil {
			iconImg = "reference"
		}

		html := fmt.Sprintf("<div onclick=%v><img src=%v></img><p>%v", onclickName, iconImg, model.Id)
		for _, class := range model.Classes {
			html = fmt.Sprintf("%v :: %v", html, class)
		}
		html = fmt.Sprintf("%v</p></div>", html)
		css := `
		p { font-size: 12px; width: 80%; padding: 0 0 0 5px; }
		img { width: 16px; height: 16px; }
		div { margin: 0 0 5px 0; }
		`
		if model.Id == o.selectedNodeId {
			css = fmt.Sprintf("%v div { background-color: #ff5 }", css)
		}
		if !isOpen {
			css = fmt.Sprintf("%v p { color: #999 }", css)
		}
		ui.LoadHTML(container, o.window, strings.NewReader(html), strings.NewReader(css), o.assets)

		if isOpen {
			for _, childModel := range model.Children {
				nodeContainer := ui.NewContainer()
				container.AddChildren(nodeContainer)
				nodeContainer.SetPadding(ui.Margin{0, 0, 0, 20})
				updateNode(childModel, nodeContainer)
			}
		}
	}

	container, ok := o.window.ElementById("overviewTree").(*ui.Container)
	if ok {
		container.RemoveAllChildren()
		updateNode(mapModel.Root, container)
	}
}
