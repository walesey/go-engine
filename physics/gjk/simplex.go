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

func (s *Simplex) ContainsOrigin(direction *vmath.Vector3) bool {
	a := s.GetLast().mPoint
	ao := a.MultiplyScalar(-1)
	size := s.Len()
	switch {
	//four points case
	case size == 4:
		b := s.Get(0).mPoint
		c := s.Get(1).mPoint
		d := s.Get(2).mPoint

		ab := b.Subtract(a)
		ac := c.Subtract(a)
		ad := d.Subtract(a)

		//region abc?
		normal := ab.Cross(ac)
		if normal.Dot(ad) > 0 {
			normal = normal.MultiplyScalar(-1)
		}
		if normal.Dot(ao) > 0 {
			s.Remove(2)
			direction.Set(normal)
		} else {
			//region acd?
			normal = ac.Cross(ad)
			if normal.Dot(ab) > 0 {
				normal = normal.MultiplyScalar(-1)
			}
			if normal.Dot(ao) > 0 {
				s.Remove(0)
				direction.Set(normal)
			} else {
				//region abd?
				normal = ab.Cross(ad)
				if normal.Dot(ac) > 0 {
					normal = normal.MultiplyScalar(-1)
				}
				if normal.Dot(ao) > 0 {
					s.Remove(1)
					direction.Set(normal)
				} else {
					//inside the simplex
					return true
				}
			}
		}
	case size == 3: //triangle case
		b := s.Get(0).mPoint
		c := s.Get(1).mPoint
		ab := b.Subtract(a)
		ac := c.Subtract(a)
		planeNormal := ab.Cross(ac)
		if planeNormal.Dot(ao) < 0 {
			direction.Set(planeNormal.MultiplyScalar(-1))
		} else {
			direction.Set(planeNormal)
		}
	default: //line case
		b := s.Get(0).mPoint
		ab := b.Subtract(a)
		abPerp := ab.Cross(ao).Cross(ab)
		direction.Set(abPerp)
		if direction.LengthSquared() <= 0.001 {
			if ab.Dot(ao) >= 0 && ab.LengthSquared() >= ao.LengthSquared() {
				//origin is on the line
				return true
			}
		}
	}
	return false
}

func (s *Simplex) Add(point SimplexPoint) {
	if len(s.points) < 4 {
		s.points = append(s.points, point)
	}
}

func (s *Simplex) Remove(index int) {
	s.points = append(s.points[:index], s.points[index+1:]...)
}

func (s *Simplex) Clear() {
	s.points = s.points[:0]
}

func (s *Simplex) Get(index int) SimplexPoint {
	return s.points[index]
}

func (s *Simplex) GetLast() SimplexPoint {
	return s.points[len(s.points)-1]
}

func (s *Simplex) Len() int {
	return len(s.points)
}
