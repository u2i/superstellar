package state

import (
	"fmt"
	"math"
	"superstellar/backend/constants"
	"superstellar/backend/pb"
	"superstellar/backend/types"
	"superstellar/backend/utils"
	"time"
)

// Direction is a type describing user input on spaceship rotation.
type Direction int

// Constants describing user input on spaceship rotation.
const (
	NONE  Direction = iota
	RIGHT
	LEFT
)

// Spaceship struct describes a spaceship.
type Spaceship struct {
	ObjectState
	InputThrust          bool
	InputBoost           bool
	InputDirection       Direction
	TargetAngle          *float64
	Fire                 bool
	LastShotTime         time.Time
	HP                   uint32
	MaxHP                uint32
	Energy               uint32
	MaxEnergy            uint32
	AutoRepairDelay      uint32
}

func NewSpaceship(clientId uint32, initialPosition *types.Point) *Spaceship {
	objectState := NewObjectState(clientId, initialPosition, types.ZeroVector())

	return &Spaceship{
		ObjectState:               *objectState,
		InputThrust:               false,
		InputDirection:            NONE,
		Fire:                      false,
		LastShotTime:              time.Now(),
		HP:                        constants.SpaceshipInitialHP,
		MaxHP:                     constants.SpaceshipInitialHP,
		Energy:                    constants.SpaceshipInitialEnergy,
		MaxEnergy:                 constants.SpaceshipInitialEnergy,
		AutoRepairDelay:           constants.AutoRepairDelay,
	}
}

// String function returns string representation.
func (s *Spaceship) String() string {
	return fmt.Sprintf("(%v, %v, %v)", s.Position, s.Velocity, s.Facing)
}

func (s *Spaceship) UpdateUserInput(userInput pb.UserInput) {
	switch userInput {
	case pb.UserInput_CENTER:
		s.InputDirection = NONE
		s.TargetAngle = nil
		s.MarkDirty()
	case pb.UserInput_LEFT:
		s.InputDirection = LEFT
		s.TargetAngle = nil
		s.MarkDirty()
	case pb.UserInput_RIGHT:
		s.InputDirection = RIGHT
		s.TargetAngle = nil
		s.MarkDirty()
	case pb.UserInput_THRUST_ON:
		s.InputThrust = true
		s.MarkDirty()
	case pb.UserInput_THRUST_OFF:
		s.InputThrust = false
		s.MarkDirty()
	case pb.UserInput_FIRE_START:
		s.Fire = true
		s.MarkDirty()
	case pb.UserInput_FIRE_STOP:
		s.Fire = false
	case pb.UserInput_BOOST_ON:
		s.InputBoost = true
		s.MarkDirty()
	case pb.UserInput_BOOST_OFF:
		s.InputBoost = false
		s.MarkDirty()
	}
}

func (s *Spaceship) NotifyAboutNewFrame() {
	s.ObjectState.NotifyAboutNewFrame()

	s.handleAutoEnergyRecharge()
	s.handleAutoRepair()
}

func (s *Spaceship) UpdateTargetAngle(angle float64) {
	s.MarkDirty()
	s.TargetAngle = &angle
	s.InputDirection = NONE
}

// ToProto returns protobuf representation
func (s *Spaceship) ToProto() *pb.Spaceship {
	return &pb.Spaceship{
		Id:              s.Id(),
		Position:        s.Position().ToProto(),
		Velocity:        s.Velocity().ToProto(),
		Facing:          s.Facing(),
		AngularVelocity: s.AngularVelocity(),
		InputDirection:  pb.Direction(s.InputDirection),
		InputThrust:     s.InputThrust,
		InputBoost:      s.InputBoost,
		MaxHp:           s.MaxHP,
		Hp:              s.HP,
		MaxEnergy:       s.MaxEnergy,
		Energy:          s.Energy,
		AutoRepairDelay: s.AutoRepairDelay,
	}
}

func (s *Spaceship) ShootIfPossible() (canShoot bool) {
	if s.Energy >= constants.BasicWeaponEnergyCost {
		canShoot = true
		s.Energy -= constants.BasicWeaponEnergyCost
		s.MarkDirty()
	} else {
		canShoot = false
	}
	return canShoot
}

func (s *Spaceship) BoostIfPossible() (canBoost bool) {
	if s.Energy >= constants.BoostPerFrameEnergyCost {
		canBoost = true
		s.Energy -= constants.BoostPerFrameEnergyCost
	} else {
		canBoost = false
	}
	return canBoost
}

func (s *Spaceship) CollideWithProjectile(projectile *Projectile) {
	if s.HP < constants.ProjectileDamage {
		s.HP = 0
	} else {
		s.HP -= constants.ProjectileDamage
	}
	s.AutoRepairDelay = constants.AutoRepairDelay

	s.MarkDirty()
}

func (s *Spaceship) AddReward(reward uint32) {
	s.HP += reward
	s.MaxHP += reward

	s.MarkDirty()
}

func (s *Spaceship) AddEnergyReward(reward uint32) {
	s.Energy += reward
	s.MaxEnergy += reward

	s.MarkDirty()
}

func (s *Spaceship) handleAutoRepair() {
	if s.AutoRepairDelay == 0 {
		s.HP = utils.Min(s.HP+constants.AutoRepairAmount, s.MaxHP)
	} else {
		s.AutoRepairDelay--
	}
}

func (s *Spaceship) handleAutoEnergyRecharge() {
	s.Energy = utils.Min(s.Energy+constants.AutoEnergyRechargeAmount, s.MaxEnergy)
}

func (s *Spaceship) LeftTurn() {
	s.SetAngularVelocityDelta(s.angularVelocityDelta())
	s.LimitAngularVelocityDelta()
}

func (s *Spaceship) RightTurn() {
	s.SetAngularVelocityDelta(-s.angularVelocityDelta())
	s.LimitAngularVelocityDelta()
}

func (s *Spaceship) TurnToTarget() {
	targetAngle := *s.TargetAngle
	offset := targetAngle - s.Facing()

	if math.Abs(offset) > math.Pi {
		offset -= math.Copysign(2*math.Pi, offset)
	}

	targetAngularVelocity := -offset * constants.SpaceshipTurnToAngleP
	s.SetAngularVelocityDelta(targetAngularVelocity - s.AngularVelocity())

	s.LimitAngularVelocityDelta()
}

func (s *Spaceship) LimitAngularVelocityDelta() {
	potentialAngularVelocity := s.AngularVelocity() + s.AngularVelocityDelta()
	diff := math.Abs(potentialAngularVelocity) - constants.SpaceshipMaxAngularVelocity

	if diff > 0 {
		s.SetAngularVelocityDelta(s.AngularVelocityDelta() - math.Copysign(diff, s.AngularVelocity()))
	}
}

func (s *Spaceship) ApplyAngularFriction() {
	s.SetAngularVelocity(s.AngularVelocity() * (1 - constants.SpaceshipAngularFriction))
}

func (s *Spaceship) angularVelocityDelta() float64 {
	nonlinearPart := constants.SpaceshipNonlinearAngularAcceleration * math.Abs(s.AngularVelocity())
	linearPart := constants.SpaceshipLinearAngularAcceleration
	return nonlinearPart + linearPart
}
