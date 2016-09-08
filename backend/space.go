package backend

import "math/rand"

type Space struct {
	Spaceships map[string]*Spaceship `json:"spaceships"`
}

func NewSpace() *Space {
	return &Space{Spaceships: make(map[string]*Spaceship)}
}

func (space *Space) AddSpaceship(clientId string, spaceship *Spaceship) {
	space.Spaceships[clientId] = spaceship
}

func (space *Space) RemoveSpaceship(clientId string) {
	delete(space.Spaceships, clientId)
}

func (space *Space) randomUpdate() {
	for _, v := range space.Spaceships {
		facingDiff := &Vector{X: rand.Float64() - 0.5, Y: rand.Float64() - 0.5}
		v.Facing = v.Facing.Add(facingDiff.Normalize().Multiply(0.15).Normalize())
		positionDiff := &Vector{X: rand.Float64() - 0.5, Y: rand.Float64() - 0.5}
		v.Position = v.Position.Add(positionDiff.Normalize().Multiply(5))
	}
}
