package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespond(t *testing.T) {
	w := httptest.NewRecorder()
	respond(w, http.StatusOK, 15)

	if w.Code != 200 {
		t.Errorf("Invalid response code")
	}

	// assert json header
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Invalid response header for Content-Type")
	}

	// assert json response
	var d int
	err := json.NewDecoder(w.Body).Decode(&d)
	if err != nil {
		t.Errorf("Invalid response body: %s", err)
	}

}
