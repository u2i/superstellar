package server

import (
	"superstellar/backend/state"
	"superstellar/backend/events"
	"superstellar/backend/leaderboard"
	"superstellar/backend/pb"
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
		sender.sendLeaderboard()
	}
}

func (sender *Sender) HandleProjectileFired(projectileFiredEvent *events.ProjectileFired) {
	sender.server.SendToAllClients(projectileFiredEvent.Projectile.ToMessage())
}


func (sender *Sender) HandleUserJoin(userJoinedEvent *events.UserJoined) {
	message := &pb.Message{
		Content: &pb.Message_PlayerJoined{
			PlayerJoined: &pb.PlayerJoined{
				Id:       userJoinedEvent.ClientID,
				Username: userJoinedEvent.UserName,
			},
		},
	}

	sender.server.SendToAllClients(message)
}

func (sender *Sender) HandleUserLeft(userLeftEvent *events.UserLeft) {
	message := &pb.Message{
		Content: &pb.Message_PlayerLeft{
			PlayerLeft: &pb.PlayerLeft{Id: userLeftEvent.ClientID},
		},
	}

	sender.server.SendToAllClients(message)
}


func (sender *Sender) sendSpace() {
	sender.server.SendToAllClients(sender.space.ToMessage())
}

func (sender *Sender) sendLeaderboard() {
	leaderboard := leaderboard.LeaderboardFromSpace(sender.space)
	sender.server.SendToAllClients(leaderboard.ToMessage())
}
