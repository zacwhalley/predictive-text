package main

// modified from https://golang.org/doc/codewalk/markov/

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

// Prefix is a markov chain prefix of one of more words
type prefix []string

// toString returns the prefix as a string (for use as a map key)
func (p prefix) toString() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the prefix and appends the given word
func (p prefix) shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes
// A prefix is a string of prefixLen words joined with spaces
// A suffix is a single word. A prefix can have multiple suffixes
type chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a string with prefixes of length prefixLen
func newChain(prefixLen int) *chain {
	return &chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and parses it into prefixes
// and suffixes stored in the chain
func (c *chain) build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.toString()
		c.chain[key] = append(c.chain[key], s)
		p.shift(s)
	}
}

func (c *chain) Generate(n int) string {
	p := make(prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.toString()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.shift(next)
	}
	return strings.Join(words, " ")
}
