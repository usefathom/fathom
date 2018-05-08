package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRequestParams(t *testing.T) {
	startDate := time.Now().AddDate(0, 0, -12)
	endDate := time.Now().AddDate(0, 0, -5)

	r, _ := http.NewRequest("GET", "/", nil)
	r.URL.Query().Add("before", string(endDate.Unix()))
	r.URL.Query().Add("after", string(startDate.Unix()))
	r.URL.Query().Add("limit", string(50))
	params := GetRequestParams(r)

	if params.Limit != 50 {
		t.Errorf("Expected %#v, got %#v", 50, params.Limit)
	}

	if startDate != params.StartDate {
		t.Errorf("Expected %#v, got %#v", startDate.Format("2006-01-02 15:04"), params.StartDate.Format("2006-01-02 15:04"))
	}

	if params.EndDate != endDate {
		t.Errorf("Expected %#v, got %#v", endDate.Format("2006-01-02 15:04"), params.EndDate.Format("2006-01-02 15:04"))
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
