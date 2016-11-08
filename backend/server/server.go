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
	"superstellar/backend/events"
	"sync/atomic"
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
		clientID:     0,
		eventsDispatcher: eventDispatcher,
	}
}

// Done sends done command to the server.
func (s *Server) Done() {
	s.doneCh <- true
}

// Err sends error to the server.
func (s *Server) Err(err error) {
	s.errCh <- err
}

// Listen runs puts server into listening mode.
func (s *Server) Listen() {
	log.Println("Listening server...")

	s.addNewClientHandler()
	s.monitor.run()
	s.eventsDispatcher.RegisterTimeTickListener(s)
}

func (s *Server) SendToAllClients(message proto.Message) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	for _, c := range s.clients {
		c.SendMessage(&bytes)
	}
}

func (s *Server) SendToClient(clientID uint32, message proto.Message) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	client, ok := s.clients[clientID]
	if ok {
		client.SendMessage(&bytes)
	} else {
		log.Println("Client %i not found", clientID)
		return
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

		client := NewClient(ws, s, s.nextClientID())
		s.handleAddNewClient(client)
		client.Listen()
	}

	http.Handle(s.pattern, websocket.Handler(onConnected))
}

func (s *Server) handleAddNewClient(client *Client) {
	log.Println("Added new client")

	s.clients[client.id] = client
	s.sendHelloMessage(client)

	log.Println("Now", len(s.clients), "clients connected.")
}

func (s *Server) deleteClient(c *Client) {
	log.Println("Delete client")

	delete(s.clients, c.id)
}

func (s *Server) nextClientID() uint32 {
	return atomic.AddUint32(&s.clientID, 1)
}







func (s *Server) readChannels() {
	select {

	case err := <-s.errCh:
		log.Println("Error:", err.Error())

	case <-s.doneCh:
		return
	default:
	}
}

func (s *Server) HandleTimeTick(e *events.TimeTick) {
	s.readChannels()
}


func (s *Server) SendJoinGameAckMessage(client *Client, joinGameAck *pb.JoinGameAck) {
	message := &pb.Message{
		Content: &pb.Message_JoinGameAck{
			JoinGameAck: joinGameAck,
		},
	}

	s.SendToClient(client.id, message)
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

	s.SendToAllClients(message)
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

	s.SendToClient(client.id, message)
}


func (s *Server) handlePhysicsUpdate() {
	before := time.Now()

	simulation.UpdatePhysics(s.space, s.eventsDispatcher)

	elapsed := time.Since(before)
	s.monitor.addPhysicsTime(elapsed)
}