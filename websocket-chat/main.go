package main

import (
	"net/http"
	"log"
	"./chat"
)

func main() {
	server := chat.NewServer("/entry")
	go server.Listen()
	http.Handle("/", http.FileServer(http.Dir(".")))
	for _, path := range []string{"lib", "js"} {
		http.Handle("/"+path+"/", http.StripPrefix("/"+path+"/", http.FileServer(http.Dir(path))))
	}
	log.Fatal(http.ListenAndServe(":8080", nil))
}
