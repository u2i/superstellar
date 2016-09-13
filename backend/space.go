package backend

import (
	"math"
	"math/rand"
)

// Space struct holds entire game state.
type Space struct {
	Spaceships map[string]*Spaceship `json:"spaceships"`
}

// NewSpace initializes new Space.
func NewSpace() *Space {
	return &Space{Spaceships: make(map[string]*Spaceship)}
}

// AddSpaceship adds new spaceship to the space.
func (space *Space) AddSpaceship(clientID string, spaceship *Spaceship) {
	space.Spaceships[clientID] = spaceship
}

// RemoveSpaceship removes spaceship from the space.
func (space *Space) RemoveSpaceship(clientID string) {
	delete(space.Spaceships, clientID)
}

func (space *Space) randomUpdate() {
	for _, e := range space.Spaceships {
		if rand.Float64() < 0.05 {
			e.InputThrust = !e.InputThrust
		}
		if rand.Float64() < 0.07 {
			e.InputDirection = Direction(rand.Int() % 3)
		}
	}
}

func (space *Space) updatePhysics() {
	for _, spaceship := range space.Spaceships {
		if spaceship.InputThrust {
			deltaVelocity := spaceship.getNormalizedFacing().Multiply(Acceleration)
			spaceship.Velocity = spaceship.Velocity.Add(deltaVelocity)
		}
		spaceship.Position = spaceship.Position.Add(spaceship.Velocity)

		angle := math.Atan2(spaceship.Facing.Y, spaceship.Facing.X)
		switch spaceship.InputDirection {
		case LEFT:
			angle -= AngularVelocity
		case RIGHT:
			angle += AngularVelocity
		}

		spaceship.Facing = NewVector(math.Cos(angle), math.Sin(angle))
	}
}
