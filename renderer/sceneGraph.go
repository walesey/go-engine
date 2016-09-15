package renderer

import (
	"sort"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl32/matstack"
)

type SceneGraph struct {
	opaqueNode        *Node
	transparentNode   *Node
	txStack           *matstack.TransformStack
	transparentBucket bucketEntries
}

//factory
func CreateSceneGraph() *SceneGraph {
	sceneGraph := &SceneGraph{
		opaqueNode:      CreateNode(),
		transparentNode: CreateNode(),
		txStack:         matstack.NewTransformStack(),
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
	sceneGraph.opaqueNode.Draw(renderer)
	for _, entry := range sceneGraph.transparentBucket {
		renderEntry(entry, renderer)
	}
}

func renderEntry(entry bucketEntry, renderer Renderer) {
	renderer.PushTransform(entry.transform)
	entry.spatial.Draw(renderer)
	renderer.PopTransform()
}

type bucketEntry struct {
	spatial     Spatial
	transform   mgl32.Mat4
	cameraDelta float32
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
	sceneGraph.txStack.Push(node.Transform)
	for _, child := range node.children {
		nextNode, found := child.(*Node)
		if found {
			sceneGraph.buildBuckets(nextNode)
		} else {
			sceneGraph.transparentBucket = append(sceneGraph.transparentBucket, bucketEntry{spatial: child, transform: sceneGraph.txStack.Peek()})
		}
	}
	sceneGraph.txStack.Pop()
}

func (sceneGraph *SceneGraph) sortBuckets(renderer Renderer) {
	for index, entry := range sceneGraph.transparentBucket {
		sceneGraph.transparentBucket[index].cameraDelta = mgl32.TransformCoordinate(entry.spatial.Centre(), entry.transform).Sub(renderer.CameraLocation()).Len()
	}
	sort.Sort(sceneGraph.transparentBucket)
}
