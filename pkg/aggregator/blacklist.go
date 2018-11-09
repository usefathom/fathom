package aggregator

import (
	"bufio"
	"bytes"
	"strings"
)

type blacklist struct {
	data []byte
}

func newBlacklist() (*blacklist, error) {
	var err error
	b := &blacklist{}
	b.data, err = Asset("blacklist.txt")
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Has returns true if the given domain appears on the blacklist
// Uses sub-string matching, so if usesfathom.com is blacklisted then this function will also return true for danny.usesfathom.com
func (b *blacklist) Has(r string) bool {
	if r == "" {
		return false
	}

	scanner := bufio.NewScanner(bytes.NewReader(b.data))
	domain := ""

	for scanner.Scan() {
		domain = scanner.Text()
		if strings.HasSuffix(r, domain) {
			return true
		}
	}

	return false
}
