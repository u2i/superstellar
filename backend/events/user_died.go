package events

import "superstellar/backend/state"

type UserDied struct {
	ClientID	uint32
	KilledBy	uint32
	ShotSpaceship   *state.Spaceship
}
