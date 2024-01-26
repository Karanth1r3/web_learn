package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var (
	PORT     = ":1234"
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!\n")
	fmt.Fprintf(w, "Please use /ws for WebSocket")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection from:", r.Host)
	// Creating ws connection based on http request handler ("Upgrading")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader.Upgrade: ", err)
		return
	}
	// Close when func executes
	defer ws.Close()
	// handling all incoming messages (Reading & writing them back [echo])
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("From", r.Host, "read", err)
			break
		}
		log.Println("Received: ", string(message))
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("WriteMessage:", err)
			break
		}
	}
}

func main() {
	// If port is typed in cmd => use it
	args := os.Args
	if len(args) != 1 {
		PORT = ":" + args[1]
	}
	// Router & server initializing. Server with specified port, mux & timeouts (read, write, idle)
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	// Linking handlers to urls
	mux.Handle("/", http.HandlerFunc(rootHandler))
	mux.Handle("/ws", http.HandlerFunc(wsHandler))

	// Starting sever
	log.Println("Listening to TCP Port", PORT)
	err := s.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}
