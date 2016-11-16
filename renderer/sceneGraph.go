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
		opaqueNode:      NewNode(),
		transparentNode: NewNode(),
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

func (sceneGraph *SceneGraph) RenderScene(renderer Renderer, cameraLocation mgl32.Vec3) {
	//setup buckets
	sceneGraph.transparentBucket = sceneGraph.transparentBucket[:0]
	sceneGraph.buildBuckets(sceneGraph.transparentNode)
	sceneGraph.sortBuckets(renderer, cameraLocation)
	//render buckets
	sceneGraph.opaqueNode.Draw(renderer, mgl32.Ident4())
	for _, entry := range sceneGraph.transparentBucket {
		renderEntry(entry, renderer)
	}
}

func renderEntry(entry bucketEntry, renderer Renderer) {
	entry.parent.DrawChild(renderer, entry.transform, entry.spatial)
}

type bucketEntry struct {
	spatial     Spatial
	parent      *Node
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
		nextNode, ok := child.(*Node)
		if ok {
			sceneGraph.buildBuckets(nextNode)
		} else {
			entry := bucketEntry{
				spatial:   child,
				parent:    node,
				transform: sceneGraph.txStack.Peek(),
			}
			sceneGraph.transparentBucket = append(sceneGraph.transparentBucket, entry)
		}
	}
	sceneGraph.txStack.Pop()
}

func (sceneGraph *SceneGraph) sortBuckets(renderer Renderer, cameraLocation mgl32.Vec3) {
	for index, entry := range sceneGraph.transparentBucket {
		sceneGraph.transparentBucket[index].cameraDelta = mgl32.TransformCoordinate(entry.spatial.Centre(), entry.transform).Sub(cameraLocation).Len()
	}
	sort.Sort(sceneGraph.transparentBucket)
}
