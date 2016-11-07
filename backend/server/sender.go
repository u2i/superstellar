package server

import (
	"superstellar/backend/state"
	"superstellar/backend/events"
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
}

func (sender *Sender) sendSpace() {
	sender.server.SendToAll(sender.space.ToMessage())
}