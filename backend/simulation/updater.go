package simulation

import (
	"superstellar/backend/events"
	"superstellar/backend/state"
)

type Updater struct {
	space *state.Space
	eventDispatcher *events.EventDispatcher
}

func NewUpdater(space *state.Space, eventDispatcher *events.EventDispatcher) *Updater {
	return &Updater{space, eventDispatcher}
}

func (updater *Updater) HandleUserInput(userInputEvent *events.UserInput) {
	spaceship, found := updater.space.Spaceships[userInputEvent.ClientID]

	if found {
		spaceship.UpdateUserInput(userInputEvent.UserInput)
	}
}

func (updater *Updater) HandleTimeTick(*events.TimeTick) {
	UpdatePhysics(updater.space, updater.eventDispatcher)
}