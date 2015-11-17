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

func (force LinearForce) DoStep(dt float64, phyObj *PhysicsObject) {
	phyObj.Velocity = phyObj.Velocity.Add(
		force.Value.DivideScalar(phyObj.Mass).
			MultiplyScalar(dt))
}

func (force GravityForce) DoStep(dt float64, phyObj *PhysicsObject) {
	phyObj.Velocity = phyObj.Velocity.Add(
		force.Value.MultiplyScalar(dt))
}

func (force TorqueForce) DoStep(dt float64, phyObj *PhysicsObject) {
	//TODO:
}

func (force PointForce) DoStep(dt float64, phyObj *PhysicsObject) {
	phyObj.Velocity = phyObj.Velocity.Add(
		force.Value.DivideScalar(phyObj.Mass).
			MultiplyScalar(dt))
}
