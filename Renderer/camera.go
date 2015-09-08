package renderer

import(
)

//A camera is an Entity
type Camera struct {
	renderer *Renderer
	lookat, up vectorMath.Vector3
}

func CreateCamera(renderer *Renderer) Camera {
	camera := Camera{
		renderer: renderer,
		lookat: vectorMath.Vector3{0,0,0},
		up: vectorMath.Vector3{0,1,0},
	}
	return 
}

func (c *Camera) SetScale( scale vectorMath.Vector3 ) {} //na

func (c *Camera) SetTranslation( translation vectorMath.Vector3 ) {
	cameraLocation := vectorMath.Vector3{5*cosine,3*sine,5*sine}
	glRenderer.Camera( translation, c.lookat, c.up )
}

func (c *Camera) SetOrientation( orientation vectorMath.Quaternion  ) {
	//TODO
}