package leaderboard

import (
	"superstellar/backend/pb"
	"superstellar/backend/state"
	"sort"
)

const LeaderboardLength = 10

type Leaderboard struct {
	ranks []Rank
}

func LeaderboardFromSpace(space *state.Space) *Leaderboard{
	size := len(space.Spaceships)
	ranks := make([]Rank, 0, size)
	for _, stateship := range space.Spaceships {
		ranks = append(ranks, Rank{stateship.ID, stateship.MaxHP})
	}
	sort.Stable(sort.Reverse(SortableByScore(ranks)))
	return &Leaderboard{ranks: ranks}
}

// ToProto returns protobuf representation
func (leaderboard *Leaderboard) ToProto() *pb.Leaderboard {
	ranks := make([]*pb.Rank, 0, len(leaderboard.ranks))
	for _, rank := range leaderboard.ranks {
		ranks = append(ranks, rank.ToProto())
	}

	return &pb.Leaderboard{Ranks: ranks}
}

// ToMessage returns protobuffer Message object with Leaderboard containing ordered Ranks.
func (leaderboard *Leaderboard) ToMessage() *pb.Message {
	return &pb.Message{
		Content: &pb.Message_Leaderboard{
			Leaderboard: leaderboard.ToProto(),
		},
	}
}
