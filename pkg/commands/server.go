package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/usefathom/fathom/pkg/api"
)

// Server starts the HTTP server, listening on the given port
func Server(port int) {
	r := api.Routes()
	log.Printf("Now serving on port %d/\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.LoggingHandler(os.Stdout, r))
	if err != nil {
		log.Println(err)
	}
}
