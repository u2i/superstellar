package ai

import "superstellar/backend/state"

type Bot interface {
	HandleStateUpdate(space *state.Space, spaceship *state.Spaceship)
}
