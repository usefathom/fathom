package api

import (
	"net/http"
	"testing"
)

func TestShouldCollect(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("User-Agent", "Mozilla/1.0")
	r.Header.Add("Referer", "http://usefathom.com/")
	if v := shouldCollect(r); v != true {
		t.Errorf("Expected %#v, got %#v", true, false)
	}
}

func TestParsePathname(t *testing.T) {
	if v := parsePathname("/"); v != "/" {
		t.Errorf("error parsing pathname. expected %#v, got %#v", "/", v)
	}

	if v := parsePathname("about"); v != "/about" {
		t.Errorf("error parsing pathname. expected %#v, got %#v", "/about", v)
	}
}

func TestParseReferrer(t *testing.T) {
	e := "https://usefathom.com"

	// normal
	if v := parseReferrer("https://usefathom.com"); v != e {
		t.Errorf("error parsing referrer. expected %#v, got %#v", e, v)
	}

	// amp in query string
	if v := parseReferrer("https://usefathom.com?amp=1&utm_source=foo"); v != e {
		t.Errorf("error parsing referrer. expected %#v, got %#v", e, v)
	}

	// amp in pathname
	if v := parseReferrer("https://usefathom.com/amp/"); v != e {
		t.Errorf("error parsing referrer. expected %#v, got %#v", e, v)
	}

	e = "https://usefathom.com/about?page_id=500"
	if v := parseReferrer("https://usefathom.com/about/amp/?amp=1&page_id=500&utm_campaign=foo"); v != e {
		t.Errorf("error parsing referrer. expected %#v, got %#v", e, v)
	}

}

func TestParseHostname(t *testing.T) {
	e := "https://usefathom.com"
	if v := parseHostname("https://usefathom.com"); v != e {
		t.Errorf("error parsing hostname. expected %#v, got %#v", e, v)
	}

	e = "http://usefathom.com"
	if v := parseHostname("http://usefathom.com"); v != e {
		t.Errorf("error parsing hostname. expected %#v, got %#v", e, v)
	}
}
