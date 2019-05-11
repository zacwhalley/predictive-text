package markov

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"

	"github.com/zacwhalley/predictivetext/util"
)

// Chain contains a map ("chain") of Prefixes to a list of suffixes
// A Prefix is a string of PrefixLen words joined with spaces
// A suffix is a single word. A Prefix can have multiple suffixes
type Chain struct {
	Chain     map[string][]string
	PrefixLen int
}

// NewChain returns a string with Prefixes of length PrefixLen
func NewChain(PrefixLen int) *Chain {
	return &Chain{make(map[string][]string), PrefixLen}
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
			c.Chain[key] = append(c.Chain[key], s)
			p.Shift(s)
		}
	}
}

func (c *Chain) getWord(key string) string {
	key = util.Clean(key)
	choices := c.Chain[key]
	if len(choices) == 0 {
		return ""
	}

	return choices[rand.Intn(len(choices))]
}
