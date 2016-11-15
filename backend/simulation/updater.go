package simulation

import (
	"superstellar/backend/events"
	"superstellar/backend/state"
	"superstellar/backend/monitor"
	"time"
	"superstellar/backend/constants"
)

type Updater struct {
	space           *state.Space
	monitor         *monitor.Monitor
	eventDispatcher *events.EventDispatcher
}

func NewUpdater(space *state.Space, monitor *monitor.Monitor, eventDispatcher *events.EventDispatcher) *Updater {
	return &Updater{
		space: space,
		monitor: monitor,
		eventDispatcher: eventDispatcher,
	}
}

func (updater *Updater) HandleUserInput(userInputEvent *events.UserInput) {
	spaceship, found := updater.space.Spaceships[userInputEvent.ClientID]

	if found {
		spaceship.UpdateUserInput(userInputEvent.UserInput)
	}
}

func (updater *Updater) HandleTimeTick(*events.TimeTick) {
	before := time.Now()

	UpdatePhysics(updater.space, updater.eventDispatcher)

	elapsed := time.Since(before)
	updater.monitor.AddPhysicsTime(elapsed)
}

func (updater *Updater) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	updater.space.NewSpaceship(userJoinedEvent.ClientID)
}

func (updater *Updater) HandleUserLeft(userLeftEvent *events.UserLeft) {
	updater.space.RemoveSpaceship(userLeftEvent.ClientID)
}

func (updater *Updater) HandleUserDied(event *events.UserDied) {
	shotSpaceshipMaxHP := event.ShotSpaceship.MaxHP
	reward := uint32(float32(shotSpaceshipMaxHP) * constants.KillRewardRatio)

	updater.space.AddRewardToClient(event.KilledBy, reward)
}
