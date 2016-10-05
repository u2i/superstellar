package backend

import (
	"math"
	"math/rand"
	"superstellar/backend/pb"
)

const (
	// Radius of playable world (in .01 units)
	WorldRadius = 100000

	// Width of boundary region (in .01 units), i.e. from WorldRadius till when no more movement is possible
	BoundaryAnnulusWidth = 20000
)

// Space struct holds entire game state.
type Space struct {
	Spaceships map[uint32]*Spaceship `json:"spaceships"`
}

// NewSpace initializes new Space.
func NewSpace() *Space {
	return &Space{Spaceships: make(map[uint32]*Spaceship)}
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
	for _, spaceship := range space.Spaceships {
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
			angle -= AngularVelocity
		case RIGHT:
			angle += AngularVelocity
		}

		spaceship.Facing = NewVector(math.Cos(angle), math.Sin(angle))
	}
}

func (space *Space) toProto() *pb.Space {
	protoSpaceships := make([]*pb.Spaceship, 0, len(space.Spaceships))
	for _, spaceship := range space.Spaceships {
		protoSpaceships = append(protoSpaceships, spaceship.toProto())
	}

  return &pb.Space{Spaceships: protoSpaceships}
}

func (space *Space) toMessage() *pb.Message {
	return &pb.Message{
		Content: &pb.Message_Space{
			Space: space.toProto(),
		},
	}
}
