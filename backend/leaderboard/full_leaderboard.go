package leaderboard

import (
	"sort"
	"superstellar/backend/constants"
	"superstellar/backend/state"
	. "superstellar/math"
)

type FullLeaderboard struct {
	ranks []Rank
}

func FullLeaderboardFromSpace(space *state.Space, clients []uint32) *FullLeaderboard {
	size := len(space.Spaceships)
	ranks := make([]Rank, 0, size)
	for _, client := range clients {
		spaceship, ok := space.Spaceships[client]
		if ok {
			ranks = append(ranks, Rank{client, spaceship.MaxHP})
		}
	}
	sort.Stable(sort.Reverse(SortableByScore(ranks)))
	return &FullLeaderboard{ranks: ranks}
}

func (fullLeaderboard *FullLeaderboard) BuildLeaderboard(userRank Rank, userPosition uint16) *Leaderboard {
	size := Min(len(fullLeaderboard.ranks), constants.LeaderboardLength)

	leaderboard := Leaderboard{
		ClientId:     userRank.clientId,
		userPosition: userPosition,
		userScore:    userRank.score,
		ranks:        fullLeaderboard.ranks[0:size:size],
	}
	return &leaderboard
}

func (fullLeaderboard *FullLeaderboard) BuildLeaderboards() []*Leaderboard {
	size := len(fullLeaderboard.ranks)
	leaderboards := make([]*Leaderboard, size, size)

	for i, userRank := range fullLeaderboard.ranks {
		leaderboard := fullLeaderboard.BuildLeaderboard(userRank, uint16(i+1))
		leaderboards[i] = leaderboard
	}

	return leaderboards
}
