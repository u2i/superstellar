package leaderboard_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLeaderboardTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LeaderboardTest Suite")
}
