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
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/renderer"
)

type MapLoadUpdate struct {
	node        *renderer.Node
	geomsLoaded int
}

func (e *Editor) saveMap(filepath string) {
	data, err := json.MarshalIndent(e.currentMap, "", "  ")
	if err != nil {
		log.Printf("Error Marshaling map file: %v\n", err)
		return
	}

	ioutil.WriteFile(filepath, data, os.ModePerm)
}

func (e *Editor) loadMap(path string) {
	e.currentMap = assets.LoadMap(path)
	e.refreshMap()
	e.overviewMenu.updateTree(e.currentMap)
}

func (e *Editor) refreshMap() {
	e.rootMapNode.RemoveAll(true)

	cd := countdown.Countdown{}
	cd.Start(countGeometries(e.currentMap.Root))
	e.openProgressBar()
	e.setProgressBar(0)
	e.setProgressTime("Loading Map...")

	mapLoadChan := loadMapToNode(e.currentMap.Root)

	var loader engine.Updatable
	loader = engine.UpdatableFunc(func(dt float64) {
		select {
		case mapLoadUpdate := <-mapLoadChan:
			if mapLoadUpdate.node == nil {
				cd.Count()
				e.setProgressBar(cd.PercentageComplete() / 5)
				e.setProgressTime(fmt.Sprintf("Loading Map... %v seconds remaining", cd.SecondsRemaining()))
			} else {
				e.rootMapNode.Add(mapLoadUpdate.node)
				e.gameEngine.RemoveUpdatable(loader)
				e.closeProgressBar()
			}
		default:
		}
	})
	e.gameEngine.AddUpdatable(loader)
}

func loadMapToNode(model *editorModels.NodeModel) chan MapLoadUpdate {

	out := make(chan MapLoadUpdate)
	geomsLoaded := 0

	var updateNode func(srcModel *editorModels.NodeModel, destNode *renderer.Node)
	updateNode = func(srcModel *editorModels.NodeModel, destNode *renderer.Node) {
		srcModel.SetNode(destNode)
		if srcModel.Geometry != nil {
			geometry, material, err := assets.ImportObjCached(*srcModel.Geometry)
			if err == nil {
				destNode.Material = material
				destNode.Add(geometry)
			}
			geomsLoaded++
			out <- MapLoadUpdate{nil, geomsLoaded}
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
			newNode := renderer.NewNode()
			destNode.Add(newNode)
			updateNode(childModel, newNode)
		}
	}

	node := renderer.NewNode()
	go func() {
		updateNode(model, node)
		out <- MapLoadUpdate{node, geomsLoaded}
	}()
	return out
}

func updateMap(model *editorModels.NodeModel) {
	node := model.GetNode()
	if node != nil {
		node.SetScale(model.Scale)
		node.SetTranslation(model.Translation)
		node.SetOrientation(model.Orientation)
	}
	for _, childModel := range model.Children {
		updateMap(childModel)
	}
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
