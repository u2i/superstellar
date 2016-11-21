package events

import "superstellar/backend/state"

type UserDied struct {
	ClientID      uint32
	Shooter       *state.Spaceship
	KilledBy      uint32
	ShotSpaceship *state.Spaceship
}
