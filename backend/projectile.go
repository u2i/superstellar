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
	ID       uint32
	ClientID uint32
	FrameID  uint32
	Facing   float32
	Origin   *Point
	Velocity *Vector
	Position *Point
	TTL      uint32
}

// NewProjectile returns new instance of Projectile
func NewProjectile(ID, frameID uint32, spaceship *Spaceship) *Projectile {
	facing := float32(math.Atan2(-spaceship.Facing.Y, spaceship.Facing.X))

	return &Projectile{
		ID:       ID,
		ClientID: spaceship.ID,
		FrameID:  frameID,
		Origin:   spaceship.Position,
		Position: spaceship.Position,
		Facing:   facing,
		Velocity: spaceship.Facing.Multiply(ProjectileSpeed).Add(spaceship.Velocity),
		TTL:      DefaultTTL,
	}
}

func (projectile *Projectile) toProto() *pb.ProjectileFired {
	return &pb.ProjectileFired{
		Id:       projectile.ID,
		FrameId:  projectile.FrameID,
		Origin:   projectile.Origin.toProto(),
		Ttl:      projectile.TTL,
		Facing:   projectile.Facing,
		Velocity: projectile.Velocity.toProto(),
	}
}
