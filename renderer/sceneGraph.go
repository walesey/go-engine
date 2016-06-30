package renderer

import (
	"sort"

	"github.com/walesey/go-engine/util"
)

type SceneGraph struct {
	opaqueNode        *Node
	transparentNode   *Node
	matStack          util.Stack
	transparentBucket bucketEntries
}

//factory
func CreateSceneGraph() *SceneGraph {
	sceneGraph := &SceneGraph{
		opaqueNode:      CreateNode(),
		transparentNode: CreateNode(),
		matStack:        util.CreateStack(),
	}
	return sceneGraph
}

func (sceneGraph *SceneGraph) AddTransparent(spatial Spatial) {
	sceneGraph.transparentNode.Add(spatial)
}

func (sceneGraph *SceneGraph) RemoveTransparent(spatial Spatial, destroy bool) {
	sceneGraph.transparentNode.Remove(spatial, destroy)
}

func (sceneGraph *SceneGraph) Add(spatial Spatial) {
	sceneGraph.opaqueNode.Add(spatial)
}

func (sceneGraph *SceneGraph) Remove(spatial Spatial, destroy bool) {
	sceneGraph.opaqueNode.Remove(spatial, destroy)
}

func (sceneGraph *SceneGraph) RenderScene(renderer Renderer) {
	//setup buckets
	sceneGraph.transparentBucket = sceneGraph.transparentBucket[:0]
	sceneGraph.buildBuckets(sceneGraph.transparentNode)
	sceneGraph.sortBuckets(renderer)
	//render buckets
	renderer.EnableDepthMask(true)
	sceneGraph.opaqueNode.Draw(renderer)
	renderer.EnableDepthMask(false)
	for _, entry := range sceneGraph.transparentBucket {
		renderEntry(entry, renderer)
	}
	renderer.EnableDepthMask(true)
}

func renderEntry(entry bucketEntry, renderer Renderer) {
	renderer.PushTransform()
	renderer.ApplyTransform(entry.transform)
	entry.spatial.Draw(renderer)
	renderer.PopTransform()
}

type bucketEntry struct {
	spatial     Spatial
	transform   Transform
	cameraDelta float64
}

type bucketEntries []bucketEntry

func (slice bucketEntries) Len() int {
	return len(slice)
}

func (slice bucketEntries) Less(i, j int) bool {
	return slice[i].cameraDelta > slice[j].cameraDelta
}

func (slice bucketEntries) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (sceneGraph *SceneGraph) buildBuckets(node *Node) {
	sceneGraph.matStack.Push(node.Transform)
	tx := ApplyAll(sceneGraph.matStack)
	for _, child := range node.children {
		nextNode, found := child.(*Node)
		if found {
			sceneGraph.buildBuckets(nextNode)
		} else {
			sceneGraph.transparentBucket = append(sceneGraph.transparentBucket, bucketEntry{spatial: child, transform: tx})
		}
	}
	sceneGraph.matStack.Pop()
}

func (sceneGraph *SceneGraph) sortBuckets(renderer Renderer) {
	for index, entry := range sceneGraph.transparentBucket {
		sceneGraph.transparentBucket[index].cameraDelta = entry.transform.TransformCoordinate(entry.spatial.Centre()).Subtract(renderer.CameraLocation()).LengthSquared()
	}
	sort.Sort(sceneGraph.transparentBucket)
}

//used to combine transformations
func ApplyAll(s util.Stack) Transform {
	result := CreateTransform()
	for i := 0; i < s.Len(); i++ {
		result.ApplyTransform(s.Get(i).(Transform))
	}
	return result
}
