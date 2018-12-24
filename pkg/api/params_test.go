package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGetRequestParams(t *testing.T) {
	startDate := time.Now().AddDate(0, 0, -12)
	endDate := time.Now().AddDate(0, 0, -5)
	limit := 50

	url := fmt.Sprintf("/?after=%d&before=%d&limit=%d", startDate.Unix(), endDate.Unix(), limit)
	r, _ := http.NewRequest("GET", url, nil)
	params := GetRequestParams(r)

	if params.Limit != 50 {
		t.Errorf("Expected %#v, got %#v", 50, params.Limit)
	}

	if startDate.Unix() != params.StartDate.Unix() {
		t.Errorf("Expected %#v, got %#v", startDate.Format("2006-01-02 15:04"), params.StartDate.Format("2006-01-02 15:04"))
	}

	if params.EndDate.Unix() != endDate.Unix() {
		t.Errorf("Expected %#v, got %#v", endDate.Format("2006-01-02 15:04"), params.EndDate.Format("2006-01-02 15:04"))
	}

}
