package assets

import (
	"log"

	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/renderer"
)

type geomImport struct {
	geometry *renderer.Geometry
	callback func(geometry *renderer.Geometry)
}

type mapImport struct {
	node     *renderer.Node
	model    *editorModels.NodeModel
	callback func(node *renderer.Node, model *editorModels.NodeModel)
}

// The loader allows asyncronous loading of obj and map files
type Loader struct {
	geoms chan geomImport
	maps  chan mapImport
}

func NewLoader() *Loader {
	return &Loader{
		geoms: make(chan geomImport, 256),
		maps:  make(chan mapImport, 256),
	}
}

func (loader *Loader) Update(dt float64) {
	for {
		select {
		case g := <-loader.geoms:
			g.callback(g.geometry)
		case m := <-loader.maps:
			m.callback(m.node, m.model)
		default:
			return
		}
	}
}

func (loader *Loader) LoadMap(path string, callback func(node *renderer.Node, model *editorModels.NodeModel)) {
	go func() {
		srcModel := LoadMap(path)
		destNode := renderer.CreateNode()
		loadedModel := LoadMapToNode(srcModel.Root, destNode)
		loader.maps <- mapImport{
			node:     destNode,
			model:    loadedModel,
			callback: callback,
		}
	}()
}

func (loader *Loader) LoadObj(path string, callback func(geometry *renderer.Geometry)) {
	go func() {
		loadedGeometry, err := ImportObjCached(path)
		if err != nil {
			log.Println("Error Loading Obj: ", err)
		} else {
			loader.geoms <- geomImport{
				geometry: loadedGeometry,
				callback: callback,
			}
		}
	}()
}
