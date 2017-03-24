package simulation

import (
	"log"
	"superstellar/backend/constants"
	"superstellar/backend/events"
	"superstellar/backend/monitor"
	"superstellar/backend/state"
	"time"
)

type Updater struct {
	space            *state.Space
	collisionManager *CollisionManager
	monitor          *monitor.Monitor
	eventDispatcher  *events.EventDispatcher
}

func NewUpdater(space *state.Space, monitor *monitor.Monitor, eventDispatcher *events.EventDispatcher) *Updater {
	return &Updater{
		space:            space,
		collisionManager: NewCollisionManager(),
		monitor:          monitor,
		eventDispatcher:  eventDispatcher,
	}
}

func (updater *Updater) HandleUserInput(userInputEvent *events.UserInput) {
	spaceship, found := updater.space.Spaceships[userInputEvent.ClientID]

	if found {
		spaceship.UpdateUserInput(userInputEvent.UserInput)
	}
}

func (updater *Updater) HandleTargetAngle(targetAngleEvent *events.TargetAngle) {
	spaceship, found := updater.space.Spaceships[targetAngleEvent.ClientID]

	if found {
		spaceship.UpdateTargetAngle(targetAngleEvent.Angle)
	}
}

func (updater *Updater) HandleTimeTick(*events.TimeTick) {
	before := time.Now()

	UpdatePhysics(updater.space, updater.eventDispatcher, updater.collisionManager)

	if updater.space.PhysicsFrameID == 1 {
		log.Println("Simulation start timestamp:", time.Now().UnixNano()/time.Millisecond.Nanoseconds())
	}

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
	energyReward := uint32(float32(shotSpaceshipMaxHP) * constants.KillEnergyRewardRatio)

	event.Shooter.AddReward(reward)
	event.Shooter.AddEnergyReward(energyReward)
}
