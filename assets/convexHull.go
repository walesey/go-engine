package assets

import (
	"math"

	"github.com/walesey/go-engine/physics"
	"github.com/walesey/go-engine/physics/gjk"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

// Converts a geometry directly into points (does threshold culling optimisation)
func ConvexSetFromGeometry(geometry renderer.Geometry, cullThreshold float64) physics.Collider {
	verticies := make([]vmath.Vector3, 0, len(geometry.Indicies))
	for i := 0; i < len(geometry.Indicies); i = i + 1 {
		index := geometry.Indicies[i]
		v := vmath.Vector3{
			float64(geometry.Verticies[index*18]),
			float64(geometry.Verticies[index*18+1]),
			float64(geometry.Verticies[index*18+2]),
		}
		//do culling
		include := true
		for _, vert := range verticies {
			if vert.Subtract(v).LengthSquared() < cullThreshold*cullThreshold {
				include = false
				break
			}
		}
		if include {
			verticies = append(verticies, v)
		}
	}
	return gjk.NewConvexSet(verticies)
}

func BoundingBoxFromGeometry(geometry renderer.Geometry) physics.Collider {
	largest := vmath.Vector3{0, 0, 0}
	for i := 0; i < len(geometry.Indicies); i = i + 1 {
		index := geometry.Indicies[i]
		v := vmath.Vector3{
			float64(geometry.Verticies[index*18]),
			float64(geometry.Verticies[index*18+1]),
			float64(geometry.Verticies[index*18+2]),
		}
		if math.Abs(v.X) > largest.X {
			largest.X = math.Abs(v.X)
		}
		if math.Abs(v.Y) > largest.Y {
			largest.Y = math.Abs(v.Y)
		}
		if math.Abs(v.Z) > largest.Z {
			largest.Z = math.Abs(v.Z)
		}
	}

	return physics.NewBoundingBox(largest.MultiplyScalar(2))
}
