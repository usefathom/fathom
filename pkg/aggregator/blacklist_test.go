package aggregator

import (
	"testing"
)

func TestBlacklistHas(t *testing.T) {
	b, err := newBlacklist()
	if err != nil {
		t.Error(err)
	}

	table := map[string]bool{
		"03e.info":      true,
		"zvetki.ru":     true,
		"usefathom.com": false,
		"foo.03e.info":  true, // sub-string match
	}

	for r, e := range table {
		if v := b.Has(r); v != e {
			t.Errorf("Expected %v, got %v", e, v)
		}
	}
}
