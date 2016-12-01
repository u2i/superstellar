package state

import (
	"math"
	"superstellar/backend/constants"
	"superstellar/backend/pb"
	"superstellar/backend/types"
)

// Projectile struct holds players' shots data.
type Projectile struct {
	ID        uint32
	ClientID  uint32
	Spaceship *Spaceship
	FrameID   uint32
	Facing    float32
	Origin    *types.Point
	Velocity  *types.Vector
	Position  *types.Point
	TTL       uint32
}

// NewProjectile returns new instance of Projectile
func NewProjectile(ID, frameID uint32, spaceship *Spaceship) *Projectile {
	facing := float32(math.Atan2(-spaceship.Facing.Y, spaceship.Facing.X))

	return &Projectile{
		ID:        ID,
		ClientID:  spaceship.ID,
		Spaceship: spaceship,
		FrameID:   frameID,
		Origin:    spaceship.Position,
		Position:  spaceship.Position,
		Facing:    facing,
		Velocity:  spaceship.Facing.Multiply(constants.ProjectileSpeed).Add(spaceship.Velocity),
		TTL:       constants.ProjectileDefaultTTL,
	}
}

// ToProto returns protobuf representation
func (projectile *Projectile) ToProto() *pb.ProjectileFired {
	return &pb.ProjectileFired{
		Id:       projectile.ID,
		FrameId:  projectile.FrameID,
		Origin:   projectile.Origin.ToProto(),
		Ttl:      projectile.TTL,
		Facing:   projectile.Facing,
		Velocity: projectile.Velocity.ToProto(),
	}
}

func (projectile *Projectile) ToMessage() *pb.Message {
	return &pb.Message{
		Content: &pb.Message_ProjectileFired{
			ProjectileFired: projectile.ToProto(),
		},
	}
}

func (projectile *Projectile) ToHitMessage() *pb.Message {
	return &pb.Message{
		Content: &pb.Message_ProjectileHit{
			ProjectileHit: &pb.ProjectileHit{
				Id: projectile.ID,
			},
		},
	}
}

func (projectile *Projectile) DetectCollision(spaceship *Spaceship) bool {
	vA := types.Point{X: projectile.Position.X - spaceship.Position.X, Y: projectile.Position.Y - spaceship.Position.Y}
	distA := vA.Length()

	endPoint := projectile.Position.Add(projectile.Velocity)
	vB := types.Point{X: endPoint.X - spaceship.Position.X, Y: endPoint.Y - spaceship.Position.Y}
	distB := vB.Length()

	return distA < constants.SpaceshipSize || distB < constants.SpaceshipSize
}
