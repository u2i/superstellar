package space

import (
	"log"
	"math"
	"math/rand"
	"superstellar/backend/pb"
	"superstellar/backend/types"
	"time"
)

const (
	// WorldRadius is the radius of playable world (in .01 units)
	WorldRadius = 100000

	// BoundaryAnnulusWidth is the width of boundary region (in .01 units), i.e. from WorldRadius till when no more movement is possible
	BoundaryAnnulusWidth = 20000

	// RandomPositionEmptyRadius describes the minimum radius around randomized
	// initial position that needs to be free of any objects.
	RandomPositionEmptyRadius = 5000.0
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

// NewSpaceship creates a new spaceship and adds it to the space.
func (space *Space) NewSpaceship(clientID uint32) {
	spaceship := &Spaceship{
		ID:             clientID,
		Position:       space.randomEmptyPosition(),
		Velocity:       types.ZeroVector(),
		Facing:         types.NewVector(0.0, 1.0),
		InputThrust:    false,
		InputDirection: NONE,
		Fire:           false,
		LastShotTime:   time.Now(),
		HP:             InitialHP,
		MaxHP:          InitialHP,
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
		if space.furtherFromAnyObject(position, RandomPositionEmptyRadius) {
			return position
		}
	}
}

func (space *Space) randomPoint() *types.Point {
	angle := rand.Float64() * 2 * math.Pi
	radius := rand.Uint32() % (WorldRadius + 1)

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
