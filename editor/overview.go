package editor

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/walesey/go-engine/assets"
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

	planetOpenImg, _ := assets.DecodeImage(bytes.NewBuffer(assets.Base64ToBytes(PlanetOpenData)))
	e.uiAssets.AddImage("planetOpen", planetOpenImg)

	planetClosedImg, _ := assets.DecodeImage(bytes.NewBuffer(assets.Base64ToBytes(PlanetClosedData)))
	e.uiAssets.AddImage("planetClosed", planetClosedImg)

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
	if node := findNodeById(e.overviewMenu.selectedNodeId, e.currentMap.Root); node != nil {
		node.Geometry = &filePath
		e.updateMap(true)
	}
}

func (e *Editor) createNewGroup(name string) {
	if node := findNodeById(e.overviewMenu.selectedNodeId, e.currentMap.Root); node != nil {
		if check := findNodeById(name, e.currentMap.Root); check != nil {
			return // id already taken
		}
		node.Children = append(node.Children, editorModels.NewNodeModel(name))
	}
}

func findNodeById(nodeId string, model *editorModels.NodeModel) *editorModels.NodeModel {
	if model.Id == nodeId {
		return model
	}
	for _, childModel := range model.Children {
		result := findNodeById(nodeId, childModel)
		if result != nil {
			return result
		}
	}
	return nil
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
