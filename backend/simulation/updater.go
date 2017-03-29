package simulation

import (
	"log"
	"superstellar/backend/events"
	"superstellar/backend/monitor"
	"superstellar/backend/state"
	"superstellar/backend/utils"
	"time"
)

type Updater struct {
	space             *state.Space
	spaceshipManager  *SpaceshipManager
	objectManager     *ObjectManager
	collisionManager  *CollisionManager
	projectileManager *ProjectileManager
	asteroidManager   *AsteroidManager
	monitor           *monitor.Monitor
	eventDispatcher   *events.EventDispatcher
	idSequencer       *utils.IdManager
}

func NewUpdater(space *state.Space, monitor *monitor.Monitor, eventDispatcher *events.EventDispatcher,
	idSequencer *utils.IdManager) *Updater {
	return &Updater{
		space:             space,
		spaceshipManager:  NewSpaceshipManager(space, eventDispatcher),
		objectManager:     NewObjectManager(space),
		collisionManager:  NewCollisionManager(space, eventDispatcher),
		projectileManager: NewProjectileManager(space, eventDispatcher),
		asteroidManager:   NewAsteroidManager(space, idSequencer),
		monitor:           monitor,
		eventDispatcher:   eventDispatcher,
		idSequencer:       idSequencer,
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

	updater.updatePhysics()

	if updater.space.PhysicsFrameID == 1 {
		log.Println("Simulation start timestamp:", time.Now().UnixNano()/time.Millisecond.Nanoseconds())
	}

	elapsed := time.Since(before)
	updater.monitor.AddPhysicsTime(elapsed)
}

func (updater *Updater) updatePhysics() {
	updater.projectileManager.detectProjectileCollisions()
	updater.asteroidManager.updateAsteroids()
	updater.spaceshipManager.updateSpaceships()
	updater.objectManager.updateObjects()
	updater.collisionManager.resolveCollisions()

	updater.space.PhysicsFrameID++
	updater.eventDispatcher.FirePhysicsReady(&events.PhysicsReady{})

	updater.projectileManager.updateProjectiles()
}

func (updater *Updater) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	updater.space.NewSpaceship(userJoinedEvent.ClientID)
}

func (updater *Updater) HandleUserLeft(userLeftEvent *events.UserLeft) {
	spaceship := updater.space.Spaceships[userLeftEvent.ClientID]

	if spaceship != nil {
		updater.space.RemoveSpaceship(spaceship.Id())

		objectDestroyedMessage := &events.ObjectDestroyed{
			DestroyedObject: spaceship,
			DestroyedBy:     spaceship,
			Timestamp:       time.Now(),
		}

		updater.eventDispatcher.FireObjectDestroyed(objectDestroyedMessage)
	}
}

func (updater *Updater) HandleObjectDestroyed(event *events.ObjectDestroyed) {
	if event.DestroyedBy != nil && event.DestroyedBy != event.DestroyedObject {
		event.DestroyedBy.ObjectDestroyed(event.DestroyedObject)
	}
}
