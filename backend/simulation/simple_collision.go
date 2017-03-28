package simulation

import (
	"math"
	"superstellar/backend/state"
	"superstellar/backend/types"
)

type SimpleCollision struct{}

func (*SimpleCollision) collide(objectA state.Object, objectB state.Object) {
	v := types.Point{
		X: objectA.Position().X - objectB.Position().X,
		Y: objectA.Position().Y - objectB.Position().Y,
	}

	transformAngle := -math.Atan2(float64(v.Y), float64(v.X))
	newV1 := objectA.Velocity().Rotate(transformAngle)
	newV2 := objectB.Velocity().Rotate(transformAngle)

	switchedV1 := types.Vector{X: newV2.X, Y: newV1.Y}
	switchedV2 := types.Vector{X: newV1.X, Y: newV2.Y}

	objectA.SetVelocity(switchedV1.Rotate(-transformAngle))
	objectB.SetVelocity(switchedV2.Rotate(-transformAngle))

	objectA.MarkDirty()
	objectB.MarkDirty()

	objectA.CollideWith(objectB)
	objectB.CollideWith(objectA)
}
