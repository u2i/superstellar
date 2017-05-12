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
	targetAngle := 0.0

	if spaceship.TurretFire && spaceship.TargetAngle != nil {
		targetAngle = *spaceship.TargetAngle
	} else {
		targetAngle = spaceship.Facing()
	}
	facingVector := types.NewVector(math.Cos(targetAngle), -math.Sin(targetAngle))

	return &Projectile{
		ID:        ID,
		ClientID:  spaceship.Id(),
		Spaceship: spaceship,
		FrameID:   frameID,
		Origin:    spaceship.Position(),
		Position:  spaceship.Position(),
		Facing:    float32(targetAngle),
		Velocity:  facingVector.Multiply(constants.ProjectileSpeed).Add(spaceship.Velocity()),
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

func (projectile *Projectile) DetectCollision(object Object) (bool, *types.Point) {
	vA := types.Point{X: projectile.Position.X - object.Position().X, Y: projectile.Position.Y - object.Position().Y}
	distA := vA.Length()

	endPoint := projectile.Position.Add(projectile.Velocity)
	vB := types.Point{X: endPoint.X - object.Position().X, Y: endPoint.Y - object.Position().Y}
	distB := vB.Length()

	if distA < constants.SpaceshipSize {
		return true, projectile.Position
	} else if distB < constants.SpaceshipSize {
		return true, endPoint
	}

	return false, nil
}
