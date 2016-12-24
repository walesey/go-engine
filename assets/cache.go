package assets

import (
	"fmt"
	"image"
	"sync"

	"github.com/walesey/go-engine/renderer"
)

type AssetCache struct {
	geometries  map[string]*renderer.Geometry
	materials   map[string]*renderer.Material
	images      map[string]image.Image
	fileMutexes map[string]*sync.Mutex
	mutex       *sync.Mutex
}

var globalCache *AssetCache

func init() {
	globalCache = NewAssetCache()
}

func (ac *AssetCache) ImportObj(path string) (geometry *renderer.Geometry, material *renderer.Material, err error) {
	ac.lockFilepath(path)
	var okGeom, okMat bool
	geometry, okGeom = ac.geometries[path]
	material, okMat = ac.materials[path]
	if !okGeom && !okMat {
		fmt.Println(path)
		geometry, material, err = ImportObj(path)
		ac.mutex.Lock()
		ac.geometries[path] = geometry
		ac.materials[path] = material
		ac.mutex.Unlock()
	}
	ac.unlockFilepath(path)
	return
}

func (ac *AssetCache) ImportImage(path string) (img image.Image, err error) {
	ac.lockFilepath(path)
	var ok bool
	img, ok = ac.images[path]
	if !ok {
		img, err = ImportImage(path)
		ac.mutex.Lock()
		ac.images[path] = img
		ac.mutex.Unlock()
	}
	ac.unlockFilepath(path)
	return
}

func (ac *AssetCache) fileMutex(path string) *sync.Mutex {
	if _, ok := ac.fileMutexes[path]; !ok {
		ac.fileMutexes[path] = &sync.Mutex{}
	}
	return ac.fileMutexes[path]
}

func (ac *AssetCache) lockFilepath(path string) {
	ac.fileMutex(path).Lock()
}

func (ac *AssetCache) unlockFilepath(path string) {
	ac.fileMutex(path).Unlock()
}

func ImportImageCached(path string) (image.Image, error) {
	return globalCache.ImportImage(path)
}

func ImportObjCached(path string) (geometry *renderer.Geometry, material *renderer.Material, err error) {
	return globalCache.ImportObj(path)
}

func NewAssetCache() *AssetCache {
	return &AssetCache{
		geometries:  make(map[string]*renderer.Geometry),
		materials:   make(map[string]*renderer.Material),
		images:      make(map[string]image.Image),
		fileMutexes: make(map[string]*sync.Mutex),
		mutex:       &sync.Mutex{},
	}
}
