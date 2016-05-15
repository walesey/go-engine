package editor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Invictus321/invictus321-countdown"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/renderer"
)

func (e *Editor) saveMap(filepath string) {
	data, err := json.Marshal(e.currentMap)
	if err != nil {
		log.Printf("Error Marshaling map file: %v\n", err)
		return
	}

	ioutil.WriteFile(filepath, data, os.ModePerm)
}

func (e *Editor) loadMap(path string) {
	e.currentMap = assets.LoadMap(path)
	e.updateMap()
	e.overviewMenu.updateTree(e.currentMap)
}

func LoadMapToNode(model *editorModels.NodeModel, node *renderer.Node, updateProgress func()) {
	var updateNode func(srcModel *editorModels.NodeModel, destNode *renderer.Node)
	updateNode = func(srcModel *editorModels.NodeModel, destNode *renderer.Node) {
		srcModel.SetNode(destNode)
		if srcModel.Geometry != nil {
			geometry, err := assets.ImportObjCached(*srcModel.Geometry)
			if err == nil {
				destNode.Add(geometry)
			}
			updateProgress()
		}
		destNode.SetScale(srcModel.Scale)
		destNode.SetTranslation(srcModel.Translation)
		destNode.SetOrientation(srcModel.Orientation)
		if srcModel.Reference != nil {
			if refModel, _ := findNodeById(*srcModel.Reference, model); refModel != nil {
				for _, childModel := range refModel.Children {
					refNode := childModel.GetNode()
					if refNode != nil {
						destNode.Add(refNode)
					}
				}
			}
		}
		for _, childModel := range srcModel.Children {
			newNode := renderer.CreateNode()
			destNode.Add(newNode)
			updateNode(childModel, newNode)
		}
	}

	updateNode(model, node)
}

func UpdateMapNode(model *editorModels.NodeModel) {
	node := model.GetNode()
	if node != nil {
		node.SetScale(model.Scale)
		node.SetTranslation(model.Translation)
		node.SetOrientation(model.Orientation)
	}
	for _, childModel := range model.Children {
		UpdateMapNode(childModel)
	}
}

func (e *Editor) updateMap() {
	e.rootMapNode.RemoveAll(true)

	cd := countdown.Countdown{}
	cd.Start(countGeometries(e.currentMap.Root))
	e.openProgressBar()
	e.setProgressBar(0)
	e.setProgressTime("Loading Map...")

	LoadMapToNode(e.currentMap.Root, e.rootMapNode, func() {
		cd.Count()
		e.setProgressBar(cd.PercentageComplete() / 5)
		e.setProgressTime(fmt.Sprintf("Loading Map... %v seconds remaining", cd.SecondsRemaining()))
	})

	e.closeProgressBar()
}

func countGeometries(nodeModel *editorModels.NodeModel) int {
	count := 0
	if nodeModel.Geometry != nil {
		count++
	}
	for _, child := range nodeModel.Children {
		count += countGeometries(child)
	}
	return count
}
