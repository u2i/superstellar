package events

import (
	"superstellar/backend/state"
	"time"
)

type ObjectDestroyed struct {
	//Id uint32
	DestroyedObject state.Object
	DestroyedBy     state.Object
	//DestroyedBy uint32
	//ShotSpaceship *state.Spaceship
	Timestamp time.Time
}
