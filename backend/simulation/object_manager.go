package simulation

import (
	"math"
	"superstellar/backend/state"
)

type ObjectManager struct {
	space *state.Space
}

func NewObjectManager(space *state.Space) *ObjectManager {
	return &ObjectManager{
		space: space,
	}
}

func (manager *ObjectManager) updateObjects() {
	for _, object := range manager.space.Objects {

		// POSITION UPDATE

		object.SetPosition(object.Position().Add(object.Velocity()))

		// APPLY ANGULAR VELOCITY

		object.SetAngularVelocity(object.AngularVelocity() + object.AngularVelocityDelta())
		object.SetAngularVelocityDelta(0.0)

		object.SetFacing(object.Facing() - object.AngularVelocity())
		if math.Abs(object.Facing()) > math.Pi {
			object.SetFacing(object.Facing() - math.Copysign(2*math.Pi, object.Facing()))
		}

		// NOTIFY ABOUT NEW FRAME

		object.NotifyAboutNewFrame()
	}
}
