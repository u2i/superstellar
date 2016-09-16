package backend

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// Server struct holds server variables.
type Server struct {
	pattern      string
	space        *Space
	clients      map[uint32]*Client
	addCh        chan *Client
	delCh        chan *Client
	moveCh       chan *Move
	doneCh       chan bool
	errCh        chan error
	updateCh     chan bool
	physicsCh    chan bool
	generateIDCh chan chan uint32
	clientID     uint32
}

// NewServer initializes a new server.
func NewServer(pattern string) *Server {
	space := NewSpace()
	clients := make(map[uint32]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	moveCh := make(chan *Move)
	doneCh := make(chan bool)
	errCh := make(chan error)
	updateCh := make(chan bool)
	physicsCh := make(chan bool)
	generateIDCh := make(chan chan uint32)

	return &Server{
		pattern,
		space,
		clients,
		addCh,
		delCh,
		moveCh,
		doneCh,
		errCh,
		updateCh,
		physicsCh,
		generateIDCh,
		0,
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

// Move sends new move command to the server.
func (s *Server) Move(move *Move) {
	s.moveCh <- move
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
	s.mainGameLoop()
}

func (s *Server) sendSpace() {
	bytes, err := json.Marshal(s.space)
	if err != nil {
		log.Println(err)
		return
	}

	str := string(bytes)

	for _, c := range s.clients {
		c.SendSpace(&str)
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
	ticker := time.NewTicker(20 * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			s.updateCh <- true
		}
	}()
}

func (s *Server) runPhysicsTicker() {
	ticker := time.NewTicker(20 * time.Millisecond)
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

		case <-s.moveCh:
			s.handleMoveCommand()

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

func (s *Server) handleAddNewClient(c *Client) {
	log.Println("Added new client")

	s.clients[c.id] = c
	spaceship := NewSpaceship(NewIntVector(0, 0))
	s.space.AddSpaceship(c.id, spaceship)

	log.Println("Now", len(s.clients), "clients connected.")
}

func (s *Server) handleDelClient(c *Client) {
	log.Println("Delete client")

	s.space.RemoveSpaceship(c.id)
	delete(s.clients, c.id)
}

func (s *Server) handleMoveCommand() {

}

func (s *Server) handleUpdate() {
	s.space.randomUpdate()
	s.sendSpace()
}

func (s *Server) handlePhysicsUpdate() {
	s.space.updatePhysics()
}

func (s *Server) handleGenerateIDCh(ch chan uint32) {
	s.clientID++
	ch <- s.clientID
}
