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

	"superstellar/backend/monitor"

	"superstellar/backend/utils"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

// Client struct holds client-specific variables.
type Client struct {
	id               uint32
	ws               *websocket.Conn
	ch               chan *[]byte
	doneCh           chan bool
	monitor          *monitor.Monitor
	eventDispatcher  *events.EventDispatcher
	userNameRegistry *utils.UserNamesRegistry
}

// NewClient initializes a new Client struct with given websocket and Server.
func NewClient(ws *websocket.Conn, monitor *monitor.Monitor, eventsDispatcher *events.EventDispatcher,
	userNameRegistry *utils.UserNamesRegistry, clientID uint32) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	ch := make(chan *[]byte, channelBufSize)
	doneCh := make(chan bool)

	return &Client{clientID, ws, ch, doneCh, monitor, eventsDispatcher, userNameRegistry}
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
		c.monitor.AddDroppedMessage()
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
	log.Println("Listening write to client")
	for {
		select {

		case bytes := <-c.ch:
			before := time.Now()
			err := websocket.Message.Send(c.ws, *bytes)
			after := time.Now()

			if err != nil {
				log.Println(err)
			} else {
				elapsed := after.Sub(before)
				c.monitor.AddSendTime(elapsed)
			}

		case <-c.doneCh:
			c.doneCh <- true
			return
		}
	}
}

func (c *Client) listenRead() {
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
	var data []byte
	err := websocket.Message.Receive(c.ws, &data)
	if err != nil {
		log.Println(err)

		c.doneCh <- true
		c.eventDispatcher.FireUserLeft(&events.UserLeft{ClientID: c.id})
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
		c.eventDispatcher.FireUserInput(userInputEvent)
	case *pb.UserMessage_TargetAngle:
		targetAngleEvent := events.TargetAngleFromProto(protoUserMessage.GetTargetAngle(), c.id)
		c.eventDispatcher.FireTargetAngle(targetAngleEvent)
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
	ok, err := validateUsername(username)

	if ok {
		c.userNameRegistry.AddUserName(c.id, username)
		c.sendJoinGameAckMessage(&pb.JoinGameAck{Success: true})
		c.eventDispatcher.FireUserJoined(&events.UserJoined{ClientID: c.id, UserName: username})
	} else {
		c.sendJoinGameAckMessage(
			&pb.JoinGameAck{Success: false, Error: err.Error()},
		)
	}
}

func validateUsername(username string) (bool, error) {
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
