package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var (
	SERVER       = ""
	PATH         = ""
	TIMESWAIT    = 0
	TIMESWAITMAX = 5
	in           = bufio.NewReader(os.Stdin)
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Provide (Server) (Path)")
		return
	}

	SERVER = args[1]
	PATH = args[2]
	fmt.Println("Connecting to:", SERVER, "at", PATH)

	// Making signal channel (Sigint?)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	input := make(chan string, 1)
	go getInput(input)
	// Forming url, setting up dial (connecting to server)
	URL := url.URL{Scheme: "ws", Host: SERVER, Path: PATH}
	c, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		log.Println("Error", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})
	//Reading data from ws connection through ReadMessage()
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("ReadMessage() error:", err)
				return
			}
			log.Printf("Received: %s", message)
		}
	}()

	for {
		select {
		// If timeout
		case <-time.After(4 * time.Second):
			log.Println("Provide input", TIMESWAIT)
			TIMESWAIT++
			if TIMESWAIT > TIMESWAITMAX {
				syscall.Kill(syscall.Getpid(), syscall.SIGINT) // Sending interrupt signal through go-code
			}
		case <-done:
			return
		case t := <-input:
			err := c.WriteMessage(websocket.TextMessage, []byte(t))
			if err != nil {
				log.Println("Write error: ", err)
				return
			}
			// Reset wait count if input was received
			TIMESWAIT = 0
			go getInput(input)
		case <-interrupt:
			log.Println("Caught interrupt signal - quitting")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, " "))
			if err != nil {
				log.Println("Write close error:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(2 * time.Second):
			}
			return
		}
	}
}

// Reading from stdin reader, sending data to channel
func getInput(input chan string) {
	result, err := in.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	input <- result
}
