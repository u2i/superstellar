package communication

import (
	"superstellar/backend/events"
	"superstellar/backend/pb"
	"superstellar/backend/persistence"
	"superstellar/backend/utils"
)

type ScoreBoardSender struct {
	server           *Server
	scoreBoardReader *persistence.ScoreBoardReader
	userNameRegistry *utils.UserNamesRegistry
}

func NewScoreBoardSender(server *Server, scoreBoardReader *persistence.ScoreBoardReader,
	userNameRegistry *utils.UserNamesRegistry) *ScoreBoardSender {
	return &ScoreBoardSender{
		server:           server,
		scoreBoardReader: scoreBoardReader,
		userNameRegistry: userNameRegistry,
	}
}

func (sender *ScoreBoardSender) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	go sender.sendScoreBoard()
}

func (sender *ScoreBoardSender) sendScoreBoard() {
	protoScoreBoard := sender.scoreBoardReader.ReadScoreBoard()
	message := &pb.Message{
		Content: &pb.Message_ScoreBoard{
			ScoreBoard: protoScoreBoard,
		},
	}

	sender.server.SendToAllClients(message)
}
