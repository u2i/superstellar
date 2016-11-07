package server

import (
	"log"
	"net/http"
	"superstellar/backend/pb"
	"superstellar/backend/simulation"
	"superstellar/backend/state"
	"time"

	"github.com/golang/protobuf/proto"

	"golang.org/x/net/websocket"
	"sort"
	"superstellar/backend/events"
)

// Server struct holds server variables.
type Server struct {
	pattern          string
	space            *state.Space
	clients          map[uint32]*Client
	monitor          *Monitor
	delCh            chan *Client
	doneCh           chan bool
	errCh            chan error
	generateIDCh     chan chan uint32
	clientID         uint32
	eventsDispatcher *events.EventDispatcher
}

// NewServer initializes a new server.
// TODO remove space
func NewServer(pattern string, eventDispatcher *events.EventDispatcher, space *state.Space) *Server {
	return &Server{
		pattern:      pattern,
		space:        space,
		clients:      make(map[uint32]*Client),
		monitor:      newMonitor(),
		delCh:        make(chan *Client),
		doneCh:       make(chan bool),
		errCh:        make(chan error),
		generateIDCh: make(chan chan uint32),
		clientID:     0,
		eventsDispatcher: eventDispatcher,
	}
}

// Del sends client delete command to the server.
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

// Done sends done command to the server.
func (s *Server) Done() {
	s.doneCh <- true
}

// Err sends error to the server.
func (s *Server) Err(err error) {
	s.errCh <- err
}

// GenerateID generates new unique ID for a client
func (s *Server) GenerateID() uint32 {
	ch := make(chan uint32)
	s.generateIDCh <- ch
	return <-ch
}

// Listen runs puts server into listening mode.
func (s *Server) Listen() {
	log.Println("Listening server...")

	s.addNewClientHandler()
	s.monitor.run()
	s.eventsDispatcher.RegisterTimeTickListener(s)
	s.eventsDispatcher.RegisterProjectileFiredListener(s)
}

func (s *Server) SendToAll(message proto.Message) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	for _, c := range s.clients {
		c.SendMessage(&bytes)
	}
}

func (s *Server) sendLeaderboard(leaderboard *Leaderboard) {
	bytes, err := proto.Marshal(leaderboard.ToMessage())
	if err != nil {
		log.Println(err)
		return
	}

	for _, c := range s.clients {
		c.SendMessage(&bytes)
	}
}

func (s *Server) addNewClientHandler() {
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.handleAddNewClient(client)
		client.Listen()
	}

	http.Handle(s.pattern, websocket.Handler(onConnected))
}

func (s *Server) readChannels() {
	select {

	case c := <-s.delCh:
		s.handleDelClient(c)

	case ch := <-s.generateIDCh:
		s.handleGenerateIDCh(ch)

	case err := <-s.errCh:
		log.Println("Error:", err.Error())

	case <-s.doneCh:
		return
	default:
	}
}

func (s *Server) HandleTimeTick(e *events.TimeTick) {
	s.readChannels()
	if (e.FrameId % 50 == 0) {
		s.handleLeaderboardUpdate()
	}
}

func (s *Server) HandleProjectileFired(e *events.ProjectileFired) {
	s.sendShot(e.Projectile)
}

func (s *Server) handleAddNewClient(client *Client) {
	log.Println("Added new client")

	s.clients[client.id] = client
	s.sendHelloMessage(client)

	log.Println("Now", len(s.clients), "clients connected.")
}

func (s *Server) SendJoinGameAckMessage(client *Client, joinGameAck *pb.JoinGameAck) {
	message := &pb.Message{
		Content: &pb.Message_JoinGameAck{
			JoinGameAck: joinGameAck,
		},
	}

	bytes, err := proto.Marshal(message)

	if err != nil {
		log.Println(err)
		return
	}

	client.SendMessage(&bytes)
}

func (s *Server) SendPlayerJoinedMessage(client *Client) {
	message := &pb.Message{
		Content: &pb.Message_PlayerJoined{
			PlayerJoined: &pb.PlayerJoined{
				Id:       client.id,
				Username: client.username,
			},
		},
	}

	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	for _, c := range s.clients {
		c.SendMessage(&bytes)
	}
}

func (s *Server) sendHelloMessage(client *Client) {
	idToUsername := make(map[uint32]string)

	for id, client := range s.clients {
		idToUsername[id] = client.username
	}

	message := &pb.Message{
		Content: &pb.Message_Hello{
			Hello: &pb.Hello{
				MyId:         client.id,
				IdToUsername: idToUsername,
			},
		},
	}

	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	client.SendMessage(&bytes)
}

func (s *Server) sendShot(shot *state.Projectile) {
	message := &pb.Message{
		Content: &pb.Message_ProjectileFired{
			ProjectileFired: shot.ToProto(),
		},
	}

	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	for _, c := range s.clients {
		c.SendMessage(&bytes)
	}
}

func (s *Server) handleDelClient(c *Client) {
	log.Println("Delete client")

	s.space.RemoveSpaceship(c.id)

	delete(s.clients, c.id)

	s.sendUserLeftMessage(c.id)
}

func (s *Server) sendUserLeftMessage(userID uint32) {
	message := &pb.Message{
		Content: &pb.Message_PlayerLeft{
			PlayerLeft: &pb.PlayerLeft{Id: userID},
		},
	}

	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	for _, c := range s.clients {
		c.SendMessage(&bytes)
	}
}

func (s *Server) handlePhysicsUpdate() {
	before := time.Now()

	simulation.UpdatePhysics(s.space, s.eventsDispatcher)

	elapsed := time.Since(before)
	s.monitor.addPhysicsTime(elapsed)
}

func (s *Server) handleLeaderboardUpdate() {
	size := len(s.space.Spaceships)
	ranks := make([]Rank, 0, size)
	for _, stateship := range s.space.Spaceships {
		// TODO: change to MaxHP?
		ranks = append(ranks, Rank{stateship.ID, stateship.HP})
	}
	sort.Stable(sort.Reverse(SortableByScore(ranks)))

	s.sendLeaderboard(&Leaderboard{ranks})
}

func (s *Server) handleGenerateIDCh(ch chan uint32) {
	s.clientID++
	ch <- s.clientID
}
