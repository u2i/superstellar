package backend

import "superstellar/backend/pb"

// UserInput struct describes client's spaceship input.
type UserInput struct {
	ClientID  uint32
	Thrust    bool
	Direction Direction
}

// NewUserInput returns new instance of UserInput
func NewUserInput(clientID uint32) *UserInput {
	return &UserInput{
		ClientID:  clientID,
		Thrust:    false,
		Direction: NONE,
	}
}

// UserInputFromProto returns new instance of UserInput basing on proto object.
func UserInputFromProto(protoUserInput *pb.UserInput, clientID uint32) *UserInput {
	return &UserInput{
		ClientID:  clientID,
		Thrust:    protoUserInput.Thrust,
		Direction: Direction(protoUserInput.Direction),
	}
}
