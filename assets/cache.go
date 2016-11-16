package assets

import (
	"image"
	"sync"

	"github.com/walesey/go-engine/renderer"
)

type AssetCache struct {
	geometries map[string]*renderer.Geometry
	materials  map[string]*renderer.Material
	images     map[string]image.Image
	mutex      *sync.Mutex
}

var globalCache *AssetCache

func init() {
	globalCache = NewAssetCache()
}

func (ac *AssetCache) ImportObj(path string) (geometry *renderer.Geometry, material *renderer.Material, err error) {
	var okGeom, okMat bool
	geometry, okGeom = ac.geometries[path]
	material, okMat = ac.materials[path]
	if !okGeom || !okMat {
		geometry, material, err = ImportObj(path)
		if err != nil {
			return
		}
		ac.mutex.Lock()
		ac.geometries[path] = geometry
		ac.materials[path] = material
		ac.mutex.Unlock()
	}
	return
}

func (ac *AssetCache) ImportImage(path string) (image.Image, error) {
	image, ok := ac.images[path]
	if !ok {
		var err error
		image, err = ImportImage(path)
		if err != nil {
			return image, err
		}
		ac.mutex.Lock()
		ac.images[path] = image
		ac.mutex.Unlock()
	}
	return image, nil
}

func ImportImageCached(path string) (image.Image, error) {
	return globalCache.ImportImage(path)
}

func ImportObjCached(path string) (geometry *renderer.Geometry, material *renderer.Material, err error) {
	return globalCache.ImportObj(path)
}

func NewAssetCache() *AssetCache {
	return &AssetCache{
		geometries: make(map[string]*renderer.Geometry),
		materials:  make(map[string]*renderer.Material),
		images:     make(map[string]image.Image),
		mutex:      &sync.Mutex{},
	}
}
