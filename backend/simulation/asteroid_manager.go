package simulation

import (
	"fmt"
	"superstellar/backend/constants"
	"superstellar/backend/state"
	"superstellar/backend/types"
)

type AsteroidManager struct {
	space  *state.Space
	lastId uint32
}

func NewAsteroidManager(space *state.Space) *AsteroidManager {
	return &AsteroidManager{
		space:  space,
		lastId: 1000,
	}
}

func (manager *AsteroidManager) updateAsteroids() {
	manager.spawnNewAsteroids()
	manager.removeObsoleteAsteroids()
}

func (manager *AsteroidManager) spawnNewAsteroids() {
	if len(manager.space.Asteroids) < 1 {

		asteroid := state.NewAsteroid(manager.lastId, types.NewPoint(1000.0, 1000.0), types.NewVector(5.0, 0.0))
		manager.space.AddAsteroid(asteroid)
		manager.lastId++
	}
}

func (manager *AsteroidManager) removeObsoleteAsteroids() {
	for _, asteroid := range manager.space.Asteroids {
		if asteroid.Position().Length() > constants.WorldRadius*2 {
			manager.space.RemoveAsteroid(asteroid.Id())
			fmt.Println(asteroid.Position().Length())
		}
	}
}
