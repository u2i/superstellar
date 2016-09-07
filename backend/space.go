package backend

type Space struct {
	Spaceships map[int]*Spaceship `json:"spaceships"`
}

func NewSpace() *Space {
	return &Space{Spaceships: make(map[int]*Spaceship)}
}

func (space *Space) AddSpaceship(clientId int, spaceship *Spaceship) {
	space.Spaceships[clientId] = spaceship
}

func (space *Space) RemoveSpaceship(clientId int) {
	delete(space.Spaceships, clientId)
}
