package gobulletutils

import (
	"bufio"
	"fmt"
	bullet "github.com/luxengine/gobullet"
	"os"
	"strings"
)

func ShapeFromWavefront(path string) (bullet.GImpactMeshShape, error) {
	file, err := os.Open(path)
	if err != nil {
		return bullet.GImpactMeshShape{}, err
	}
	defer file.Close()

	tmpposition := make([][3]float32, 0)

	positionsindices := make([][3]uint16, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "v ") {
			var p1, p2, p3 float32
			fmt.Sscanf(line, "v %f %f %f", &p1, &p2, &p3)
			tmpposition = append(tmpposition, [3]float32{p1, p2, p3})
		} else if strings.HasPrefix(line, "f ") {
			var p1, p2, p3, vt1, vt2, vt3, vn1, vn2, vn3 uint16
			fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d", &p1, &vt1, &vn1, &p2, &vt2, &vn2, &p3, &vt3, &vn3)
			positionsindices = append(positionsindices, [3]uint16{p1, p2, p3})
		}
	}
	mesh := bullet.NewTriangleMesh(false, false)
	for _, x := range positionsindices {
		mesh.AddTriangle(&tmpposition[x[0]-1], &tmpposition[x[1]-1], &tmpposition[x[2]-1], true)
	}
	shape := bullet.NewGImpactMeshShapeFromTriangleMesh(mesh)
	return shape, nil
}

func ShapeFromWavefrontStatic(path string) (bullet.TriangleMeshShape, error) {
	file, err := os.Open(path)
	if err != nil {
		return bullet.TriangleMeshShape{}, err
	}
	defer file.Close()

	tmpposition := make([][3]float32, 0)

	positionsindices := make([][3]uint16, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "v ") {
			var p1, p2, p3 float32
			fmt.Sscanf(line, "v %f %f %f", &p1, &p2, &p3)
			tmpposition = append(tmpposition, [3]float32{p1, p2, p3})
		} else if strings.HasPrefix(line, "f ") {
			var p1, p2, p3, vt1, vt2, vt3, vn1, vn2, vn3 uint16
			fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d", &p1, &vt1, &vn1, &p2, &vt2, &vn2, &p3, &vt3, &vn3)
			positionsindices = append(positionsindices, [3]uint16{p1, p2, p3})
		}
	}
	mesh := bullet.NewTriangleMesh(false, false)
	for _, x := range positionsindices {
		mesh.AddTriangle(&tmpposition[x[0]-1], &tmpposition[x[1]-1], &tmpposition[x[2]-1], true)
	}
	shape := mesh.NewTriangleMeshShape(true, true)
	return shape, nil
}
