package gjk

import (
	vmath "github.com/walesey/go-engine/vectormath"
)

type SimplexPoint struct {
	mPoint vmath.Vector3
}

type Simplex struct {
	points []SimplexPoint
}

func NewSimplex() *Simplex {
	return &Simplex{make([]SimplexPoint, 0, 4)}
}

func (s *Simplex) Add(point SimplexPoint) {
	if len(s.points) < 4 {
		s.points = append(s.points, point)
	}
}

func (s *Simplex) Remove(index int) {
	s.points = append(s.points[:index], s.points[index+1:]...)
}

func (s *Simplex) Get(index int) SimplexPoint {
	return s.points[index]
}

func (s *Simplex) Len() int {
	return len(s.points)
}
