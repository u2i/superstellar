package communication

import (
	"superstellar/backend/events"
	"superstellar/backend/pb"
	"superstellar/backend/persistence"
	"superstellar/backend/utils"
)

type ScoreBoardSender struct {
	server                 *Server
	scoreBoardReader       *persistence.ScoreBoardReader
	userNameRegistry       *utils.UserNamesRegistry
	protoScoreBoardMessage *pb.Message
	minScore               uint32
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
	go sender.sendToOne(userJoinedEvent.ClientID)
}

func (sender *ScoreBoardSender) HandleScoreSent(scoreSent *events.ScoreSent) {
	if scoreSent.Score >= sender.minScore {
		go sender.refreshAndSendToAll()
	}
}

func (sender *ScoreBoardSender) refreshScoreBoard() {
	protoScoreBoard := sender.scoreBoardReader.ReadScoreBoard()

	sender.protoScoreBoardMessage = &pb.Message{
		Content: &pb.Message_ScoreBoard{
			ScoreBoard: protoScoreBoard,
		},
	}
}

func (sender *ScoreBoardSender) sendToOne(clientId uint32) {
	if sender.protoScoreBoardMessage == nil {
		sender.refreshScoreBoard()
	}

	sender.server.SendToClient(clientId, sender.protoScoreBoardMessage)
}

func (sender *ScoreBoardSender) refreshAndSendToAll() {
	sender.refreshScoreBoard()
	sender.server.SendToAllClients(sender.protoScoreBoardMessage)
}
