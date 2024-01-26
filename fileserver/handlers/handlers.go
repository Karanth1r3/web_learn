package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := fmt.Sprintf("Visiting default handler'n")
	fmt.Fprintf(w, "%s", Body)
}
