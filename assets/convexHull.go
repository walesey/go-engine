package assets

import (
	"math"

	"github.com/luxengine/gobullet"
	"github.com/walesey/go-engine/physics/collision"
	"github.com/walesey/go-engine/renderer"
	vmath "github.com/walesey/go-engine/vectormath"
)

func CollisionShapeFromGeometry(geometry renderer.Geometry, cullThreshold float64) gobullet.CollisionShape {
	verts := PointsFromGeometry(geometry, cullThreshold)
	convexHull := gobullet.NewConvexHullShape()
	for _, vert := range *verts {
		convexHull.AddVertex(float32(vert.X), float32(vert.Y), float32(vert.Z))
	}
	return gobullet.CollisionShape(convexHull)
}

// Converts a geometry to a ConvexSet Collider
func ConvexSetFromGeometry(geometry renderer.Geometry, cullThreshold float64) collision.Collider {
	return collision.NewConvexSet(PointsFromGeometry(geometry, cullThreshold))
}

// Converts a geometry directly into points (does threshold culling optimisation)
func PointsFromGeometry(geometry renderer.Geometry, cullThreshold float64) *[]vmath.Vector3 {
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
	return &verticies
}

func BoundingBoxFromGeometry(geometry renderer.Geometry) collision.Collider {
	r := RadiusFromGeometry(geometry)
	return BoundingBoxFromRadius(r)
}

func BoundingBoxFromRadius(radius float64) collision.Collider {
	return collision.NewBoundingBox(vmath.Vector3{radius, radius, radius}.MultiplyScalar(2))
}

func RadiusFromGeometry(geometry renderer.Geometry) float64 {
	largest := 0.0
	for i := 0; i < len(geometry.Indicies); i = i + 1 {
		index := geometry.Indicies[i]
		v := vmath.Vector3{
			float64(geometry.Verticies[index*18]),
			float64(geometry.Verticies[index*18+1]),
			float64(geometry.Verticies[index*18+2]),
		}
		if math.Abs(v.X) > largest {
			largest = math.Abs(v.X)
		}
		if math.Abs(v.Y) > largest {
			largest = math.Abs(v.Y)
		}
		if math.Abs(v.Z) > largest {
			largest = math.Abs(v.Z)
		}
	}

	return largest
}
