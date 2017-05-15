package ai

import (
	"fmt"
	"superstellar/backend/events"
	"superstellar/backend/state"
	"superstellar/backend/utils"
)

type BotManager struct {
	idToBot           map[uint32]Bot
	space             *state.Space
	idManager         *utils.IdManager
	dispatcher        *events.EventDispatcher
	userNamesRegistry *utils.UserNamesRegistry
}

func NewBotManager(dispatcher *events.EventDispatcher, space *state.Space, idManager *utils.IdManager,
	userNamesRegistry *utils.UserNamesRegistry) *BotManager {
	return &BotManager{
		idToBot:           make(map[uint32]Bot),
		space:             space,
		idManager:         idManager,
		dispatcher:        dispatcher,
		userNamesRegistry: userNamesRegistry,
	}
}

func (m *BotManager) CreateBots(numberOfBots int) {
	for i := 0; i < numberOfBots; i++ {
		m.CreateNewBot()
	}
}

func (m *BotManager) CreateNewBot() {
	id := m.idManager.NextBotId()
	userName := fmt.Sprintf("Bot %d", id)
	m.space.NewSpaceship(id)
	m.idToBot[id] = NewCleverBot()
	m.dispatcher.FireUserJoined(&events.UserJoined{ClientID: id, UserName: userName})
	m.userNamesRegistry.AddUserName(id, userName)
}

func (m *BotManager) HandleTimeTick(event *events.TimeTick) {
	for id, bot := range m.idToBot {
		bot.HandleStateUpdate(m.space, m.space.Spaceships[id])
	}
}

func (m *BotManager) HandleObjectDestroyed(event *events.ObjectDestroyed) {
	destroyedId := event.DestroyedObject.Id()
	_, botDied := m.idToBot[destroyedId]

	if botDied {
		delete(m.idToBot, destroyedId)
		m.CreateNewBot()
	}
}
