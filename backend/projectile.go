package backend

import (
	"math"
	"superstellar/backend/pb"
)

const (
	// DefaultTTL describes the default number of frames the projectile lives.
	DefaultTTL = 50

	// ProjectileSpeed describes projectile speed. Captain Obvious.
	ProjectileSpeed = 2000
)

// Projectile struct holds players' shots data.
type Projectile struct {
	ClientID uint32
	FrameID  uint32
	Facing   float32
	Origin   *Point
	Velocity *Vector
	Position *Point
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
		Velocity: spaceship.Facing.Multiply(ProjectileSpeed).Add(spaceship.Velocity),
	}
}

func (projectile *Projectile) toProto() *pb.ProjectileFired {
	return &pb.ProjectileFired{
		FrameId:  projectile.FrameID,
		Origin:   projectile.Origin.toProto(),
		Ttl:      DefaultTTL,
		Facing:   projectile.Facing,
		Velocity: projectile.Velocity.toProto(),
	}
}
