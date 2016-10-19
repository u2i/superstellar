package space

import (
	"math/rand"
	"superstellar/backend/pb"
)

const (
	// WorldRadius is the radius of playable world (in .01 units)
	WorldRadius = 100000

	// BoundaryAnnulusWidth is the width of boundary region (in .01 units), i.e. from WorldRadius till when no more movement is possible
	BoundaryAnnulusWidth = 20000
)

// Space struct holds entire game state.
type Space struct {
	ShotsCh               chan *Projectile
	Spaceships            map[uint32]*Spaceship `json:"spaceships"`
	Projectiles           map[*Projectile]bool
	PhysicsFrameID        uint32
	NextProjectileIDValue uint32
}

// NewSpace initializes new Space.
func NewSpace(shotsCh chan *Projectile) *Space {
	return &Space{
		ShotsCh:               shotsCh,
		Spaceships:            make(map[uint32]*Spaceship),
		Projectiles:           make(map[*Projectile]bool),
		PhysicsFrameID:        0,
		NextProjectileIDValue: 0,
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

// AddProjectile adds projectile to the space.``
func (space *Space) AddProjectile(projectile *Projectile) {
	space.Projectiles[projectile] = true
}

// RemoveProjectile removes projectile from the space.
func (space *Space) RemoveProjectile(projectile *Projectile) {
	delete(space.Projectiles, projectile)
}

// UpdateUserInput updates user input in correct spaceship
func (space *Space) UpdateUserInput(userInput *UserInput) {
	spaceship, found := space.Spaceships[userInput.ClientID]

	if found {
		spaceship.updateUserInput(userInput)
	}
}

// NextProjectileID returns next unused projectile ID.
func (space *Space) NextProjectileID() uint32 {
	ID := space.NextProjectileIDValue
	space.NextProjectileIDValue++
	return ID
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

// ToProto returns protobuf representation
func (space *Space) ToProto() *pb.Space {
	protoSpaceships := make([]*pb.Spaceship, 0, len(space.Spaceships))
	for _, spaceship := range space.Spaceships {
		protoSpaceships = append(protoSpaceships, spaceship.ToProto())
	}

	return &pb.Space{Spaceships: protoSpaceships, PhysicsFrameID: space.PhysicsFrameID}
}

// ToMessage returns protobuffer Message object with Space set.
func (space *Space) ToMessage() *pb.Message {
	return &pb.Message{
		Content: &pb.Message_Space{
			Space: space.ToProto(),
		},
	}
}
