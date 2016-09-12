package backend

import "math/rand"

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
	for _, v := range space.Spaceships {
		facingDiff := &Vector{X: rand.Float64() - 0.5, Y: rand.Float64() - 0.5}
		v.Facing = v.Facing.Add(facingDiff.Normalize().Multiply(0.15).Normalize())
		positionDiff := &Vector{X: rand.Float64() - 0.5, Y: rand.Float64() - 0.5}
		v.Position = v.Position.Add(positionDiff.Normalize().Multiply(5))
	}
}
