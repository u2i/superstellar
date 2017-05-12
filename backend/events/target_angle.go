package events

import "superstellar/backend/pb"

// TargetAngle struct describes client's spaceship angle steering input.
type TargetAngle struct {
	ClientID uint32
	Angle    float64
}

// NewTargetAngle returns new instance of TargetAngle.
func NewTargetAngle(clientID uint32) *TargetAngle {
	return &TargetAngle{
		ClientID: clientID,
		Angle:    0.0,
	}
}

// TargetAngleFromProto returns new instance of TargetAngle basing on proto object.
func TargetAngleFromProto(targetAngle *pb.TargetAngle, clientID uint32) *TargetAngle {
	return &TargetAngle{
		ClientID: clientID,
		Angle:    float64(targetAngle.Angle),
	}
}
