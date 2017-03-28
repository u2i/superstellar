package state

import (
	"superstellar/backend/constants"
	"superstellar/backend/pb"
	"superstellar/backend/types"
)

type Asteroid struct {
	ObjectState
}

func NewAsteroid(id uint32, initialPosition *types.Point, initialVelocity *types.Vector) *Asteroid {
	objectState := NewObjectState(id, initialPosition, initialVelocity)

	return &Asteroid{
		ObjectState: *objectState,
	}
}

// ToProto returns protobuf representation
func (asteroid *Asteroid) ToProto() *pb.Asteroid {
	return &pb.Asteroid{
		Id:              asteroid.Id(),
		Position:        asteroid.Position().ToProto(),
		Velocity:        asteroid.Velocity().ToProto(),
		Facing:          asteroid.Facing(),
		AngularVelocity: asteroid.AngularVelocity(),
	}
}

func (asteroid *Asteroid) AddToProtoSpace(space *pb.Space) {
	space.Asteroids = append(space.Asteroids, asteroid.ToProto())
}

func (asteroid *Asteroid) CollideWith(other Object) {

}

func (asteroid *Asteroid) DamageValue() uint32 {
	return constants.AsteroidDamageValue
}
