package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Handler is our custom HTTP handler with error returns
type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		HandleError(w, r, err)
	}
}

// HandlerFunc takes a custom Handler func and converts it to http.HandlerFunc
func HandlerFunc(fn Handler) http.HandlerFunc {
	return http.HandlerFunc(Handler(fn).ServeHTTP)
}

// HandleError handles errors
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	log.WithFields(log.Fields{
		"request": r.Method + " " + r.RequestURI,
		"error":   err,
	}).Error("error handling request")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("false"))
}

type envelope struct {
	Data  interface{} `json:",omitempty"`
	Error interface{} `json:",omitempty"`
}

func respond(w http.ResponseWriter, d interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(d)
	return err
}
