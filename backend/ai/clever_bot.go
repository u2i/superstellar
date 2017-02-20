package ai

import "superstellar/backend/state"

type CleverBot struct {

}

func NewCleverBot() *CleverBot {
	return &CleverBot{}
}

func (b *CleverBot) HandleStateUpdate(space *state.Space, spaceship *state.Spaceship) {

}
