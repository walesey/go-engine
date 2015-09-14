package renderer

import(
	"github.com/Walesey/goEngine/vectorMath"
)

//A camera is an Entity
type Camera struct {
	renderer Renderer
	Translation, Lookat vectorMath.Vector3
	Up vectorMath.Vector3
}

func CreateCamera(renderer Renderer) *Camera {
	cam := Camera{
		renderer: renderer,
		Translation: vectorMath.Vector3{1,0,0},
		Lookat: vectorMath.Vector3{0,0,0},
		Up: vectorMath.Vector3{0,1,0},
	}
	return &cam
}

func (c *Camera) SetScale( scale vectorMath.Vector3 ) {} //na

func (c *Camera) SetTranslation( translation vectorMath.Vector3 ) {
	c.Translation = translation
	c.renderer.Camera( c.Translation, c.Lookat, c.Up )
}

func (c *Camera) SetOrientation( orientation vectorMath.Quaternion  ) {
	direction := orientation.Apply( vectorMath.Vector3{1,0,0} )
	c.Lookat = c.Translation.Add(direction)
	c.renderer.Camera( c.Translation, c.Lookat, c.Up )
}