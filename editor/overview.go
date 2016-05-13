package editor

import (
	"fmt"
	"strings"

	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

var uniqueIdCounter int

type Overview struct {
	window         *ui.Window
	assets         ui.HtmlAssets
	selectedNodeId string
	closedNodes    map[string]bool
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
			e.updateMap(true)
		}
	})

	e.uiAssets.AddCallback("deleteGroup", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			e.deleteGroup()
			e.overviewMenu.updateTree(e.currentMap)
			e.updateMap(true)
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

	e.controllerManager.AddController(ui.NewUiController(window))
	ui.LoadHTML(container, window, strings.NewReader(overviewMenuHtml), strings.NewReader(globalCss), e.uiAssets)

	e.gameEngine.AddOrtho(window)
	e.overviewMenu = &Overview{
		window:      window,
		assets:      e.uiAssets,
		closedNodes: make(map[string]bool),
	}
	e.overviewMenu.updateTree(e.currentMap)
}

func (e *Editor) setGeametry(filePath string) {
	if node, _ := e.overviewMenu.getSelectedNode(e.currentMap.Root); node != nil {
		node.Geometry = &filePath
		e.updateMap(true)
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

func (e *Editor) resetGroup() {
	selectedNode, ok := e.nodeIndex[e.overviewMenu.selectedNodeId]
	if node, _ := findNodeById(e.overviewMenu.selectedNodeId, e.currentMap.Root); node != nil && ok {
		node.Scale = vmath.Vector3{1, 1, 1}
		node.Translation = vmath.Vector3{}
		node.Orientation = vmath.IdentityQuaternion()
		selectedNode.SetScale(node.Scale)
		selectedNode.SetTranslation(node.Translation)
		selectedNode.SetOrientation(node.Orientation)
	}
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
					o.closedNodes[model.Id] = !o.closedNodes[model.Id]
				}
				o.selectedNodeId = model.Id
				o.updateTree(mapModel)
			}
		})
		isClosed := o.closedNodes[model.Id]

		iconImg := "planetOpen"
		if isClosed {
			iconImg = "planetClosed"
		}

		html := fmt.Sprintf("<div onclick=%v><img src=%v></img><p>%v</p>", onclickName, iconImg, model.Id)
		for _, class := range model.Classes {
			html = fmt.Sprintf("%v :: <p>%v</p>", html, class)
		}
		html = fmt.Sprintf("%v</div>", html)
		css := `
		p { font-size: 12px; width: 80%; padding: 0 0 0 5px; }
		img { width: 16px; height: 16px; }
		div { margin: 0 0 5px 0; }
		`
		if model.Id == o.selectedNodeId {
			css = fmt.Sprintf("%v div { background-color: #ff5 }", css)
		}
		if isClosed {
			css = fmt.Sprintf("%v p { color: #999 }", css)
		}
		ui.LoadHTML(container, o.window, strings.NewReader(html), strings.NewReader(css), o.assets)

		if !isClosed {
			for _, childModel := range model.Children {
				nodeContainer := ui.NewContainer()
				container.AddChildren(nodeContainer)
				nodeContainer.SetPadding(ui.Margin{0, 0, 0, 20})
				updateNode(childModel, nodeContainer)
			}
		}
	}

	elem := o.window.ElementById("overviewTree")
	container, ok := elem.(*ui.Container)
	if ok {
		container.RemoveAllChildren()
		updateNode(mapModel.Root, container)
	}
}
