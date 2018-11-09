package aggregator

import (
	"net/url"
	"testing"
)

func TestParseReferrer(t *testing.T) {
	testsValid := map[string]*url.URL{
		"https://www.usefathom.com/?utm_source=github": &url.URL{
			Scheme: "https",
			Host:   "www.usefathom.com",
			Path:   "/",
		},
		"https://www.usefathom.com/privacy/amp/?utm_source=github": &url.URL{
			Scheme: "https",
			Host:   "www.usefathom.com",
			Path:   "/privacy/",
		},
	}
	testsErr := []string{
		"mysite.com",
		"foobar",
		"",
	}

	for r, e := range testsValid {
		v, err := parseReferrer(r)
		if err != nil {
			t.Error(err)
		}

		if v.Host != e.Host {
			t.Errorf("Invalid Host: expected %s, got %s", e.Host, v.Host)
		}

		if v.Scheme != e.Scheme {
			t.Errorf("Invalid Scheme: expected %s, got %s", e.Scheme, v.Scheme)
		}

		if v.Path != e.Path {
			t.Errorf("Invalid Path: expected %s, got %s", e.Path, v.Path)
		}

	}

	for _, r := range testsErr {
		v, err := parseReferrer(r)
		if err == nil {
			t.Errorf("Expected err, got %#v", v)
		}
	}

}
