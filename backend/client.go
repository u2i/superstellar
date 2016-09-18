package backend

import (
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

// Client struct holds client-specific variables.
type Client struct {
	id      uint32
	ws      *websocket.Conn
	server  *Server
	ch      chan *[]byte
	doneCh  chan bool
	monitor *Monitor
}

// NewClient initializes a new Client struct with given websocket and Server.
func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	ch := make(chan *[]byte, channelBufSize)
	doneCh := make(chan bool)
	id := server.GenerateID()
	monitor := server.monitor

	return &Client{id, ws, server, ch, doneCh, monitor}
}

// Conn returns client's websocket.Conn struct.
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

// SendSpace sends game state to the client.
func (c *Client) SendSpace(bytes *[]byte) {
	select {
	case c.ch <- bytes:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected", c.id)
		c.server.Err(err)
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
				c.monitor.addSendTime(elapsed)
			}

		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var move Move
			err := websocket.JSON.Receive(c.ws, &move)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.Move(&move)
			}
		}
	}
}
