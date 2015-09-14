package renderer

import(
	"github.com/walesey/go-engine/vectormath"
)

//A camera is an Entity
type Camera struct {
	renderer Renderer
	Translation, Lookat vectormath.Vector3
	Up vectormath.Vector3
}

func CreateCamera(renderer Renderer) *Camera {
	cam := Camera{
		renderer: renderer,
		Translation: vectormath.Vector3{1,0,0},
		Lookat: vectormath.Vector3{0,0,0},
		Up: vectormath.Vector3{0,1,0},
	}
	return &cam
}

func (c *Camera) SetScale( scale vectormath.Vector3 ) {} //na

func (c *Camera) SetTranslation( translation vectormath.Vector3 ) {
	c.Translation = translation
	c.renderer.Camera( c.Translation, c.Lookat, c.Up )
}

func (c *Camera) SetOrientation( orientation vectormath.Quaternion  ) {
	direction := orientation.Apply( vectormath.Vector3{1,0,0} )
	c.Lookat = c.Translation.Add(direction)
	c.renderer.Camera( c.Translation, c.Lookat, c.Up )
}