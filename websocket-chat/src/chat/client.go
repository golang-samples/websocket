package chat

import (
	"io"
	"log"

	"code.google.com/p/go.net/websocket"
)

const channelBufSize = 100

// Chat client.
type Client struct {
	ws     *websocket.Conn
	server *Server
	ch     chan *Message
	doneCh chan bool
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	} else if server == nil {
		panic("server cannot be nil")
	}

	ch := make(chan *Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{ws, server, ch, doneCh}
}

// Get websocket connection.
func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

// Get Write channel
func (c *Client) Write() chan<- *Message {
	return c.ch
}

// Get done channel.
func (c *Client) DoneCh() chan<- bool {
	return c.doneCh
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

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.ws, msg)

		// receive done request
		case <-c.doneCh:
			c.server.DelCh() <- c
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
			c.server.DelCh() <- c
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.ErrCh() <- err
			} else {
				c.server.SendAllCh() <- &msg
			}
		}
	}
}
