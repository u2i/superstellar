package communication

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"superstellar/backend/constants"
	"superstellar/backend/events"
	"superstellar/backend/pb"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/gorilla/websocket"
)

const channelBufSize = 100

// Client struct holds client-specific variables.
type Client struct {
	id     uint32
	ws     *websocket.Conn
	ch     chan *[]byte
	doneCh chan bool
	server *Server
}

// NewClient initializes a new Client struct with given websocket.
func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	ch := make(chan *[]byte, channelBufSize)
	doneCh := make(chan bool)

	return &Client{server.idManager.NextPlayerId(), ws, ch, doneCh, server}
}

// Conn returns client's websocket.Conn struct.
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

// SendMessage sends game state to the client.
func (c *Client) SendMessage(bytes *[]byte) {
	select {
	case c.ch <- bytes:
	default:
		c.server.monitor.AddDroppedMessage()
	}
}

// Done sends done message to the Client which closes the conection.
func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	defer func() {
		err := c.ws.Close()
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()

	log.Println("Listening write to client")
	for {
		select {

		case bytes := <-c.ch:
			before := time.Now()
			err := c.ws.WriteMessage(websocket.BinaryMessage, *bytes)
			after := time.Now()

			if err != nil {
				log.Println(err)
			} else {
				elapsed := after.Sub(before)
				c.server.monitor.AddSendTime(elapsed)
			}

		case <-c.doneCh:
			c.doneCh <- true
			return
		}
	}
}

func (c *Client) listenRead() {
	defer func() {
		err := c.ws.Close()
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()

	log.Println("Listening read from client")
	for {
		select {

		case <-c.doneCh:
			c.doneCh <- true
			return

		default:
			c.readFromWebSocket()
		}
	}
}

func (c *Client) readFromWebSocket() {
	messageType, data, err := c.ws.ReadMessage()
	if err != nil {
		log.Println(err)

		c.doneCh <- true
		c.server.eventsDispatcher.FireUserLeft(&events.UserLeft{ClientID: c.id})
	} else if messageType != websocket.BinaryMessage {
		log.Println("Non binary message recived, ignoring")
	} else {
		c.unmarshalUserInput(data)
	}
}

func (c *Client) unmarshalUserInput(data []byte) {
	protoUserMessage := &pb.UserMessage{}
	if err := proto.Unmarshal(data, protoUserMessage); err != nil {
		log.Fatalln("Failed to unmarshal UserInput:", err)
		return
	}

	switch x := protoUserMessage.Content.(type) {
	case *pb.UserMessage_UserAction:
		userInputEvent := events.UserInputFromProto(protoUserMessage.GetUserAction(), c.id)
		c.server.eventsDispatcher.FireUserInput(userInputEvent)
	case *pb.UserMessage_TargetAngle:
		targetAngleEvent := events.TargetAngleFromProto(protoUserMessage.GetTargetAngle(), c.id)
		c.server.eventsDispatcher.FireTargetAngle(targetAngleEvent)
	case *pb.UserMessage_JoinGame:
		c.tryToJoinGame(protoUserMessage.GetJoinGame())
	case *pb.UserMessage_Ping:
		c.sendPong(protoUserMessage.GetPing().Id)
	default:
		log.Fatalln("Unknown message type %T", x)
	}
}

func (c *Client) tryToJoinGame(joinGameMsg *pb.JoinGame) {
	username := strings.TrimSpace(joinGameMsg.Username)
	ok, err := c.validateUser(username)

	if ok {
		c.server.userNameRegistry.AddUserName(c.id, username)
		c.sendJoinGameAckMessage(&pb.JoinGameAck{Success: true})
		c.server.eventsDispatcher.FireUserJoined(&events.UserJoined{ClientID: c.id, UserName: username})
	} else {
		c.sendJoinGameAckMessage(
			&pb.JoinGameAck{Success: false, Error: err.Error()},
		)
	}
}

func (c *Client) validateUser(username string) (bool, error) {
	if len(c.server.clients) > constants.MaxNumberOfClients {
		return false, errors.New("We reached the maximum number of users, Captain.")
	}

	length := len(username)

	if length < constants.MinimumUsernameLength {
		return false, errors.New(fmt.Sprintf("I doubt your name is shorter than %d characters, Captain.", constants.MinimumUsernameLength))
	}

	if length > constants.MaximumUsernameLength {
		return false, errors.New(fmt.Sprintf("Space fleet doesn't allow names longer than %d characters!", constants.MaximumUsernameLength))
	}

	return true, nil
}

func (c *Client) sendJoinGameAckMessage(joinGameAck *pb.JoinGameAck) {
	message := &pb.Message{
		Content: &pb.Message_JoinGameAck{
			JoinGameAck: joinGameAck,
		},
	}

	c.SendMessage(marshalMessage(message))
}

func (c *Client) sendPong(id uint32) {
	message := &pb.Message{
		Content: &pb.Message_Pong{
			Pong: &pb.Pong{Id: id},
		},
	}

	c.SendMessage(marshalMessage(message))
}
