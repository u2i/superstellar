package backend

import (
	"log"
	"net/http"
	"superstellar/backend/pb"
	"time"

	"github.com/golang/protobuf/proto"

	"golang.org/x/net/websocket"
)

// Server struct holds server variables.
type Server struct {
	pattern      string
	space        *Space
	clients      map[uint32]*Client
	monitor      *Monitor
	addCh        chan *Client
	delCh        chan *Client
	inputCh      chan *UserInput
	doneCh       chan bool
	errCh        chan error
	updateCh     chan bool
	physicsCh    chan bool
	generateIDCh chan chan uint32
	clientID     uint32
}

const (
	// PhysicsTickInterval equals how often the physics is updated
	PhysicsTickInterval = 20 * time.Millisecond

	// BroadcastStateTickInterval equals how often we broadcast state to clients
	BroadcastStateTickInterval = 20 * time.Millisecond
)

// NewServer initializes a new server.
func NewServer(pattern string) *Server {
	return &Server{
		pattern:      pattern,
		space:        NewSpace(),
		clients:      make(map[uint32]*Client),
		monitor:      newMonitor(),
		addCh:        make(chan *Client),
		delCh:        make(chan *Client),
		inputCh:      make(chan *UserInput),
		doneCh:       make(chan bool),
		errCh:        make(chan error),
		updateCh:     make(chan bool),
		physicsCh:    make(chan bool),
		generateIDCh: make(chan chan uint32),
		clientID:     0,
	}
}

// Add sends client add command to the server.
func (s *Server) Add(c *Client) {
	s.addCh <- c
}

// Del sends client delete command to the server.
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

// UserInput sends new move command to the server.
func (s *Server) UserInput(userInput *UserInput) {
	s.inputCh <- userInput
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
	s.runSenderTicker()
	s.runPhysicsTicker()
	s.monitor.run()
	s.mainGameLoop()
}

func (s *Server) sendSpace() {
	bytes, err := proto.Marshal(s.space.toMessage())
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
		s.Add(client)
		client.Listen()
	}

	http.Handle(s.pattern, websocket.Handler(onConnected))
}

func (s *Server) runSenderTicker() {
	ticker := time.NewTicker(BroadcastStateTickInterval)
	go func() {
		for _ = range ticker.C {
			s.updateCh <- true
		}
	}()
}

func (s *Server) runPhysicsTicker() {
	ticker := time.NewTicker(PhysicsTickInterval)
	go func() {
		for _ = range ticker.C {
			s.physicsCh <- true
		}
	}()
}

func (s *Server) mainGameLoop() {
	for {
		select {

		case c := <-s.addCh:
			s.handleAddNewClient(c)

		case c := <-s.delCh:
			s.handleDelClient(c)

		case input := <-s.inputCh:
			s.handleUserInput(input)

		case <-s.updateCh:
			s.handleUpdate()

		case <-s.physicsCh:
			s.handlePhysicsUpdate()

		case ch := <-s.generateIDCh:
			s.handleGenerateIDCh(ch)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func (s *Server) handleAddNewClient(client *Client) {
	log.Println("Added new client")

	s.clients[client.id] = client
	spaceship := NewSpaceship(client.id, NewIntVector(0, 0))
	s.space.AddSpaceship(client.id, spaceship)

	s.sendHelloMessage(client)

	log.Println("Now", len(s.clients), "clients connected.")
}

func (s *Server) sendHelloMessage(client *Client) {
	message := &pb.Message{
		Content: &pb.Message_Hello{
			Hello: &pb.Hello{MyId: client.id},
		},
	}

	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	client.SendMessage(&bytes)
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

func (s *Server) handleUserInput(userInput *UserInput) {
	s.space.UpdateUserInput(userInput)
}

func (s *Server) handleUpdate() {
	s.sendSpace()
}

func (s *Server) handlePhysicsUpdate() {
	before := time.Now()

	s.space.updatePhysics()

	elapsed := time.Since(before)
	s.monitor.addPhysicsTime(elapsed)
}

func (s *Server) handleGenerateIDCh(ch chan uint32) {
	s.clientID++
	ch <- s.clientID
}
