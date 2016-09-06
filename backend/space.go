package backend

type Space struct {
	Spaceships map[int]*Spaceship `json:"spaceships"`
}
