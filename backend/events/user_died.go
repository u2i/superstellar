package events

import (
	"superstellar/backend/state"
	"time"
)

type UserDied struct {
	ClientID      uint32
	Shooter       *state.Spaceship
	KilledBy      uint32
	ShotSpaceship *state.Spaceship
	Timestamp     time.Time
}
