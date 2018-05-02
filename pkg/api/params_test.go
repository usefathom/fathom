package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequestedLimit(t *testing.T) {
	r, _ := http.NewRequest("GET", "", nil)
	limit := getRequestedLimit(r)

	if limit != defaultLimit {
		t.Errorf("Expected limit of %d does not match %d", defaultLimit, limit)
	}

	r, _ = http.NewRequest("GET", "?limit=50", nil)
	limit = getRequestedLimit(r)
	if limit != 50 {
		t.Errorf("Expected limit of %d does not match %d", defaultLimit, limit)
	}
}

func TestGetRequestedPeriods(t *testing.T) {
	r, _ := http.NewRequest("GET", "?before=500&after=100", nil)
	before, after := getRequestedPeriods(r)

	if before != 500 || after != 100 {
		t.Error("Expected URl argument for `before` or `after` does not match")
	}
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
