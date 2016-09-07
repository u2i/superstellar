package backend

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	pattern      string
	space        *Space
	clients      map[int]*Client
	addCh        chan *Client
	delCh        chan *Client
	moveCh       chan *Move
	doneCh       chan bool
	errCh        chan error
	updateCh     chan bool
	nextClientId int
}

func NewServer(pattern string) *Server {
	space := NewSpace()
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	moveCh := make(chan *Move)
	doneCh := make(chan bool)
	errCh := make(chan error)
	updateCh := make(chan bool)
	nextClientId := 0

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
		nextClientId
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

func (s *Server) sendSpace() {
	for _, c := range s.clients {
		c.SendSpace(s.space)
	}
}

func (s *Server) nextClientId() int {
	clientId := s.nextClientId
	s.nextClientId += 1
	return clientId
}

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

	ticker := time.NewTicker(20 * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			s.updateCh <- true
		}
	}()

	s.mainGameLoop()
}

func (s *Server) mainGameLoop() {
	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			clientId := s.nextClientId()
			s.clients[clientId] = c
			spaceship := NewSpaceship(NewVector(400.0, 300.0))
			c.id = clientId
			s.space.AddSpaceship(clientId, spaceship)
			log.Println("Now", len(s.clients), "clients connected.")

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			s.space.RemoveSpaceship(c.id)
			delete(s.clients, c.id)

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
			s.sendSpace()

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
