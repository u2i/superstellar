package leaderboard

import (
	"superstellar/backend/pb"
)

type Leaderboard struct {
	ClientId     uint32
	userScore    uint32
	userPosition uint16
	ranks        []Rank
}

// ToProto returns protobuf representation
func (leaderboard *Leaderboard) ToProto() *pb.Leaderboard {
	ranks := make([]*pb.Rank, 0, len(leaderboard.ranks))
	for _, rank := range leaderboard.ranks {
		ranks = append(ranks, rank.ToProto())
	}

	return &pb.Leaderboard{
		Ranks:        ranks,
		ClientId:     leaderboard.ClientId,
		UserScore:    leaderboard.userScore,
		UserPosition: uint32(leaderboard.userPosition),
	}
}

// ToMessage returns protobuffer Message object with Leaderboard containing ordered Ranks.
func (leaderboard *Leaderboard) ToMessage() *pb.Message {
	return &pb.Message{
		Content: &pb.Message_Leaderboard{
			Leaderboard: leaderboard.ToProto(),
		},
	}
}
