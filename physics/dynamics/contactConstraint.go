package dynamics

import vmath "github.com/walesey/go-engine/vectormath"

type ContactConstraint struct {
	BodyIndex1, BodyIndex2       int
	Body1, Body2                 *PhysicsObjectImpl
	LocalContact1, LocalContact2 vmath.Vector3
	Normal                       vmath.Vector3
}
