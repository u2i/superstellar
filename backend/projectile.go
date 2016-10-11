package backend

import "superstellar/backend/pb"

const (
	DefaultTtl = 50
	ProjectileSpeed = 2000
)

// Shot struct holds players' shots data.
type Projectile struct {
	ClientID uint32
	FrameID  uint32
	Origin   *IntVector
	Facing   *Vector
	Range    uint32
	Position *IntVector
}

// NewProjectile returns new instance of Projectile
func NewProjectile(clientID, frameID uint32, origin *IntVector, facing *Vector,
	shotRange uint32) *Projectile {
	return &Projectile{
		ClientID: clientID,
		FrameID:  frameID,
		Origin:   origin,
		Facing:   facing,
		Range:    shotRange,
		Position: origin,
	}
}

func (shot *Projectile) toProto() *pb.ProjectileFired {
	return &pb.ProjectileFired{
		FrameId: shot.FrameID,
		Origin:  shot.Origin.toProto(),
		Facing:  float32(shot.Facing.Radians()),
		Ttl: DefaultTtl,
		Speed: ProjectileSpeed,
	}
}
