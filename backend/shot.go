package backend

import "superstellar/backend/pb"

// Shot struct holds players' shots data.
type Shot struct {
	ClientID uint32
	FrameID  uint32
	Origin   *IntVector
	Facing   *Vector
	Range    uint32
	Position *IntVector
}

// NewShot returns new instance of Shot
func NewShot(clientID, frameID uint32, origin *IntVector, facing *Vector,
	shotRange uint32) *Shot {
	return &Shot{
		ClientID: clientID,
		FrameID:  frameID,
		Origin:   origin,
		Facing:   facing,
		Range:    shotRange,
		Position: origin,
	}
}

func (shot *Shot) toProto() *pb.Shot {
	return &pb.Shot{
		FrameId: shot.FrameID,
		Origin:  shot.Origin.toProto(),
		Facing:  float32(shot.Facing.Radians()),
		Range:   shot.Range,
	}
}
