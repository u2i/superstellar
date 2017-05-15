package communication

import (
	"superstellar/backend/constants"
	"superstellar/backend/events"
	"superstellar/backend/leaderboard"
	"superstellar/backend/pb"
	"superstellar/backend/state"
	"superstellar/backend/utils"
)

type Sender struct {
	server             *Server
	space              *state.Space
	userNameRegistry   *utils.UserNamesRegistry
	leaderboardCounter int32
}

func NewSender(server *Server, space *state.Space, userNameRegistry *utils.UserNamesRegistry) *Sender {
	return &Sender{
		server:             server,
		space:              space,
		userNameRegistry:   userNameRegistry,
		leaderboardCounter: 0,
	}
}

func (sender *Sender) HandlePhysicsReady(physicsReadyEvent *events.PhysicsReady) {
	sender.sendSpace()
	sender.leaderboardCounter++

	if sender.leaderboardCounter%10 == 0 {
		sender.leaderboardCounter = 0
		sender.sendLeaderboard()
	}
}

func (sender *Sender) HandleProjectileFired(projectileFiredEvent *events.ProjectileFired) {
	sender.server.SendToAllClients(projectileFiredEvent.Projectile.ToMessage())
}

func (sender *Sender) HandleProjectileHit(projectileHitEvent *events.ProjectileHit) {
	sender.server.SendToAllClients(projectileHitEvent.Projectile.ToHitMessage())
}

func (sender *Sender) HandleUserConnected(userConnectedEvent *events.UserConnected) {
	sender.sendConstantMessage(userConnectedEvent.ClientID)
}

func (sender *Sender) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	sender.sendHelloMessage(userJoinedEvent.ClientID)
	sender.sendUserJoinedMessage(userJoinedEvent.ClientID, userJoinedEvent.UserName)
	sender.server.SendToClient(userJoinedEvent.ClientID, sender.space.ToMessage(true))
}

func (sender *Sender) HandleUserLeft(userLeftEvent *events.UserLeft) {
	message := &pb.Message{
		Content: &pb.Message_PlayerLeft{
			PlayerLeft: &pb.PlayerLeft{Id: userLeftEvent.ClientID},
		},
	}

	sender.server.SendToAllClients(message)
}

func (sender *Sender) HandleObjectDestroyed(userDiedEvent *events.ObjectDestroyed) {
	sender.sendObjectDestroyed(userDiedEvent.DestroyedObject.Id(), userDiedEvent.DestroyedBy.Id())
}

func (sender *Sender) sendSpace() {
	message := sender.space.ToMessage(false)
	if len(message.GetSpace().GetSpaceships()) > 0 || len(message.GetSpace().GetAsteroids()) > 0 {
		sender.server.SendToAllClients(message)
	}
}

func (sender *Sender) sendLeaderboard() {
	fullLeaderboard := leaderboard.FullLeaderboardFromSpace(sender.space, sender.server.ClientIDs())
	leaderboards := fullLeaderboard.BuildLeaderboards()

	for _, l := range leaderboards {
		sender.server.SendToClient(l.ClientId, l.ToMessage())
	}
}

func (sender *Sender) sendUserJoinedMessage(clientID uint32, userName string) {
	message := &pb.Message{
		Content: &pb.Message_PlayerJoined{
			PlayerJoined: &pb.PlayerJoined{
				Id:       clientID,
				Username: userName,
			},
		},
	}

	sender.server.SendToAllClients(message)
}

func (sender *Sender) sendConstantMessage(clientID uint32) {
	constantsProto := &pb.Constants{
		WorldRadius:                           constants.WorldRadius,
		BoundaryAnnulusWidth:                  constants.BoundaryAnnulusWidth,
		SpaceshipAcceleration:                 constants.SpaceshipAcceleration,
		FrictionCoefficient:                   constants.FrictionCoefficient,
		SpaceshipNonlinearAngularAcceleration: constants.SpaceshipNonlinearAngularAcceleration,
		SpaceshipLinearAngularAcceleration:    constants.SpaceshipLinearAngularAcceleration,
		SpaceshipMaxAngularVelocity:           constants.SpaceshipMaxAngularVelocity,
		SpaceshipAngularFriction:              constants.SpaceshipAngularFriction,
		SpaceshipMaxSpeed:                     constants.SpaceshipMaxSpeed,
		SpaceshipBoostFactor:                  constants.SpaceshipBoostFactor,
		AutoRepairInterval:                    constants.AutoRepairInterval,
		AutoRepairAmount:                      constants.AutoRepairAmount,
		AutoEnergyRechargeAmount:              constants.AutoEnergyRechargeAmount,
		BoostPerFrameEnergyCost:               constants.BoostPerFrameEnergyCost,
	}

	message := &pb.Message{
		Content: &pb.Message_Constants{
			Constants: constantsProto,
		},
	}

	sender.server.SendToClient(clientID, message)
}

func (sender *Sender) sendHelloMessage(clientID uint32) {
	idToUsername := make(map[uint32]string)

	for id := range sender.space.Spaceships {
		idToUsername[id] = sender.userNameRegistry.GetUserName(id)
	}

	message := &pb.Message{
		Content: &pb.Message_Hello{
			Hello: &pb.Hello{
				MyId:         clientID,
				IdToUsername: idToUsername,
			},
		},
	}

	sender.server.SendToClient(clientID, message)
}

func (sender *Sender) sendObjectDestroyed(clientID uint32, killedBy uint32) {
	message := &pb.Message{
		Content: &pb.Message_PlayerDied{
			PlayerDied: &pb.PlayerDied{
				Id:       clientID,
				KilledBy: killedBy,
			},
		},
	}

	sender.server.SendToAllClients(message)
}
