package renderer

import (
	"sort"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl32/matstack"
	"github.com/walesey/go-engine/util"
)

type SceneGraph struct {
	opaqueNode      *Node
	transparentNode *Node
	txStack         *matstack.TransformStack
	bucketCount     int
	bucket          bucketEntries
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
	sceneGraph.bucketCount = 0
	sceneGraph.buildBuckets(sceneGraph.transparentNode, cameraLocation)
	sort.Sort(sceneGraph.bucket)
	//render buckets
	sceneGraph.opaqueNode.Draw(renderer, mgl32.Ident4())
	for i := 0; i < sceneGraph.bucketCount; i++ {
		entry := sceneGraph.bucket[i]
		entry.parent.DrawChild(renderer, entry.transform, entry.spatial)
	}
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

func (sceneGraph *SceneGraph) buildBuckets(node *Node, cameraLocation mgl32.Vec3) {
	sceneGraph.txStack.Push(node.Transform)
	for _, child := range node.children {
		nextNode, ok := child.(*Node)
		if ok {
			sceneGraph.buildBuckets(nextNode, cameraLocation)
		} else {
			transform := sceneGraph.txStack.Peek()
			entry := bucketEntry{
				spatial:     child,
				parent:      node,
				transform:   transform,
				cameraDelta: util.Vec3LenSq(mgl32.TransformCoordinate(child.Center(), transform).Sub(cameraLocation)),
			}
			if sceneGraph.bucketCount < len(sceneGraph.bucket) {
				sceneGraph.bucket[sceneGraph.bucketCount] = entry
			} else {
				sceneGraph.bucket = append(sceneGraph.bucket, entry)
			}
			sceneGraph.bucketCount = sceneGraph.bucketCount + 1
		}
	}
	sceneGraph.txStack.Pop()
}
