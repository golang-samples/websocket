package chat

import (
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

// Chat server.
type Server struct {
	pattern   string
	messages  []*Message
	clients   []*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []*Message{}
	clients := []*Client{}
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{pattern, messages, clients, addCh, delCh, sendAllCh, doneCh, errCh}
}

func (s *Server) AddCh() chan<- *Client {
	return s.addCh
}

func (s *Server) DelCh() chan<- *Client {
	return s.delCh
}

func (s *Server) SendAllCh() chan<- *Message {
	return s.sendAllCh
}

func (s *Server) DoneCh() chan<- bool {
	return s.doneCh
}

func (s *Server) ErrCh() chan<- error {
	return s.errCh
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write() <- msg
	}
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		c.Write() <- msg
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer ws.Close()

		client := NewClient(ws, s)
		s.addCh <- client
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients = append(s.clients, c)
			s.sendPastMessages(c)
			log.Println("Now", len(s.clients), "clients connected.")

			// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			for i := range s.clients {
				if s.clients[i] == c {
					s.clients = append(s.clients[:i], s.clients[i+1:]...)
					break
				}
			}

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
