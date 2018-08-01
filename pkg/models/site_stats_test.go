package models

import (
	"testing"
)

func TestSiteStatsFormattedDuration(t *testing.T) {
	s := SiteStats{
		AvgDuration: 100.00,
	}
	e := "1:40"
	if v := s.FormattedDuration(); v != e {
		t.Errorf("FormattedDuration: expected %s, got %s", e, v)
	}

	s.AvgDuration = 1040.22
	e = "17:20"
	if v := s.FormattedDuration(); v != e {
		t.Errorf("FormattedDuration: expected %s, got %s", e, v)
	}
}
