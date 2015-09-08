package controller

import (
)

type BasicMovementActor interface {
	StartMovingUp()
	StartMovingDown()
	StartMovingRight()
	StartMovingLeft()
	StopMovingUp()
	StopMovingDown()
	StopMovingRight()
	StopMovingLeft()
}