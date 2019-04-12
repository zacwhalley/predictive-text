package markov

// modified from https://golang.org/doc/codewalk/markov/

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"github.com/zacwhalley/reddit-simulator/util"
)

// Prefix is a markov chain prefix of one of more words
type prefix []string

// toString returns the prefix as a string (for use as a map key)
func (p prefix) toString() string {
	s := strings.Join(p, " ")

	return util.Clean(s)
}

func (p prefix) last() string {
	return p[len(p)-1]
}

// Shift removes the first word from the prefix and appends the given word
func (p prefix) shift(word string) {
	endChars := []string{".", "!", "?"}
	if util.DoesEndWith(p.last(), endChars) {
		// word ends with one of ?.! -> end of sentence
		p.clear()
	} else {
		copy(p, p[1:])
	}
	p[len(p)-1] = util.Clean(word)
}

// Remove the last non-empty word from the prefix with ""
func (p prefix) reduce() {
	for i := 0; i < len(p); i++ {
		if p[i] != "" {
			p[i] = ""
			break
		}
	}
}

// clear removes all words from the prefix
func (p prefix) clear() {
	p = make([]string, len(p))
}

// Chain contains a map ("chain") of prefixes to a list of suffixes
// A prefix is a string of prefixLen words joined with spaces
// A suffix is a single word. A prefix can have multiple suffixes
type Chain struct {
	Chain     map[string][]string
	PrefixLen int
}

// NewChain returns a string with prefixes of length prefixLen
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and parses it into prefixes
// and suffixes stored in the chain
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(prefix, c.PrefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		s = util.Filter(s)
		if s != "" {
			// If s was filtered out
			key := p.toString()
			c.Chain[key] = append(c.Chain[key], s)
			p.shift(s)
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

// Generate returns a string of n words generated from the chain
func (c *Chain) Generate(n int) string {
	p := make(prefix, c.PrefixLen)
	var words []string
	for i := 0; i < n; i++ {
		next := c.getWord(p.toString())
		for next == "" {
			// No more options. Shorten prefix
			p.reduce()
			next = c.getWord(p.toString())
		}
		words = append(words, next)
		p.shift(next)
	}

	// Capitalize first word
	if len(words) != 0 {
		words[0] = strings.Title(words[0])
	}

	return strings.Join(words, " ")
}
