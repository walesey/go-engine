package assets

import (
	"image"
	"sync"

	"github.com/walesey/go-engine/renderer"
)

type AssetCache struct {
	geometries map[string]*renderer.Geometry
	images     map[string]image.Image
	mutex      *sync.Mutex
}

var globalCache *AssetCache

func init() {
	globalCache = NewAssetCache()
}

func (ac *AssetCache) ImportObj(path string) (*renderer.Geometry, error) {
	geometry, ok := ac.geometries[path]
	if !ok {
		var err error
		geometry, err = ImportObj(path)
		if err != nil {
			return geometry, err
		}
		ac.mutex.Lock()
		ac.geometries[path] = geometry
		ac.mutex.Unlock()
	}
	return geometry, nil
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

func ImportObjCached(path string) (*renderer.Geometry, error) {
	return globalCache.ImportObj(path)
}

func NewAssetCache() *AssetCache {
	return &AssetCache{
		geometries: make(map[string]*renderer.Geometry),
		images:     make(map[string]image.Image),
		mutex:      &sync.Mutex{},
	}
}
