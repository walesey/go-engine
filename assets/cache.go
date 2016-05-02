package assets

import "github.com/walesey/go-engine/renderer"

type AssetCache struct {
	geometries map[string]*renderer.Geometry
}

var globalCache *AssetCache

func init() {
	globalCache = NewAssetCache()
}

func (ac *AssetCache) ImportObj(path string) (*renderer.Geometry, error) {
	_, ok := ac.geometries[path]
	if !ok {
		geometry, err := ImportObj(path)
		if err != nil {
			return geometry, err
		}
		ac.geometries[path] = geometry
	}
	return ac.geometries[path], nil
}

func ImportObjCached(path string) (*renderer.Geometry, error) {
	return globalCache.ImportObj(path)
}

func NewAssetCache() *AssetCache {
	return &AssetCache{
		geometries: make(map[string]*renderer.Geometry),
	}
}
