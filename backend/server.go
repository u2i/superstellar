package backend

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// Chat server.
type Server struct {
	pattern   string
	positions map[int]*Position
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	moveCh    chan *Move
	doneCh    chan bool
	errCh     chan error
	updateCh  chan bool
}

// Create new chat server.
func NewServer(pattern string) *Server {
	positions := make(map[int]*Position)
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	moveCh := make(chan *Move)
	doneCh := make(chan bool)
	errCh := make(chan error)
	updateCh := make(chan bool)

	return &Server{
		pattern,
		positions,
		clients,
		addCh,
		delCh,
		moveCh,
		doneCh,
		errCh,
		updateCh,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) Move(move *Move) {
	s.moveCh <- move
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendGameState() {
	positions := make([]Position, 0, len(s.positions))

	for _, value := range s.positions {
		positions = append(positions, *value)
	}

	for _, c := range s.clients {
		c.SendPositions(&GameState{Positions: positions})
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	log.Println("Listening server...")

	// websocket handler
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
	log.Println("Created handler")

	ticker := time.NewTicker(30 * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			s.updateCh <- true
		}
	}()

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			s.positions[c.id] = &Position{X: 400, Y: 300}
			log.Println("Now", len(s.clients), "clients connected.")

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)
			delete(s.positions, c.id)

			// broadcast message for all clients
		case move := <-s.moveCh:
			log.Println("New move:", move)

			clientPosition := s.positions[move.ClientID]
			switch move.Direction {
			case "up":
				clientPosition.Y -= 20
				clientPosition.Angle = 3.14 * 1.5
			case "down":
				clientPosition.Y += 20
				clientPosition.Angle = 3.14 * 0.5
			case "left":
				clientPosition.X -= 20
				clientPosition.Angle = 3.14
			case "right":
				clientPosition.X += 20
				clientPosition.Angle = 0.0
			}

		case <-s.updateCh:
			s.sendGameState()

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
