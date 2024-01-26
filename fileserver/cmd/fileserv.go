package main

import (
	"fmt"
	"net/http"
	"web_learn/fileserver/handlers"
)

var PORT = ":8765"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.DefaultHandler)

	fileServer := http.FileServer(http.Dir("/tmp"))
	// Static does not exist in file space => it is necessary to trim it while handling data. Before this all URLS starting with static will be handleded by fileServer
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	fmt.Println("Starting server on", PORT)
	// Launching server
	err := http.ListenAndServe(PORT, mux)
	fmt.Println(err)
}
