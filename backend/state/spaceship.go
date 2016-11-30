package state

import (
	"fmt"
	"math"
	"superstellar/backend/constants"
	"superstellar/backend/pb"
	"superstellar/backend/types"
	"time"
)

// Direction is a type describing user input on spaceship rotation.
type Direction int

// Constants describing user input on spaceship rotation.
const (
	NONE Direction = iota
	RIGHT
	LEFT
)

// Spaceship struct describes a spaceship.
type Spaceship struct {
	ID                      uint32
	Position                *types.Point
	Velocity                *types.Vector
	Facing                  *types.Vector
	AngularSpeed            float64
	InputThrust             bool
	InputDirection          Direction
	TargetAngle             *float64
	Fire                    bool
	LastShotTime            time.Time
	HP                      uint32
	MaxHP                   uint32
	Energy                  uint32
	MaxEnergy               uint32
	AutoRepairDelay         uint32
	AutoEnergyRechargeDelay uint32
}

func NewSpaceship(clientId uint32, initialPosition *types.Point) *Spaceship {
	return &Spaceship{
		ID:              clientId,
		Position:        initialPosition,
		Velocity:        types.ZeroVector(),
		Facing:          types.NewVector(0.0, 1.0),
		AngularSpeed:    0,
		InputThrust:     false,
		InputDirection:  NONE,
		Fire:            false,
		LastShotTime:    time.Now(),
		HP:              constants.SpaceshipInitialHP,
		MaxHP:           constants.SpaceshipInitialHP,
		Energy:          constants.SpaceshipInitialEnergy,
		MaxEnergy:       constants.SpaceshipInitialEnergy,
		AutoRepairDelay: constants.AutoRepairDelay,
	}
}

// String function returns string representation.
func (s *Spaceship) String() string {
	return fmt.Sprintf("(%v, %v, %v)", s.Position, s.Velocity, s.Facing)
}

// NormalizedFacing return normalized facing vector.
func (s *Spaceship) NormalizedFacing() *types.Vector {
	return s.Facing.Normalize()
}

func (s *Spaceship) UpdateUserInput(userInput pb.UserInput) {
	switch userInput {
	case pb.UserInput_CENTER:
		s.InputDirection = NONE
		s.TargetAngle = nil
	case pb.UserInput_LEFT:
		s.InputDirection = LEFT
		s.TargetAngle = nil
	case pb.UserInput_RIGHT:
		s.InputDirection = RIGHT
		s.TargetAngle = nil
	case pb.UserInput_THRUST_ON:
		s.InputThrust = true
	case pb.UserInput_THRUST_OFF:
		s.InputThrust = false
	case pb.UserInput_FIRE_START:
		s.Fire = true
	case pb.UserInput_FIRE_STOP:
		s.Fire = false
	}
}

func (s *Spaceship) UpdateTargetAngle(angle float64) {
	s.TargetAngle = &angle
	s.InputDirection = NONE
}

// ToProto returns protobuf representation
func (s *Spaceship) ToProto() *pb.Spaceship {
	return &pb.Spaceship{
		Id:          s.ID,
		Position:    s.Position.ToProto(),
		Velocity:    s.Velocity.ToProto(),
		Facing:      float32(s.Facing.Radians()),
		InputThrust: s.InputThrust,
		MaxHp:       s.MaxHP,
		Hp:          s.HP,
		MaxEnergy:   s.MaxEnergy,
		Energy:      s.Energy,
	}
}

// DetectCollision returns true if receiver spaceship collides with other spaceship.
func (s *Spaceship) DetectCollision(other *Spaceship) bool {
	v := types.Point{X: s.Position.X - other.Position.X, Y: s.Position.Y - other.Position.Y}
	dist := v.Length()

	return dist < 2 * constants.SpaceshipSize
}

// Collide transforms colliding ships' parameters.
func (s *Spaceship) Collide(other *Spaceship) {
	v := types.Point{
		X: s.Position.X - other.Position.X,
		Y: s.Position.Y - other.Position.Y,
	}

	transformAngle := -math.Atan2(float64(v.Y), float64(v.X))
	newV1 := s.Velocity.Rotate(transformAngle)
	newV2 := other.Velocity.Rotate(transformAngle)

	switchedV1 := types.Vector{X: newV2.X, Y: newV1.Y}
	switchedV2 := types.Vector{X: newV1.X, Y: newV2.Y}

	s.Velocity = switchedV1.Rotate(-transformAngle)
	other.Velocity = switchedV2.Rotate(-transformAngle)
}

func (s *Spaceship) ShootIfPossible() (canShoot bool) {
	if s.Energy >= constants.BasicWeaponEnergyCost {
		s.Energy -= constants.BasicWeaponEnergyCost;
		canShoot = true
	} else {
		canShoot = false
	}
	return canShoot
}

func (s *Spaceship) CollideWithProjectile(projectile *Projectile) {
	if s.HP < 50 {
		s.HP = 0
	} else {
		s.HP -= 50
	}
	s.AutoRepairDelay = constants.AutoRepairDelay
}

func (s *Spaceship) AddReward(reward uint32) {
	s.HP += reward
	s.MaxHP += reward
}

func (s *Spaceship) AddEnergyReward(reward uint32) {
	s.Energy += reward
	s.MaxEnergy += reward
}

func (s *Spaceship) AutoRepair() {
	s.HP += constants.AutoRepairAmount

	if (s.HP > s.MaxHP) {
		s.HP = s.MaxHP
	}
	s.AutoRepairDelay = constants.AutoRepairInterval
}

func (s *Spaceship) AutoEnergyRecharge() {
	s.Energy += constants.AutoEnergyRechargeAmount

	if (s.Energy > s.MaxEnergy) {
		s.Energy = s.MaxEnergy
	}
	s.AutoEnergyRechargeDelay = constants.AutoEnergyRechargeInterval
}

func (s *Spaceship) LeftTurn() {
	s.AngularSpeed += s.angularSpeedDelta()
	s.LimitAngularSpeed()
}

func (s *Spaceship) RightTurn() {
	s.AngularSpeed -= s.angularSpeedDelta()
	s.LimitAngularSpeed()
}

func (s *Spaceship) TurnToTarget() {
	targetAngle := *s.TargetAngle
	offset := targetAngle - s.Facing.Radians()

	if math.Abs(offset) > math.Pi {
		offset -= math.Copysign(2 * math.Pi, offset)
	}

	s.AngularSpeed = -offset * constants.SpaceshipTurnToAngleP

	s.LimitAngularSpeed()
}

func (s *Spaceship) LimitAngularSpeed() {
	if (math.Abs(s.AngularSpeed) > constants.SpaceshipMaxAngularSpeed) {
		s.AngularSpeed = math.Copysign(constants.SpaceshipMaxAngularSpeed, s.AngularSpeed)
	}
}

func (s *Spaceship) ApplyAngularFriction() {
	s.AngularSpeed *= (1 - constants.SpaceshipAngularFriction)
}

func (s *Spaceship) angularSpeedDelta() float64 {
	nonlinearPart := constants.SpaceshipNonlinearAngularAcceleration * math.Abs(s.AngularSpeed)
	linearPart := constants.SpaceshipLinearAngularAcceleration
	return nonlinearPart + linearPart
}
