package dynamics

import vmath "github.com/walesey/go-engine/vectormath"

type ContactConstraint struct {
	PenetrationVector            vmath.Vector3
	LocalContact1, LocalContact2 vmath.Vector3
	Object1, Object2             *PhysicsObject
}

func (cc *ContactConstraint) Solve() {

	//collision normal
	var norm vmath.Vector3
	if cc.PenetrationVector.LengthSquared() > 0 {
		norm = cc.PenetrationVector.Normalize()
	} else if !cc.Object2.Position.ApproxEqual(cc.Object1.Position, 0.00001) {
		norm = cc.Object2.Position.Subtract(cc.Object1.Position).Normalize()
	} else {
		norm = vmath.Vector3{1, 0, 0}
	}

	//velocities
	angularV1 := cc.Object1.AngularVelocityVector()
	angularV2 := cc.Object2.AngularVelocityVector()
	radialV1 := cc.LocalContact1.Cross(angularV1)
	radialV2 := cc.LocalContact2.Cross(angularV2)
	contactV1 := radialV1.Add(cc.Object1.Velocity)
	contactV2 := radialV2.Add(cc.Object2.Velocity)

	mR1 := 0.4 * cc.Object1.Mass * cc.Object1.Radius
	mR2 := 0.4 * cc.Object2.Mass * cc.Object2.Radius
	tensor1 := vmath.Matrix3{
		mR1, 0.0, 0.0,
		0.0, mR1, 0.0,
		0.0, 0.0, mR1,
	}
	tensor2 := vmath.Matrix3{
		mR2, 0.0, 0.0,
		0.0, mR2, 0.0,
		0.0, 0.0, mR2,
	}

	contactV := contactV1.Subtract(contactV2)
	norm = contactV.Normalize()
	relativeV := norm.Dot(contactV)
	velocityImpulse := -relativeV

	vel1 := tensor1.Inverse().Transform(cc.LocalContact1.Cross(norm))
	vel1 = vel1.Cross(cc.LocalContact1)
	impulseDenom1 := (1.0 / cc.Object1.Mass) + norm.Dot(vel1)
	vel2 := tensor2.Inverse().Transform(cc.LocalContact2.Cross(norm))
	vel2 = vel2.Cross(cc.LocalContact2)
	impulseDenom2 := (1.0 / cc.Object2.Mass) + norm.Dot(vel2)
	impulseDenom := impulseDenom1
	if !cc.Object2.Static {
		impulseDenom = impulseDenom + impulseDenom2
	}
	NormalImpulse := velocityImpulse / impulseDenom
	NormalImpulse = NormalImpulse + (cc.PenetrationVector.Length() / impulseDenom)

	if NormalImpulse > 0.0 {
		NormalImpulse = 0.0
	}

	impulseVector1 := norm.MultiplyScalar(NormalImpulse)
	impulseVector2 := impulseVector1.MultiplyScalar(-1)
	torqueImpulse1 := impulseVector1.Cross(cc.LocalContact1)
	torqueImpulse2 := impulseVector2.Cross(cc.LocalContact2)

	cc.Object1.Velocity = cc.Object1.Velocity.Add(impulseVector1.DivideScalar(cc.Object1.Mass))
	newAngularV1 := angularV1.Add(tensor1.Inverse().Transform(torqueImpulse1))
	cc.Object1.SetAngularVelocityVector(newAngularV1)

	if !cc.Object2.Static {
		cc.Object2.Velocity = cc.Object2.Velocity.Add(impulseVector2.DivideScalar(cc.Object2.Mass))
		newAngularV2 := angularV2.Add(tensor2.Inverse().Transform(torqueImpulse2))
		cc.Object2.SetAngularVelocityVector(newAngularV2)
	}
}
