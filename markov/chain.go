package markov

import (
	"bufio"
	"fmt"
	"io"

	"github.com/zacwhalley/predictivetext/common"

	"github.com/zacwhalley/predictivetext/util"
)

// Chain contains a map ("chain") of Prefixes to a list of suffixes
// A Prefix is a string of PrefixLen words joined with spaces
// A suffix is a single word. A Prefix can have multiple suffixes
type Chain struct {
	Chain     common.SetMap
	PrefixLen int
}

// NewChain returns a string with Prefixes of length PrefixLen
func NewChain(PrefixLen int) *Chain {
	return &Chain{make(common.SetMap), PrefixLen}
}

// Build reads text from the provided Reader and parses it into Prefixes
// and suffixes stored in the chain
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.PrefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		s = util.Filter(s)
		if s != "" {
			// If s was filtered out
			key := p.ToString()
			c.Chain.Add(key, s)
			p.Shift(s)
		}
	}
}
