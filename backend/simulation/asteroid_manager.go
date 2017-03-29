package simulation

import (
	"math"
	"math/rand"
	"superstellar/backend/constants"
	"superstellar/backend/state"
	"superstellar/backend/types"
	"superstellar/backend/utils"
	"time"
)

type AsteroidManager struct {
	space     *state.Space
	idManager *utils.IdManager
	rand      *rand.Rand
}

func NewAsteroidManager(space *state.Space, idManager *utils.IdManager) *AsteroidManager {
	return &AsteroidManager{
		space:     space,
		idManager: idManager,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (manager *AsteroidManager) updateAsteroids() {
	manager.spawnNewAsteroids()
	manager.removeObsoleteAsteroids()
}

func (manager *AsteroidManager) spawnNewAsteroids() {
	if len(manager.space.Asteroids) < constants.AsteroidCountLimit {
		manager.space.AddAsteroid(manager.newAsteroid())
	}
}

func (manager *AsteroidManager) newAsteroid() *state.Asteroid {
	circleAngle := manager.rand.Float64() * 2 * math.Pi
	circlePosition := types.NewPointFromPolar(circleAngle, constants.AsteroidSpawnRadius)

	directionRange := (manager.rand.Float64() - 0.5) * 0.25 * math.Pi
	directionAngle := circleAngle - math.Pi + directionRange

	direction := types.NewVector(constants.AsteroidVelocity, 0.0).Rotate(directionAngle)

	return state.NewAsteroid(manager.idManager.NextAsteroidsId(), circlePosition, direction)
}

func (manager *AsteroidManager) removeObsoleteAsteroids() {
	for _, asteroid := range manager.space.Asteroids {
		if asteroid.Position().Length() > constants.AsteroidRemoveRadius {
			manager.space.RemoveAsteroid(asteroid.Id())
		}
	}
}
