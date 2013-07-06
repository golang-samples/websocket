package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
)

func echoHandler(ws *websocket.Conn) {
	msg := make([]byte, 512)
	_, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg)

	_, err = ws.Write(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", msg)
}

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
