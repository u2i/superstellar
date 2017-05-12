package leaderboard

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"superstellar/backend/constants"
)

var _ = Describe("FullLeaderboard", func() {
	Describe("BuildLeaderboard", func() {
		var ranks []Rank
		var (
			fullLeaderboard FullLeaderboard
			userRank        Rank
		)

		JustBeforeEach(func() {
			fullLeaderboard = FullLeaderboard{ranks: ranks}
		})

		Context("when leaderboard is not empty", func() {
			BeforeEach(func() {
				ranks = []Rank{{clientId: 2, score: 100}}
			})

			Context("when user is on a list", func() {
				BeforeEach(func() {
					userRank = Rank{clientId: 2, score: 300}
				})

				It("creates leaderboard with no empty ranks", func() {
					subject := fullLeaderboard.BuildLeaderboard(userRank, 1)
					Expect(subject.ranks).To(Equal(ranks))
					Expect(subject.ClientId).To(Equal(userRank.clientId))
					Expect(subject.userScore).To(Equal(userRank.score))
					Expect(subject.userPosition).To(Equal(uint16(1)))
				})
			})
		})

		Context("when leaderboard size greater than 10", func() {
			BeforeEach(func() {
				ranks = []Rank{
					{clientId: 22, score: 990},
					{clientId: 23, score: 950},
					{clientId: 24, score: 900},
					{clientId: 25, score: 800},
					{clientId: 26, score: 700},
					{clientId: 27, score: 600},
					{clientId: 28, score: 500},
					{clientId: 29, score: 400},
					{clientId: 30, score: 300},
					{clientId: 31, score: 200},
					{clientId: 32, score: 100},
					{clientId: 33, score: 100},
				}
				userRank = Rank{clientId: 2, score: 300}
			})

			It("creates leaderboard with no empty ranks", func() {
				subject := fullLeaderboard.BuildLeaderboard(userRank, 1)
				Expect(subject.ClientId).To(Equal(userRank.clientId))
				Expect(subject.userScore).To(Equal(userRank.score))
				Expect(subject.userPosition).To(Equal(uint16(1)))
				Expect(subject.ranks).To(HaveLen(constants.LeaderboardLength))
			})
		})
	})

	Describe("BuildLeaderboards", func() {
		var ranks []Rank
		var fullLeaderboard FullLeaderboard

		empty := make([]Rank, 0, 0)

		JustBeforeEach(func() {
			fullLeaderboard = FullLeaderboard{ranks: ranks}
		})

		Context("when leaderboard is empty", func() {
			BeforeEach(func() {
				ranks = empty
			})

			It("creates no leaderboards", func() {
				subject := fullLeaderboard.BuildLeaderboards()
				Expect(subject).To(BeEmpty())
			})
		})

		Context("when leaderboard is not empty", func() {
			BeforeEach(func() {
				ranks = []Rank{
					{clientId: 2, score: 130},
					{clientId: 3, score: 100},
				}
			})

			It("creates leaderboard for each rank", func() {
				subject := fullLeaderboard.BuildLeaderboards()
				Expect(subject).To(HaveLen(2))
				Expect(subject).To(HaveCap(2))
				Expect(subject[0].ClientId).To(Equal(uint32(2)))
				Expect(subject[0].userPosition).To(Equal(uint16(1)))
				Expect(subject[0].userScore).To(Equal(uint32(130)))
				Expect(subject[0].ranks).To(HaveLen(2))
				Expect(subject[0].ranks).To(HaveCap(2))
				Expect(subject[1].ClientId).To(Equal(uint32(3)))
				Expect(subject[1].userPosition).To(Equal(uint16(2)))
				Expect(subject[1].userScore).To(Equal(uint32(100)))
				Expect(subject[1].ranks).To(HaveLen(2))
				Expect(subject[1].ranks).To(HaveCap(2))
			})
		})

		Context("when leaderboard limits to top ten", func() {
			BeforeEach(func() {
				ranks = []Rank{
					{clientId: 22, score: 990},
					{clientId: 23, score: 950},
					{clientId: 24, score: 900},
					{clientId: 25, score: 800},
					{clientId: 26, score: 700},
					{clientId: 27, score: 600},
					{clientId: 28, score: 500},
					{clientId: 29, score: 400},
					{clientId: 30, score: 300},
					{clientId: 31, score: 200},
					{clientId: 32, score: 100},
					{clientId: 33, score: 100},
				}
			})

			It("creates leaderboard for each rank", func() {
				subject := fullLeaderboard.BuildLeaderboards()
				Expect(subject).To(HaveLen(len(ranks)))
				Expect(subject).To(HaveCap(len(ranks)))

				var ids = []int{22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

				for i, id := range ids {
					Expect(subject[i].ClientId).To(Equal(uint32(id)))
					Expect(subject[i].userPosition).To(Equal(uint16(i + 1)))
					Expect(subject[i].ranks).To(HaveLen(constants.LeaderboardLength))
					Expect(subject[i].ranks).To(HaveCap(constants.LeaderboardLength))
				}
			})
		})

	})
})
