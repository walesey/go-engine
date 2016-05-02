package editor

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/Invictus321/invictus321-countdown"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/renderer"
)

func (e *Editor) loadMap(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error Reading map file: %v\n", err)
		return
	}

	var mapModel editorModels.MapModel
	err = json.Unmarshal(data, &mapModel)
	if err != nil {
		log.Printf("Error unmarshaling map model: %v\n", err)
		return
	}

	e.currentMap = &mapModel
	e.updateMap(true)
}

func (e *Editor) updateMap(clearMemory bool) {
	e.rootMapNode.RemoveAll(clearMemory)

	cd := countdown.Countdown{}
	cd.Start(countGeometries(e.currentMap.Root))
	e.openProgressBar()
	e.setProgressBar(0)

	updateProgress := func() {
		cd.Count()
		e.setProgressBar(cd.PercentageComplete() / 5)
		//TODO: display cd.SecondsRemaining()
	}

	var updateNode func(srcModel *editorModels.NodeModel, destNode *renderer.Node)
	updateNode = func(srcModel *editorModels.NodeModel, destNode *renderer.Node) {
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
		for _, childModel := range srcModel.Children {
			newNode := renderer.CreateNode()
			destNode.Add(newNode)
			updateNode(childModel, newNode)
		}
	}

	updateNode(e.currentMap.Root, e.rootMapNode)
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
