package common

import (
	"bufio"
	"fmt"
	"io"

	"github.com/zacwhalley/predictivetext/domain"
	"github.com/zacwhalley/predictivetext/util"
)

// Chain contains a map ("chain") of Prefixes to a list of suffixes
// A Prefix is a string of PrefixLen words joined with spaces
// A suffix is a single word. A Prefix can have multiple suffixes
type Chain struct {
	data      domain.SetMap
	prefixLen int
}

// NewChain returns a string with Prefixes of length PrefixLen
func NewChain(prefixLen int) Chain {
	return Chain{make(SetMap), prefixLen}
}

// GetData returns the chain's chain Data value
func (c Chain) GetData() domain.SetMap {
	return c.data
}

// GetPrefixLen returns the chain's prefixLen value
func (c Chain) GetPrefixLen() int {
	return c.prefixLen
}

// Get returns the value in the chain indexed by key
func (c Chain) Get(key string) (domain.Set, bool) {
	set, ok := c.data.Get(key)
	return set, ok
}

// Build reads text from the provided Reader and parses it into Prefixes
// and suffixes stored in the chain
func (c Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		s = util.Filter(s)
		if s != "" {
			// If s was filtered out
			key := p.ToString()
			c.data.Add(key, s)
			p.Shift(s)
		}
	}
}
