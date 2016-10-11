package backend

import (
	"math"
	"math/rand"
	"superstellar/backend/pb"
	"time"
)

const (
	// WorldRadius is the radius of playable world (in .01 units)
	WorldRadius = 100000

	// BoundaryAnnulusWidth is the width of boundary region (in .01 units), i.e. from WorldRadius till when no more movement is possible
	BoundaryAnnulusWidth = 20000
)

// Space struct holds entire game state.
type Space struct {
	ShotsCh        chan *Projectile
	Spaceships     map[uint32]*Spaceship `json:"spaceships"`
	PhysicsFrameID uint32
}

// NewSpace initializes new Space.
func NewSpace(shotsCh chan *Projectile) *Space {
	return &Space{
		ShotsCh:        shotsCh,
		Spaceships:     make(map[uint32]*Spaceship),
		PhysicsFrameID: 0,
	}
}

// AddSpaceship adds new spaceship to the space.
func (space *Space) AddSpaceship(clientID uint32, spaceship *Spaceship) {
	space.Spaceships[clientID] = spaceship
}

// RemoveSpaceship removes spaceship from the space.
func (space *Space) RemoveSpaceship(clientID uint32) {
	delete(space.Spaceships, clientID)
}

// UpdateUserInput updates user input in correct spaceship
func (space *Space) UpdateUserInput(userInput *UserInput) {
	spaceship, found := space.Spaceships[userInput.ClientID]

	if found {
		spaceship.updateUserInput(userInput)
	}
}

func (space *Space) randomUpdate() {
	for _, e := range space.Spaceships {
		if rand.Float64() < 0.05 {
			e.InputThrust = !e.InputThrust
		}
		if rand.Float64() < 0.03 {
			e.InputDirection = Direction(rand.Int() % 3)
		}
	}
}

func (space *Space) updatePhysics() {
	now := time.Now()

	for _, spaceship := range space.Spaceships {
		if spaceship.Fire {
			timeSinceLastShot := now.Sub(spaceship.LastShotTime)
			if timeSinceLastShot >= MinFireInterval {
				projectile := NewProjectile(spaceship, space.PhysicsFrameID)
				space.ShotsCh <- projectile
				spaceship.LastShotTime = now
			}
		}

		if spaceship.InputThrust {
			deltaVelocity := spaceship.getNormalizedFacing().Multiply(Acceleration)
			spaceship.Velocity = spaceship.Velocity.Add(deltaVelocity)
			if spaceship.Velocity.Length() > MaxSpeed {
				spaceship.Velocity = spaceship.Velocity.Normalize().Multiply(MaxSpeed)
			}
		}

		spaceship.Position = spaceship.Position.Add(spaceship.Velocity)
		if spaceship.Position.Length() > WorldRadius {
			outreachLength := spaceship.Position.Length() - WorldRadius
			gravityAcceleration := -(outreachLength / BoundaryAnnulusWidth) * Acceleration
			deltaVelocity := spaceship.Position.Normalize().Multiply(gravityAcceleration)
			spaceship.Velocity = spaceship.Velocity.Add(deltaVelocity)
		}

		angle := math.Atan2(spaceship.Facing.Y, spaceship.Facing.X)
		switch spaceship.InputDirection {
		case LEFT:
			angle += AngularVelocity
		case RIGHT:
			angle -= AngularVelocity
		}

		spaceship.Facing = NewVector(math.Cos(angle), math.Sin(angle))
	}

	space.PhysicsFrameID++
}

func (space *Space) toProto() *pb.Space {
	protoSpaceships := make([]*pb.Spaceship, 0, len(space.Spaceships))
	for _, spaceship := range space.Spaceships {
		protoSpaceships = append(protoSpaceships, spaceship.toProto())
	}

	return &pb.Space{Spaceships: protoSpaceships, PhysicsFrameID: space.PhysicsFrameID}
}

func (space *Space) toMessage() *pb.Message {
	return &pb.Message{
		Content: &pb.Message_Space{
			Space: space.toProto(),
		},
	}
}
