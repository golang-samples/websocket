package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
)

type Client struct {
	ws *websocket.Conn
}

func (c *Client) Read() string {
	data := make([]byte, 512)
	n, err := c.ws.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	return string(data[:n])
}

func (c *Client) Write(data string) int {
	n, err := c.ws.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func NewClient(ws *websocket.Conn) *Client {
	client := &Client{
		ws: ws,
	}
	return client
}

type Clients struct {
	clients []*Client
}

func (c *Clients) Add(client *Client) {
	c.clients = append(c.clients, client)
}

func (c *Clients) Broadcast(data string) {
	for _, c := range c.clients {
		log.Println(c)
		c.Write(data)
	}
}

func NewClients() *Clients {
	return &Clients{
		clients: []*Client{},
	}
}

var clients = NewClients()

func echoHandler(ws *websocket.Conn) {
	client := NewClient(ws)
	clients.Add(client)
	log.Println(len(clients.clients))

	for {
		message := client.Read()
		fmt.Printf("Receive: %s\n", message)

		clients.Broadcast(message)
		fmt.Printf("Send: %s\n", message)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
