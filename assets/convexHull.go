package assets

import (
	"github.com/walesey/go-engine/physics"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
	"math"
)

// Converts a geometry directly into triangles (does no optimisation and assumes geometry is already convex)
func ConvexHullFromGeometry(geometry renderer.Geometry) physics.Collider {
	triangles := make([]physics.Triangle, 0, len(geometry.Indicies)/3)
	for i := 0; i < len(geometry.Indicies); i = i + 3 {
		i1, i2, i3 := geometry.Indicies[i], geometry.Indicies[i+1], geometry.Indicies[i+2]
		v11, v12, v13 := geometry.Verticies[i1*18], geometry.Verticies[i1*18+1], geometry.Verticies[i1*18+2]
		v21, v22, v23 := geometry.Verticies[i2*18], geometry.Verticies[i2*18+1], geometry.Verticies[i2*18+2]
		v31, v32, v33 := geometry.Verticies[i3*18], geometry.Verticies[i3*18+1], geometry.Verticies[i3*18+2]
		triangles = append(triangles, physics.Triangle{
			Point1: vmath.Vector3{float64(v11), float64(v12), float64(v13)},
			Point2: vmath.Vector3{float64(v21), float64(v22), float64(v23)},
			Point3: vmath.Vector3{float64(v31), float64(v32), float64(v33)},
		})
	}
	return physics.NewConvexHull(triangles)
}

// Converts a geometry directly into points (does no optimisation and assumes geometry is already convex)
func ConvexSetFromGeometry(geometry renderer.Geometry) physics.Collider {
	verticies := make([]vmath.Vector3, 0, len(geometry.Indicies))
	for i := 0; i < len(geometry.Indicies); i = i + 1 {
		index := geometry.Indicies[i]
		v := vmath.Vector3{
			float64(geometry.Verticies[index*18]),
			float64(geometry.Verticies[index*18+1]),
			float64(geometry.Verticies[index*18+2]),
		}
		verticies = append(verticies, v)
	}
	return physics.NewConvexSet(verticies)
}

func BoundingBoxFromGeometry(geometry renderer.Geometry) physics.Collider {
	largest := vmath.Vector3{0, 0, 0}
	largestSide := 0.0
	for i := 0; i < len(geometry.Indicies); i = i + 1 {
		index := geometry.Indicies[i]
		v := vmath.Vector3{
			float64(geometry.Verticies[index*18]),
			float64(geometry.Verticies[index*18+1]),
			float64(geometry.Verticies[index*18+2]),
		}
		if math.Abs(v.X) > largestSide {
			largest = v
			largestSide = v.X
		}
		if math.Abs(v.Y) > largestSide {
			largest = v
			largestSide = v.Y
		}
		if math.Abs(v.Z) > largestSide {
			largest = v
			largestSide = v.Z
		}
	}
	return physics.NewBoundingBox(largest)
}
