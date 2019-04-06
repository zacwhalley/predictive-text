package markov

// modified from https://golang.org/doc/codewalk/markov/

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"regexp"
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

// Remove the last non-empty word from the prefix with ""
func (p prefix) reduce() {
	for i := 0; i < len(p); i++ {
		if p[i] != "" {
			p[i] = ""
			break
		}
	}
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
		s = filter(s)
		key := p.toString()
		c.Chain[key] = append(c.Chain[key], s)
		p.shift(s)
	}
}

// Generate returns a string of n words generated from the chain
func (c *Chain) Generate(n int) string {
	p := make(prefix, c.PrefixLen)
	var words []string
	var next string
	for i := 0; i < n; i++ {
		choices := c.Chain[p.toString()]
		for len(choices) == 0 {
			// No more options. Shorten prefix
			p.reduce()
			choices = c.Chain[p.toString()]
		}
		next = choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.shift(next)
	}

	// Capitalize first word
	if len(words) != 0 {
		words[0] = strings.Title(words[0])
	}

	return strings.Join(words, " ")
}

// filter removes non-word characters, and converts to lowercase
func filter(s string) string {
	// Remove all characters that are not part of words
	symPattern := `[^a-zA-Z0-9 ']`
	symRegex, err := regexp.Compile(symPattern)
	if err != nil {
		panic(err)
	}
	s = symRegex.ReplaceAllString(s, "")

	// Convert to lowercase
	return strings.ToLower(s)
}
