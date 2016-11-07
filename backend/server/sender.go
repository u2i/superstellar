package server

import (
	"superstellar/backend/state"
	"superstellar/backend/events"
	"superstellar/backend/leaderboard"
)

type Sender struct {
	server 	*Server
	space 	*state.Space
}

func NewSender(server *Server, space *state.Space) *Sender {
	return &Sender{
		server: server,
		space: space,
	}
}

func (sender *Sender) HandleTimeTick(timeTickEvent *events.TimeTick) {
	sender.sendSpace()
	if (timeTickEvent.FrameId % 50 == 0) {
		sender.handleLeaderboardUpdate()
	}
}

func (sender *Sender) sendSpace() {
	sender.server.SendToAll(sender.space.ToMessage())
}

func (sender *Sender) handleLeaderboardUpdate() {
	leaderboard := leaderboard.LeaderboardFromSpace(sender.space)
	sender.server.SendToAll(leaderboard.ToMessage())
}