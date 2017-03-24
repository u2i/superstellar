package state

import (
	"superstellar/backend/types"
)

type Asteroid struct {
	ObjectState
}

func NewAsteroid(id uint32, initialPosition *types.Point, initialVelocity *types.Vector) *Asteroid {
	objectState := NewObjectState(id, initialPosition, initialVelocity)

	return &Asteroid{
		ObjectState:               *objectState,
	}
}