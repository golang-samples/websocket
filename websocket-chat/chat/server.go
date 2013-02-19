package chat

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"log"
)

// Chat server.
type Server struct {
	path string
	clients []*Client
	addClient chan *Client
	removeClient chan *Client
	sendAll chan *Message
	messages []*Message
}

// Create new chat server.
func NewServer(path string) *Server {
	clients := make([]*Client, 0)
	addClient := make(chan *Client)
	removeClient := make(chan *Client)
	sendAll := make(chan *Message)
	messages := make([]*Message, 0)
	return &Server{path, clients, addClient, removeClient, sendAll, messages}
}

func (self *Server) AddClient() chan<- *Client {
	return (chan<- *Client)(self.addClient)
}

func (self *Server) RemoveClient() chan<- *Client {
	return (chan<- *Client)(self.removeClient)
}

func (self *Server) SendAll() chan<-*Message {
	return (chan<-*Message)(self.sendAll)
}

func (self *Server) Messages() []*Message {
	msgs := make([]*Message, len(self.messages))
	copy(msgs, self.messages)
	return msgs
}

// Listen and serve.
// It serves client connection and broadcast request.
func (self *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		client := NewClient(ws, self)
		self.addClient <- client
		client.Listen()
		defer ws.Close()
	}
	http.Handle(self.path, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-self.addClient:
			log.Println("Added new client")
			self.clients = append(self.clients, c)
			for _, msg := range self.messages {
				c.Write() <- msg
			}
			log.Println("Now", len(self.clients), "clients connected.")

		// remove a client
		case c := <-self.removeClient:
			log.Println("Remove client")
			for i := range self.clients {
				if self.clients[i] == c {
					self.clients = append(self.clients[:i], self.clients[i+1:]...)
					break
				}
			}

		// broadcast message for all clients
		case msg := <-self.sendAll:
			log.Println("Send all:", msg)
			self.messages = append(self.messages, msg)
			for _, c := range self.clients {
				c.Write() <- msg
			}
		}
	}
}
