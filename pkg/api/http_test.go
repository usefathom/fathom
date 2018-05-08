package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestRespond(t *testing.T) {
	w := httptest.NewRecorder()
	respond(w, 15)

	// assert json header
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Invalid Content-Type header")
	}

	// assert json response
	var d int
	err := json.NewDecoder(w.Body).Decode(&d)
	if err != nil {
		t.Errorf("Invalid response body: %s", err)
	}

}
