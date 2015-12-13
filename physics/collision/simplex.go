package collision

import (
	vmath "github.com/walesey/go-engine/vectormath"
)

type SimplexPoint struct {
	mPoint     vmath.Vector3
	point      vmath.Vector3
	otherPoint vmath.Vector3
}

type SimplexFace struct {
	index      int
	p1, p2, p3 int
}

type Simplex struct {
	points []SimplexPoint
	faces  []SimplexFace
}

func NewSimplex() *Simplex {
	return &Simplex{
		points: make([]SimplexPoint, 0, 4),
		faces:  make([]SimplexFace, 0, 4),
	}
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

func (s *Simplex) Add(point SimplexPoint) int {
	s.points = append(s.points, point)
	return len(s.points) - 1
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

func (s *Simplex) containsPoint(point vmath.Vector3, epsilon float64) bool {
	for _, p := range s.points {
		if p.mPoint.ApproxEqual(point, epsilon) {
			return true
		}
	}
	return false
}

/////////
//Faces

func (s *Simplex) ClosestFace() (SimplexFace, float64) {
	min := 999999999999999999999999999999999.9
	var closestFace SimplexFace
	for _, face := range s.faces {
		dist := vmath.PointToPlaneDist(
			s.Get(face.p1).mPoint,
			s.Get(face.p2).mPoint,
			s.Get(face.p3).mPoint,
			vmath.Vector3{0, 0, 0})
		if dist < min {
			min = dist
			closestFace = face
		}
	}
	return closestFace, min
}

func (s *Simplex) AddPointToFace(point SimplexPoint, faceIndex int) {
	face := s.GetFace(faceIndex)
	s.RemoveFace(faceIndex)
	pIndex := s.Add(point)
	s.AddFace(SimplexFace{p1: face.p1, p2: face.p2, p3: pIndex})
	s.AddFace(SimplexFace{p1: face.p2, p2: face.p3, p3: pIndex})
	s.AddFace(SimplexFace{p1: face.p3, p2: face.p1, p3: pIndex})
}

func (s *Simplex) FaceNormal(index int) vmath.Vector3 {
	f := s.GetFace(index)
	p1 := s.Get(f.p1).mPoint
	p2 := s.Get(f.p2).mPoint
	p3 := s.Get(f.p3).mPoint
	norm := p2.Subtract(p1).Cross(p3.Subtract(p1))
	if p1.Dot(norm) < 0 {
		return norm.MultiplyScalar(-1)
	}
	return norm
}

func (s *Simplex) GetFace(index int) SimplexFace {
	return s.faces[index]
}

func (s *Simplex) RemoveFace(index int) {
	s.faces = append(s.faces[:index], s.faces[index+1:]...)
	//update all indicies
	for i, _ := range s.faces {
		s.faces[i].index = i
	}
}

func (s *Simplex) ClearFaces() {
	s.faces = s.faces[:0]
}

func (s *Simplex) AddFace(face SimplexFace) {
	s.faces = append(s.faces, face)
	index := len(s.faces) - 1
	s.faces[index].index = index
}
