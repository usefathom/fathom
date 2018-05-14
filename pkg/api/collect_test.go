package api

import (
	"net/http"
	"testing"
)

func TestShouldCollect(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("User-Agent", "Mozilla/1.0")
	r.Header.Add("Referer", "http://usefathom.com/")
	if v := ShouldCollect(r); v != true {
		t.Errorf("Expected %#v, got %#v", true, false)
	}
}
