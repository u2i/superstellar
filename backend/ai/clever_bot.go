package ai

import (
	"superstellar/backend/state"
	"math/rand"
	"superstellar/backend/types"
	"math"
)

type CleverBot struct {
	targetSelected    bool
	targetSpaceshipId uint32
	changeTargetIn    uint32
	fireIn            uint32
}

func NewCleverBot() *CleverBot {
	return &CleverBot{false, 0, 300, 20}
}

func (b *CleverBot) HandleStateUpdate(space *state.Space, spaceship *state.Spaceship) {
	target := b.selectTarget(space, spaceship)

	// TODO: fix this hack
	if spaceship == nil {
		return
	}

	targetPosition := target.ObjectState.Position()
	botPosition := spaceship.ObjectState.Position()

	targetVector := types.NewVector(
		float64(targetPosition.X-botPosition.X)+(rand.Float64()-0.5)*10,
		-float64(targetPosition.Y - botPosition.Y)+(rand.Float64()-0.5)*10,
	)

	newTargetAngle := math.Atan2(targetVector.Y, targetVector.X)

	spaceship.UpdateTargetAngle(newTargetAngle)
	spaceship.TurnToTarget()

	if targetVector.Length() > 40000 {
		if spaceship.InputThrust != true {
			spaceship.InputThrust = true
			spaceship.MarkDirty()
		}

		if spaceship.Fire != false {
			spaceship.Fire = false
			spaceship.MarkDirty()
		}
	} else {
		if spaceship.InputThrust != false {
			spaceship.InputThrust = false
			spaceship.MarkDirty()
		}

		b.fireIn--
		if b.fireIn == 0 {
			if spaceship.Fire != true {
				spaceship.Fire = true
				spaceship.MarkDirty()
			}
			b.fireIn = uint32(rand.Intn(80)) + 20
		} else {
			spaceship.Fire = false
			spaceship.MarkDirty()
		}
	}
}

func (b *CleverBot) selectTarget(space *state.Space, botSpaceship *state.Spaceship) *state.Spaceship {
	target, targetExists := space.Spaceships[b.targetSpaceshipId]

	for !targetExists || !b.targetSelected || botSpaceship == target || b.changeTargetIn == 0 {
		b.changeTargetRandomly(space)
		b.changeTargetIn = uint32(rand.Intn(300)) + 300
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
