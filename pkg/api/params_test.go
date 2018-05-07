package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequestParams(t *testing.T) {
	// TODO: Implement this
}

func TestParseMajorMinor(t *testing.T) {
	actual := parseMajorMinor("50.0.0")
	expected := "50.0"
	if actual != expected {
		t.Errorf("Return value should be %s, is %s instead", expected, actual)
	}

	actual = parseMajorMinor("1.1")
	expected = "1.1"
	if actual != expected {
		t.Errorf("Return value should be %s is %s instead", expected, actual)
	}
}

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
