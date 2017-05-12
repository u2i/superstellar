package leaderboard

import (
	"superstellar/backend/pb"
)

type Rank struct {
	clientId uint32
	score    uint32
}

type SortableByScore []Rank

func (a SortableByScore) Len() int {
	return len(a)
}

func (a SortableByScore) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a SortableByScore) Less(i, j int) bool {
	if a[i].score == a[j].score {
		return a[i].clientId > a[j].clientId
	}
	return a[i].score < a[j].score
}

// ToProto returns protobuf representation
func (rank *Rank) ToProto() *pb.Rank {
	return &pb.Rank{Id: rank.clientId, Score: rank.score}
}
