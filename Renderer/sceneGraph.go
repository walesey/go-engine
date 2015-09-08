package renderer

import(
	"sort"
	"github.com/Walesey/goEngine/vectorMath"
)

//Sorts the scene into buckets and renders nodes in the correct order - opaque, transparent, ortho
type SceneGraph struct {
	root Node
	matStack Stack
	opaqueBucket bucketEntries
	transparentBucket bucketEntries
	orthoBucket bucketEntries
}

//factory
func CreateSceneGraph() SceneGraph{
	sceneGraph := SceneGraph{ root: CreateNode(), matStack: CreateStack() }
	return sceneGraph
}

type bucketEntry struct {
	spatial Spatial
	transform Transform
	cameraDelta float64
}

type bucketEntries []bucketEntry

func (slice bucketEntries) Len() int {
    return len(slice)
}

func (slice bucketEntries) Less(i, j int) bool {
    return slice[i].cameraDelta > slice[j].cameraDelta;
}

func (slice bucketEntries) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

func (sceneGraph *SceneGraph) Add( spatial Spatial ) {
	sceneGraph.root.Add(spatial)
}

func (sceneGraph *SceneGraph) RenderScene( renderer Renderer ) {
	sceneGraph.opaqueBucket = make(bucketEntries,0,0)
	sceneGraph.transparentBucket = make(bucketEntries,0,0)
	sceneGraph.orthoBucket = make(bucketEntries,0,0)
	sceneGraph.buildBuckets(&sceneGraph.root)
	sceneGraph.sortBuckets(renderer)
	sceneGraph.drawBuckets(renderer)
}

func (sceneGraph *SceneGraph) buildBuckets( node *Node ) {
	if node.BucketType == BUCKET_OPAQUE {
		sceneGraph.opaqueBucket = append( sceneGraph.opaqueBucket, bucketEntry{spatial: node, transform: sceneGraph.matStack.ApplyAll()} )
	} else if node.BucketType == BUCKET_ORTHO {
		sceneGraph.orthoBucket = append( sceneGraph.orthoBucket, bucketEntry{spatial: node, transform: sceneGraph.matStack.ApplyAll()} )
	} else if node.BucketType == BUCKET_TRANSPARENT {
		sceneGraph.matStack.Push( node.Transform )
		tx := sceneGraph.matStack.ApplyAll()
	    for _,child := range node.children {
	    	nextNode, found := child.(*Node)
	    	if found {
		    	sceneGraph.buildBuckets(nextNode)
	    	} else {
		        sceneGraph.transparentBucket = append( sceneGraph.transparentBucket, bucketEntry{spatial: child, transform: tx} )
		    }
	    }
		sceneGraph.matStack.Pop()
	}
}

func (sceneGraph *SceneGraph) sortBuckets( renderer Renderer ) {
	for index,entry := range sceneGraph.transparentBucket {
		sceneGraph.transparentBucket[index].cameraDelta = entry.transform.TransformCoordinate(vectorMath.Vector3{0,0,0}).Subtract(renderer.CameraLocation()).LengthSquared()
	}
	sort.Sort(sceneGraph.transparentBucket)
}

func (sceneGraph *SceneGraph) drawBuckets( renderer Renderer ) {
	for _,entry := range sceneGraph.opaqueBucket {
		renderEntry(entry, renderer)
	}
	for _,entry := range sceneGraph.transparentBucket {
		renderEntry(entry, renderer)
	}
	for _,entry := range sceneGraph.orthoBucket {
		renderEntry(entry, renderer)
	}
}

func renderEntry( entry bucketEntry, renderer Renderer ) {
	renderer.PushTransform()
	renderer.ApplyTransform(entry.transform)
	entry.spatial.Draw( renderer )
	renderer.PopTransform()
}

//used to combine transformations
func (s *Stack) ApplyAll() Transform {
	result := CreateTransform()
	for i:=0 ; i<s.size ; i++ {
		result.ApplyTransform(s.Get(i).(Transform))
	}
	return result
}
