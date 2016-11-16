package state

import (
	"log"
	"math"
	"math/rand"
	"superstellar/backend/constants"
	"superstellar/backend/pb"
	"superstellar/backend/types"
	"time"
)

// Space struct holds entire game state.
type Space struct {
	Spaceships            map[uint32]*Spaceship
	Projectiles           map[*Projectile]bool
	PhysicsFrameID        uint32
	NextProjectileIDValue uint32
}

// NewSpace initializes new Space.
func NewSpace() *Space {
	return &Space{
		Spaceships:            make(map[uint32]*Spaceship),
		Projectiles:           make(map[*Projectile]bool),
		PhysicsFrameID:        0,
		NextProjectileIDValue: 0,
	}
}

// NewSpaceship creates a new spaceship and adds it to the space.
func (space *Space) NewSpaceship(clientID uint32) {
	spaceship := &Spaceship{
		ID:              clientID,
		Position:        space.randomEmptyPosition(),
		Velocity:        types.ZeroVector(),
		Facing:          types.NewVector(0.0, 1.0),
		InputThrust:     false,
		InputDirection:  NONE,
		Fire:            false,
		LastShotTime:    time.Now(),
		HP:              constants.SpaceshipInitialHP,
		MaxHP:           constants.SpaceshipInitialHP,
		AutoRepairDelay: constants.AutoRepairDelay,
	}

	space.AddSpaceship(clientID, spaceship)
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

// NextProjectileID returns next unused projectile ID.
func (space *Space) NextProjectileID() uint32 {
	ID := space.NextProjectileIDValue
	space.NextProjectileIDValue++
	return ID
}

func (space *Space) AddRewardToClient(clientId uint32, reward uint32) {
	space.Spaceships[clientId].AddReward(reward)
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

// RandomEmptyPosition return randomized position within the space that is
// no closer to any object than given radius.
func (space *Space) randomEmptyPosition() *types.Point {
	for {
		position := space.randomPoint()
		if space.furtherFromAnyObject(position, constants.RandomPositionEmptyRadius) {
			return position
		}
	}
}

func (space *Space) randomPoint() *types.Point {
	angle := rand.Float64() * 2 * math.Pi
	radius := rand.Uint32() % (constants.WorldRadius + 1)

	return types.NewPointFromPolar(angle, radius)
}

func (space *Space) furtherFromAnyObject(position *types.Point, radius float64) bool {
	objectsPositions := space.allObjectsPositions()
	for _, objectPosition := range objectsPositions {
		if objectPosition.Distance(position) < radius {
			return false
		}
	}

	return true
}

func (space *Space) allObjectsPositions() []*types.Point {
	var positions []*types.Point

	for _, spaceship := range space.Spaceships {
		log.Printf("%d %d", spaceship.Position.X, spaceship.Position.Y)
		positions = append(positions, spaceship.Position)
	}

	for projectile := range space.Projectiles {
		log.Printf("%d %d", projectile.Position.X, projectile.Position.Y)
		positions = append(positions, projectile.Position)
	}

	return positions
}
