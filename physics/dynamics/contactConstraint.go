package dynamics

import vmath "github.com/walesey/go-engine/vectormath"

type ContactConstraint struct {
	PenetrationVector            vmath.Vector3
	LocalContact1, LocalContact2 vmath.Vector3
	Object1, Object2             *PhysicsObject
	InContact                    bool
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

	if cc.InContact && cc.Object2.Static {
		//restrict velocity to the tangent plane between the two objects
		normalVelocity1 := norm.MultiplyScalar(cc.Object1.Velocity.Dot(norm))
		cc.Object1.Velocity = cc.Object1.Velocity.Subtract(normalVelocity1)
	}

	//velocities
	angularV1 := cc.Object1.AngularVelocityVector()
	angularV2 := cc.Object2.AngularVelocityVector()
	radialV1 := cc.LocalContact1.Cross(angularV1)
	radialV2 := cc.LocalContact2.Cross(angularV2)
	contactV1 := radialV1.Add(cc.Object1.Velocity)
	contactV2 := radialV2.Add(cc.Object2.Velocity)

	if cc.Object2.Static {
		cc.Object2.Mass = 99999999999999999.9
	}

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
	relativeV := norm.Dot(contactV)
	velocityImpulse := -relativeV

	vel1 := tensor1.Inverse().Transform(cc.LocalContact1.Cross(norm))
	vel1 = vel1.Cross(cc.LocalContact1)
	impulseDenom1 := (1.0 / cc.Object1.Mass) + norm.Dot(vel1)
	vel2 := tensor2.Inverse().Transform(cc.LocalContact2.Cross(norm))
	vel2 = vel2.Cross(cc.LocalContact2)
	impulseDenom2 := (1.0 / cc.Object2.Mass) + norm.Dot(vel2)
	impulseDenom := impulseDenom1 + impulseDenom2
	NormalImpulse := velocityImpulse / impulseDenom

	if NormalImpulse > 0.0 {
		NormalImpulse = 0.0
	}

	impulseVector1 := norm.MultiplyScalar(NormalImpulse)
	impulseVector2 := impulseVector1.MultiplyScalar(-1)
	//project impulse onto the localcontact normal
	localContactNorm1 := cc.LocalContact1.Normalize()
	localContactNorm2 := cc.LocalContact2.Normalize()
	linearImpulse1 := localContactNorm1.MultiplyScalar(impulseVector1.Dot(localContactNorm1))
	linearImpulse2 := localContactNorm2.MultiplyScalar(impulseVector2.Dot(localContactNorm2))

	torqueImpulse1 := impulseVector1.Cross(cc.LocalContact1)
	torqueImpulse2 := impulseVector2.Cross(cc.LocalContact2)

	cc.Object1.Velocity = cc.Object1.Velocity.Add(linearImpulse1.DivideScalar(cc.Object1.Mass))
	newAngularV1 := angularV1.Add(tensor1.Inverse().Transform(torqueImpulse1))
	cc.Object1.SetAngularVelocityVector(newAngularV1)

	if !cc.Object2.Static {
		cc.Object2.Velocity = cc.Object2.Velocity.Add(linearImpulse2.DivideScalar(cc.Object2.Mass))
		newAngularV2 := angularV2.Add(tensor2.Inverse().Transform(torqueImpulse2))
		cc.Object2.SetAngularVelocityVector(newAngularV2)
	}

}
