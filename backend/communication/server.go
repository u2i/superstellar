package communication

import (
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"

	"superstellar/backend/events"
	"superstellar/backend/monitor"
	"superstellar/backend/utils"

	"github.com/gorilla/websocket"
)

// Server struct holds server variables.
type Server struct {
	pattern          string
	clients          map[uint32]*Client
	monitor          *monitor.Monitor
	clientID         uint32
	eventsDispatcher *events.EventDispatcher
	idManager        *utils.IdManager
	userNameRegistry *utils.UserNamesRegistry
	upgrader         *websocket.Upgrader
}

// NewServer initializes a new server.
func NewServer(pattern string, monitor *monitor.Monitor, eventDispatcher *events.EventDispatcher,
	idManager *utils.IdManager, userNameRegistry *utils.UserNamesRegistry) *Server {
	return &Server{
		pattern:          pattern,
		clients:          make(map[uint32]*Client),
		monitor:          monitor,
		clientID:         0,
		eventsDispatcher: eventDispatcher,
		idManager:        idManager,
		userNameRegistry: userNameRegistry,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			}},
	}
}

// Listen runs puts server into listening mode.
func (s *Server) Listen() {
	log.Println("Listening server...")

	handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := NewClient(conn, s)
		s.clients[client.id] = client

		s.eventsDispatcher.FireUserConnected(&events.UserConnected{ClientID: client.id})

		log.Println("Added new client. Now", len(s.clients), "clients connected.")
		client.Listen()
	}

	healthCheckHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	http.HandleFunc(s.pattern, handler)
	http.HandleFunc("/healthz", healthCheckHandler)
}

func (s *Server) SendToAllClients(message proto.Message) {
	bytes := marshalMessage(message)
	for _, c := range s.clients {
		c.SendMessage(bytes)
	}
}

func (s *Server) SendToClient(clientID uint32, message proto.Message) {
	bytes := marshalMessage(message)

	client, ok := s.clients[clientID]
	if ok {
		client.SendMessage(bytes)
	} else {
		log.Printf("Client %d not found\n", clientID)
		return
	}
}

func (s *Server) HandleUserLeft(userLeftEvent *events.UserLeft) {
	delete(s.clients, userLeftEvent.ClientID)
}

func (s *Server) ClientIDs() []uint32 {
	clientIDs := make([]uint32, 0, len(s.clients))
	for k := range s.clients {
		clientIDs = append(clientIDs, k)
	}

	return clientIDs
}

func (s *Server) GetClient(clientId uint32) (*Client, bool) {
	client, ok := s.clients[clientId]
	return client, ok
}
