package backend

import (
	"math"
	"superstellar/backend/pb"
)

const (
	DefaultTtl      = 50
	ProjectileSpeed = 2000
	// TODO: remove Konarski's factor
	KonarskiFactor = 100
)

// Shot struct holds players' shots data.
type Projectile struct {
	ClientID uint32
	FrameID  uint32
	Facing   float32
	Origin   *IntVector
	Velocity *Vector
	Position *IntVector
}

// NewProjectile returns new instance of Projectile
func NewProjectile(spaceship *Spaceship, frameID uint32) *Projectile {
	facing := float32(math.Atan2(-spaceship.Facing.Y, spaceship.Facing.X))

	return &Projectile{
		ClientID: spaceship.ID,
		FrameID:  frameID,
		Origin:   spaceship.Position,
		Position: spaceship.Position,
		Facing:   facing,
		// TODO: remove Konarski's factor
		Velocity: spaceship.Facing.Multiply(ProjectileSpeed).Add(spaceship.Velocity.Multiply(KonarskiFactor)),
	}
}

func (projectile *Projectile) toProto() *pb.ProjectileFired {
	return &pb.ProjectileFired{
		FrameId:  projectile.FrameID,
		Origin:   projectile.Origin.toProto(),
		Ttl:      DefaultTtl,
		Facing:   projectile.Facing,
		Velocity: projectile.Velocity.toProto(),
	}
}
