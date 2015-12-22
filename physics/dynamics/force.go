package dynamics

import (
	vmath "github.com/walesey/go-engine/vectormath"
)

type ForceStore struct {
	forces map[string]Force
}

type Force interface {
	DoStep(dt float64, phyObj *PhysicsObjectImpl)
}

type LinearForce struct {
	Value vmath.Vector3
}

//acceleration due to gravity
type GravityForce struct {
	Value vmath.Vector3
}

//Force due to friction proportional to velocity
type FrictionForce struct {
	LinearCoefficient, AngularCoefficient float64
}

type TorqueForce struct {
	Value vmath.Vector3
}

// force applied to a fixed local point on the object
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

func (fs *ForceStore) DoStep(dt float64, phyObj *PhysicsObjectImpl) {
	for _, force := range fs.forces {
		force.DoStep(dt, phyObj)
	}
}

func (force *LinearForce) DoStep(dt float64, phyObj *PhysicsObjectImpl) {
	phyObj.Velocity = phyObj.Velocity.Add(
		force.Value.DivideScalar(phyObj.Mass).
			MultiplyScalar(dt))
}

func (force *GravityForce) DoStep(dt float64, phyObj *PhysicsObjectImpl) {
	phyObj.Velocity = phyObj.Velocity.Add(
		force.Value.MultiplyScalar(dt))
}

func (force *FrictionForce) DoStep(dt float64, phyObj *PhysicsObjectImpl) {
	//linear
	forceValue := phyObj.Velocity.MultiplyScalar(-force.LinearCoefficient)
	phyObj.Velocity = phyObj.Velocity.Add(forceValue.MultiplyScalar(dt))

	//angular
	if !vmath.ApproxEqual(phyObj.AngularVelocity.W, 0, 0.00001) {
		angV := phyObj.GetAngularVelocityVector()
		torque := angV.MultiplyScalar(-force.AngularCoefficient)
		newAngV := angV.Add(phyObj.InertiaTensor().Inverse().Transform(torque))
		phyObj.SetAngularVelocityVector(newAngV)
	}
}

func (force *TorqueForce) DoStep(dt float64, phyObj *PhysicsObjectImpl) {
	if !vmath.ApproxEqual(phyObj.AngularVelocity.W, 0, 0.00001) {
		angV := phyObj.GetAngularVelocityVector()
		newAngV := angV.Add(phyObj.InertiaTensor().Inverse().Transform(force.Value))
		phyObj.SetAngularVelocityVector(newAngV)
	}
}

func (force *PointForce) DoStep(dt float64, phyObj *PhysicsObjectImpl) {
	value := phyObj.Orientation.Apply(force.Value)
	if force.Position.ApproxEqual(vmath.Vector3{0, 0, 0}, 0.00001) {
		//linear only
		phyObj.Velocity = phyObj.Velocity.Add(value.DivideScalar(phyObj.Mass).MultiplyScalar(dt))
	} else {
		//update position length
		if !vmath.ApproxEqual(force.Position.LengthSquared(), force.positionLength, 0.001) {
			force.positionLength = force.Position.Length()
		}

		// linear
		position := phyObj.Orientation.Apply(force.Position)
		positionNorm := position.DivideScalar(force.positionLength).MultiplyScalar(-1)
		linearForce := positionNorm.MultiplyScalar(value.Dot(positionNorm))
		phyObj.Velocity = phyObj.Velocity.Add(linearForce.DivideScalar(phyObj.Mass).MultiplyScalar(dt))

		//angular
		angV := phyObj.GetAngularVelocityVector()
		torque := value.Cross(position)
		newAngV := angV.Add(phyObj.InertiaTensor().Inverse().Transform(torque))
		phyObj.SetAngularVelocityVector(newAngV)
	}
}
