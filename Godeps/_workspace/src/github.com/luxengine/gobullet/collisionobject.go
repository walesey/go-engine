package gobullet

type CollisionObject interface {
	GetActivationState() ActivationTag
	SetActivationState(ActivationTag)
	SetDeactivationTime(time float32)
	GetDeactivationTime() float32
	ForceActivationState(ActivationTag)
	Activate(bool)
	IsActive() bool
}

type ActivationTag int

const (
	ACTIVE_TAG           ActivationTag = 1
	ISLAND_SLEEPING                    = 2
	WANTS_DEACTIVATION                 = 3
	DISABLE_DEACTIVATION               = 4
	DISABLE_SIMULATION                 = 5
)
