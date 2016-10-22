package server

import (
	"superstellar/backend/pb"
)

const LeaderboardLength = 10

type Leaderboard struct {
	ranks []Rank
}

type Rank struct {
	clientId uint32
	score uint32
}

type SortableByScore []Rank

func (a SortableByScore) Len() int           { return len(a) }
func (a SortableByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortableByScore) Less(i, j int) bool { return a[i].score < a[j].score }

// ToProto returns protobuf representation
func (rank *Rank) ToProto() *pb.Rank {
	return &pb.Rank{Id: rank.clientId, Score: rank.score}
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
