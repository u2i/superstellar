package simulation

import (
	"superstellar/backend/state"
	"superstellar/backend/types"
)

type AsteroidManager struct {
	asteroidsLeft uint32
}

func NewAsteroidManager() *AsteroidManager {
	return &AsteroidManager{asteroidsLeft: 1}
}

func (manager *AsteroidManager) update(space *state.Space) {
	if manager.asteroidsLeft > 0 {
		manager.asteroidsLeft--

		asteroid := state.NewAsteroid(1000, types.NewPoint(1000.0, 1000.0), types.NewVector(0.0, 0.0))
		space.AddAsteroid(asteroid)
	}
}

