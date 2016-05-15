package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/renderer"
)

func LoadMap(path string) *editorModels.MapModel {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error Reading map file: %v\n", err)
		return nil
	}

	var mapModel editorModels.MapModel
	err = json.Unmarshal(data, &mapModel)
	if err != nil {
		log.Printf("Error unmarshaling map model: %v\n", err)
		return nil
	}

	return &mapModel
}

func LoadMapToNode(srcModel *editorModels.NodeModel, destNode *renderer.Node) *editorModels.NodeModel {
	copy := srcModel.Copy(func(name string) string { return name })
	loadMapRecursive(srcModel, copy, destNode)
	return copy
}

func loadMapRecursive(srcModel, model *editorModels.NodeModel, destNode *renderer.Node) {
	srcModel.SetNode(destNode)
	if srcModel.Geometry != nil {
		geometry, err := ImportObjCached(*srcModel.Geometry)
		if err == nil {
			destNode.Add(geometry)
		}
	}
	destNode.SetScale(srcModel.Scale)
	destNode.SetTranslation(srcModel.Translation)
	destNode.SetOrientation(srcModel.Orientation)
	if srcModel.Reference != nil {
		if refModel := findNodeById(*srcModel.Reference, model); refModel != nil {
			for _, childModel := range refModel.Children {
				srcModel.Children = append(srcModel.Children, childModel.Copy(func(name string) string {
					return fmt.Sprintf("%v::%v", srcModel.Reference, name)
				}))
			}
		}
	}
	for _, childModel := range srcModel.Children {
		newNode := renderer.CreateNode()
		destNode.Add(newNode)
		loadMapRecursive(childModel, model, newNode)
	}
}

func findNodeById(nodeId string, model *editorModels.NodeModel) *editorModels.NodeModel {
	if model.Id == nodeId {
		return model
	}
	for _, childModel := range model.Children {
		node := findNodeById(nodeId, childModel)
		if node != nil {
			return node
		}
	}
	return nil
}
