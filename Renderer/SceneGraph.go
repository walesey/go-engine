package renderer

//Interfaces
type SceneGraph interface {
	RenderScene()
	Add( spatial Spatial )
}

//DefaultSceneGraph
type DefaultSceneGraph struct {
	SceneRenderer Renderer
	root *Node
}

//factory
func CreateSceneGraph( renderer Renderer ) SceneGraph{
	return &DefaultSceneGraph{ renderer, CreateNode() }
}

func (sceneGraph *DefaultSceneGraph) RenderScene() {
	sceneGraph.root.load( sceneGraph.SceneRenderer )
	sceneGraph.root.draw( sceneGraph.SceneRenderer )
}

func (sceneGraph *DefaultSceneGraph) Add( spatial Spatial ) {
	sceneGraph.root.Add(spatial)
}
