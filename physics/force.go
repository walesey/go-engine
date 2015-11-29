package physics

import (
	vmath "github.com/walesey/go-engine/vectormath"
)

type ForceStore struct {
	forces map[string]Force
}

type Force interface {
	DoStep(dt float64, phyObj *PhysicsObject)
}

type LinearForce struct {
	Value vmath.Vector3
}

//acceleration due to gravity
type GravityForce struct {
	Value vmath.Vector3
}

type TorqueForce struct {
	Value vmath.Quaternion
}

type PointForce struct {
	Value, Position vmath.Vector3
	positionLength  float64
}

func NewForceStore() *ForceStore {
	return &ForceStore{make(map[string]Force)}
}

func (fs *ForceStore) GetForce(name string) Force {
	return fs.forces[name]
}

func (fs *ForceStore) AddForce(name string, force Force) *ForceStore {
	fs.forces[name] = force
	return fs
}

func (fs *ForceStore) RemoveForce(name string) *ForceStore {
	delete(fs.forces, name)
	return fs
}

func (fs *ForceStore) DoStep(dt float64, phyObj *PhysicsObject) {
	for _, force := range fs.forces {
		force.DoStep(dt, phyObj)
	}
}

func (force *LinearForce) DoStep(dt float64, phyObj *PhysicsObject) {
	phyObj.Velocity = phyObj.Velocity.Add(
		force.Value.DivideScalar(phyObj.Mass).
			MultiplyScalar(dt))
}

func (force *GravityForce) DoStep(dt float64, phyObj *PhysicsObject) {
	phyObj.Velocity = phyObj.Velocity.Add(
		force.Value.MultiplyScalar(dt))
}

func (force *TorqueForce) DoStep(dt float64, phyObj *PhysicsObject) {
	//TODO:
}

func (force *PointForce) DoStep(dt float64, phyObj *PhysicsObject) {
	//update position length
	if !vmath.ApproxEqual(force.Position.LengthSquared(), force.positionLength, 0.001) {
		force.positionLength = force.Position.Length()
	}
	// linear
	positionNorm := force.Position.DivideScalar(force.positionLength)
	linearForce := positionNorm.MultiplyScalar(force.Value.Dot(positionNorm))
	phyObj.Velocity = phyObj.Velocity.Add(linearForce.DivideScalar(phyObj.Mass).MultiplyScalar(dt))
	//angular
	tangentForce := force.Value.Subtract(linearForce)
	axis := tangentForce.Cross(force.Position)
	if axis.LengthSquared() < 0.000001 {
		axis = vmath.Vector3{1, 0, 0}
	} else {
		axis = axis.Normalize()
	}
	magnitude := tangentForce.Length()
	angularVelocity := vmath.Vector3{phyObj.AngularVelocity.X, phyObj.AngularVelocity.Y, phyObj.AngularVelocity.Z}
	if !vmath.ApproxEqual(angularVelocity.LengthSquared(), 1, 0.00001) {
		if vmath.ApproxEqual(angularVelocity.LengthSquared(), 0, 0.00001) {
			angularVelocity = axis
		} else {
			angularVelocity = angularVelocity.Normalize()
		}
	}
	angularVelocity.MultiplyScalar(phyObj.AngularVelocity.W)
	angularVelocity = vmath.AngleAxis(magnitude, axis).Apply(angularVelocity)
	phyObj.AngularVelocity.W = angularVelocity.Length()
	//normalize
	if !vmath.ApproxEqual(phyObj.AngularVelocity.W, 0, 0.00001) {
		angularVelocity = angularVelocity.DivideScalar(phyObj.AngularVelocity.W)
	}
	phyObj.AngularVelocity.X = angularVelocity.X
	phyObj.AngularVelocity.Y = angularVelocity.Y
	phyObj.AngularVelocity.Z = angularVelocity.Z
}
