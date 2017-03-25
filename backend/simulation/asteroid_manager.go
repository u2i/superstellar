package simulation

import (
	"superstellar/backend/state"
	"superstellar/backend/types"
)

type AsteroidManager struct {
	space         *state.Space
	asteroidsLeft uint32
}

func NewAsteroidManager(space *state.Space) *AsteroidManager {
	return &AsteroidManager{
		space:         space,
		asteroidsLeft: 1,
	}
}

func (manager *AsteroidManager) updateAsteroids() {
	if manager.asteroidsLeft > 0 {
		manager.asteroidsLeft--

		asteroid := state.NewAsteroid(1000, types.NewPoint(1000.0, 1000.0), types.NewVector(0.0, 0.0))
		manager.space.AddAsteroid(asteroid)
	}
}
