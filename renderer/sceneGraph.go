package renderer

import (
	"sort"
)

type SceneGraph struct {
	backgroundNode, transparentNode, orthoNode Node
	matStack                                   Stack
	transparentBucket                          bucketEntries
}

//factory
func CreateSceneGraph() SceneGraph {
	sceneGraph := SceneGraph{
		backgroundNode:  CreateNode(),
		transparentNode: CreateNode(),
		orthoNode:       CreateNode(),
		matStack:        CreateStack(),
	}
	return sceneGraph
}

func (sceneGraph *SceneGraph) Add(spatial Spatial) {
	sceneGraph.transparentNode.Add(spatial)
}

func (sceneGraph *SceneGraph) AddBackGround(spatial Spatial) {
	sceneGraph.backgroundNode.Add(spatial)
}

func (sceneGraph *SceneGraph) AddOrtho(spatial Spatial) {
	sceneGraph.orthoNode.Add(spatial)
}

func (sceneGraph *SceneGraph) RenderScene(renderer Renderer) {
	//setup buckets
	sceneGraph.transparentBucket = sceneGraph.transparentBucket[:0]
	sceneGraph.buildBuckets(&sceneGraph.transparentNode)
	sceneGraph.sortBuckets(renderer)
	//render buckets
	//renderer.EnableDepthTest(true)
	sceneGraph.backgroundNode.Draw(renderer)
	//renderer.EnableDepthTest(false)
	for _, entry := range sceneGraph.transparentBucket {
		renderEntry(entry, renderer)
	}
	sceneGraph.orthoNode.Draw(renderer)
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
	tx := sceneGraph.matStack.ApplyAll()
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
func (s *Stack) ApplyAll() Transform {
	result := CreateTransform()
	for i := 0; i < s.size; i++ {
		result.ApplyTransform(s.Get(i).(Transform))
	}
	return result
}
