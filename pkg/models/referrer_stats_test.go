package models

import "testing"

func TestReferrerStatsHandlePageview(t *testing.T) {
	s := ReferrerStats{}
	p1 := &Pageview{
		Duration:     100,
		IsBounce:     false,
		IsNewVisitor: true,
	}
	p2 := &Pageview{
		Duration:     60,
		IsNewVisitor: false,
		IsBounce:     true,
	}
	p3 := &Pageview{
		IsNewSession: true,
		IsBounce:     true,
	}

	// add first pageview & test
	s.HandlePageview(p1)
	if s.Pageviews != 1 {
		t.Errorf("Pageviews: expected %d, got %d", 1, s.Pageviews)
	}
	if s.Visitors != 1 {
		t.Errorf("Visitors: expected %d, got %d", 1, s.Visitors)
	}
	if s.AvgDuration != 100 {
		t.Errorf("AvgDuration: expected %.2f, got %.2f", 100.00, s.AvgDuration)
	}
	if s.BounceRate != 0.00 {
		t.Errorf("BounceRate: expected %.2f, got %.2f", 0.00, s.BounceRate)
	}

	// add second pageview
	s.HandlePageview(p2)
	if s.Pageviews != 2 {
		t.Errorf("Pageviews: expected %d, got %d", 2, s.Pageviews)
	}
	if s.Visitors != 1 {
		t.Errorf("Visitors: expected %d, got %d", 1, s.Visitors)
	}
	if s.AvgDuration != 80 {
		t.Errorf("AvgDuration: expected %.2f, got %.2f", 80.00, s.AvgDuration)
	}
	if s.BounceRate != 0.50 {
		t.Errorf("BounceRate: expected %.2f, got %.2f", 0.50, s.BounceRate)
	}

	// add third pageview
	s.HandlePageview(p3)
	if int64(100.00*s.BounceRate) != 66 {
		t.Errorf("BounceRate: expected %.2f, got %.2f", 0.67, s.BounceRate)
	}

}
