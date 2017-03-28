package state

import (
	"math"
	"math/rand"
	"superstellar/backend/constants"
	"superstellar/backend/pb"
	"superstellar/backend/types"
)

// Space struct holds entire game state.
type Space struct {
	Objects               map[uint32]Object
	Spaceships            map[uint32]*Spaceship
	Asteroids             map[uint32]*Asteroid
	Projectiles           map[*Projectile]bool
	PhysicsFrameID        uint32
	NextProjectileIDValue uint32
}

// NewSpace initializes new Space.
func NewSpace() *Space {
	return &Space{
		Objects:               make(map[uint32]Object),
		Spaceships:            make(map[uint32]*Spaceship),
		Asteroids:             make(map[uint32]*Asteroid),
		Projectiles:           make(map[*Projectile]bool),
		PhysicsFrameID:        0,
		NextProjectileIDValue: 0,
	}
}

// NewSpaceship creates a new spaceship and adds it to the space.
func (space *Space) NewSpaceship(clientID uint32) *Spaceship {
	spaceship := NewSpaceship(clientID, space.randomEmptyPosition())
	space.AddSpaceship(clientID, spaceship)
	return spaceship
}

// AddSpaceship adds new spaceship to the space.
func (space *Space) AddSpaceship(clientID uint32, spaceship *Spaceship) {
	space.Spaceships[clientID] = spaceship
	space.Objects[clientID] = spaceship
}

// AddAsteroid adds new spaceship to the space.
func (space *Space) AddAsteroid(asteroid *Asteroid) {
	space.Asteroids[asteroid.id] = asteroid
	space.Objects[asteroid.id] = asteroid
}

// RemoveSpaceship removes spaceship from the space.
func (space *Space) RemoveSpaceship(id uint32) {
	space.RemoveObject(id)

}

// RemoveAsteroid removes spaceship from the space.
func (space *Space) RemoveAsteroid(id uint32) {
	space.RemoveObject(id)
}

// RemoveAsteroid removes spaceship from the space.
func (space *Space) RemoveObject(id uint32) {
	delete(space.Asteroids, id)
	delete(space.Spaceships, id)
	delete(space.Objects, id)
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

// ToProto returns protobuf representation
func (space *Space) ToProto(fullUpdate bool) *pb.Space {
	protoSpaceships := make([]*pb.Spaceship, 0, len(space.Spaceships))
	protoAsteroids := make([]*pb.Asteroid, 0, len(space.Asteroids))

	protoSpace := &pb.Space{Spaceships: protoSpaceships, Asteroids: protoAsteroids, PhysicsFrameID: space.PhysicsFrameID}

	for _, object := range space.Objects {
		if fullUpdate || object.Dirty() {
			object.AddToProtoSpace(protoSpace)
			if !fullUpdate {
				object.MarkClean()
			}
		}
	}

	return protoSpace
}

// ToMessage returns protobuffer Message object with Space set.
func (space *Space) ToMessage(fullUpdate bool) *pb.Message {
	return &pb.Message{
		Content: &pb.Message_Space{
			Space: space.ToProto(fullUpdate),
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
		positions = append(positions, spaceship.Position())
	}

	for projectile := range space.Projectiles {
		positions = append(positions, projectile.Position)
	}

	return positions
}
