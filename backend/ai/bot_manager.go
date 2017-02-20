package ai

import (
	"superstellar/backend/events"
	"superstellar/backend/state"
	"superstellar/backend/utils"
)

type BotManager struct {
	idToBot     map[uint32]Bot
	space       *state.Space
	clientIdSeq *utils.IdSequencer
}

func NewBotManager(space *state.Space, clientIdSeq *utils.IdSequencer) *BotManager {
	return &BotManager{
		idToBot: make(map[uint32]Bot),
		space: space,
		clientIdSeq: clientIdSeq,
	}
}

func (m *BotManager) CreateNewBot() {
	id := m.clientIdSeq.NextId()
	m.space.NewSpaceship(id)
	m.idToBot[id] = NewCleverBot()
}

func (m *BotManager) HandleTimeTick(event *events.TimeTick) {
	for id, bot := range m.idToBot {
		bot.HandleStateUpdate(m.space, m.space.Spaceships[id])
	}
}
