package events

import "superstellar/backend/pb"

// UserInput struct describes client's spaceship input.
type UserInput struct {
	ClientID  uint32
	UserInput pb.UserInput
}

// NewUserInput returns new instance of UserInput
func NewUserInput(clientID uint32) *UserInput {
	return &UserInput{
		ClientID:  clientID,
		UserInput: pb.UserInput_CENTER,
	}
}

// UserInputFromProto returns new instance of UserInput basing on proto object.
func UserInputFromProto(userAction *pb.UserAction, clientID uint32) *UserInput {
	return &UserInput{
		ClientID:  clientID,
		UserInput: userAction.UserInput,
	}
}
