package ai

import (
	"math"
	"math/rand"
	"superstellar/backend/state"
	"superstellar/backend/types"
)

const (
	AverageChangeTargetAfter          = 300
	ChangeTargetAfterVariance         = 300
	AverageFireDelay                  = 60
	FireDelayVariance                 = 40
	GoCloserIfDistanceBiggerThan      = 40000
	TargetingAngleDifferenceThreshold = math.Pi / 20
)

type CleverBot struct {
	targetSelected    bool
	targetSpaceshipId uint32
	changeTargetIn    int
	fireIn            int
	changeDirectionIn int
}

func NewCleverBot() *CleverBot {
	return &CleverBot{
		targetSelected:    false,
		targetSpaceshipId: 0,
		changeTargetIn:    AverageChangeTargetAfter,
		fireIn:            AverageFireDelay,
	}
}

func (b *CleverBot) HandleStateUpdate(space *state.Space, spaceship *state.Spaceship) {
	// TODO: fix this hack
	if spaceship == nil {
		return
	}

	target := b.selectTarget(space, spaceship)

	if target == nil {
		return
	}

	// it looks like after projectile hit spaceship gets pretty high angular velocity
	// and it interferes somehow with bot's logic causing 'nietoperzyca' disease
	if math.Abs(spaceship.AngularVelocity()) > 0.5 {
		return
	}

	targetPosition := target.ObjectState.Position()
	botPosition := spaceship.ObjectState.Position()
	facingAngle := spaceship.ObjectState.Facing()

	targetVector := types.NewVector(
		float64(targetPosition.X-botPosition.X),
		float64(targetPosition.Y-botPosition.Y),
	)

	targetDirection := targetVector.Radians()

	angleDifference := targetDirection - facingAngle
	if angleDifference > math.Pi {
		angleDifference -= 2 * math.Pi
	}

	absAngleDifference := math.Abs(angleDifference)

	if absAngleDifference > TargetingAngleDifferenceThreshold {
		if angleDifference < 0 {
			b.changeDirection(spaceship, state.LEFT)
		} else {
			b.changeDirection(spaceship, state.RIGHT)
		}
	} else {
		b.changeDirection(spaceship, state.NONE)
	}

	b.fireIn--
	if targetVector.Length() > GoCloserIfDistanceBiggerThan {
		b.changeThrust(spaceship, true)

		b.changeFire(spaceship, false)
	} else {
		b.changeThrust(spaceship, false)

		if b.fireIn <= 0 {
			b.changeFire(spaceship, true)
			b.fireIn = rand.Intn(FireDelayVariance) - FireDelayVariance/2 + AverageFireDelay
		} else {
			b.changeFire(spaceship, false)
		}
	}
}

func (b *CleverBot) changeDirection(spaceship *state.Spaceship, direction state.Direction) {
	if spaceship.InputDirection != direction {
		spaceship.InputDirection = direction
		spaceship.MarkDirty()
	}
}

func (b *CleverBot) changeThrust(spaceship *state.Spaceship, thrustEnabled bool) {
	if spaceship.InputThrust != thrustEnabled {
		spaceship.InputThrust = thrustEnabled
		spaceship.MarkDirty()
	}
}

func (b *CleverBot) changeFire(spaceship *state.Spaceship, fireEnabled bool) {
	if spaceship.StraightFire != fireEnabled {
		spaceship.StraightFire = fireEnabled
		spaceship.MarkDirty()
	}
}

func (b *CleverBot) selectTarget(space *state.Space, botSpaceship *state.Spaceship) *state.Spaceship {
	if len(space.Spaceships) <= 1 {
		return nil
	}

	target, targetExists := space.Spaceships[b.targetSpaceshipId]

	for !targetExists || !b.targetSelected || botSpaceship == target || b.changeTargetIn == 0 {
		b.changeTargetRandomly(space)
		b.changeTargetIn = rand.Intn(ChangeTargetAfterVariance) - ChangeTargetAfterVariance/2 + AverageChangeTargetAfter
		target, targetExists = space.Spaceships[b.targetSpaceshipId]
	}

	b.changeTargetIn--

	return target
}

func (b *CleverBot) changeTargetRandomly(space *state.Space) {
	numberOfSpaceships := len(space.Spaceships)
	randomSpaceship := rand.Intn(numberOfSpaceships)

	for spaceshipId := range space.Spaceships {
		randomSpaceship -= 1

		if randomSpaceship <= 0 {
			b.targetSpaceshipId = spaceshipId
			b.targetSelected = true
			break
		}
	}
}
