package communication

import (
	"superstellar/backend/state"
	"superstellar/backend/events"
	"superstellar/backend/leaderboard"
	"superstellar/backend/pb"
	"superstellar/backend/constants"
		"time"
)

type Sender struct {
	server *Server
	space  *state.Space
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

func (sender *Sender) HandleProjectileHit(projectileHitEvent *events.ProjectileHit) {
	sender.server.SendToAllClients(projectileHitEvent.Projectile.ToHitMessage())
}

func (sender *Sender) HandleUserJoined(userJoinedEvent *events.UserJoined) {
	sender.sendHelloMessage(userJoinedEvent.ClientID)
	sender.sendUserJoinedMessage(userJoinedEvent.ClientID, userJoinedEvent.UserName)
}

func (sender *Sender) HandleUserLeft(userLeftEvent *events.UserLeft) {
	message := &pb.Message{
		Content: &pb.Message_PlayerLeft{
			PlayerLeft: &pb.PlayerLeft{Id: userLeftEvent.ClientID},
		},
	}

	sender.server.SendToAllClients(message)
}

func (sender *Sender) HandleUserDied(userDiedEvent *events.UserDied) {
	sender.sendUserDiedMessage(userDiedEvent.ClientID, userDiedEvent.KilledBy)
}

func (sender *Sender) sendSpace() {
	sender.server.SendToAllClients(sender.space.ToMessage())
}

func (sender *Sender) sendLeaderboard() {
	fullLeaderboard := leaderboard.FullLeaderboardFromSpace(sender.space)
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

func (sender *Sender) sendHelloMessage(clientID uint32) {
	idToUsername := make(map[uint32]string)

	for id, client := range sender.server.clients {
		idToUsername[id] = client.username
	}

	message := &pb.Message{
		Content: &pb.Message_Hello{
			Hello: &pb.Hello{
				MyId:         clientID,
				IdToUsername: idToUsername,
				WorldRadius:  constants.WorldRadius / 100,
				BoundaryAnnulusWidth: constants.BoundaryAnnulusWidth / 100,
				FirstPhysicsFrameTimestamp: uint64(sender.space.FirstPhysicsFrameTimestamp),
				PhysicsFrameRate: uint32(time.Second.Nanoseconds() / constants.PhysicsFrameDuration.Nanoseconds()),
			},
		},
	}

	sender.server.SendToClient(clientID, message)
}

func (sender *Sender) sendUserDiedMessage(clientID uint32, killedBy uint32) {
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
