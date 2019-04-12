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
	s := strings.Join(p, " ")

	return clean(s)
}

// Shift removes the first word from the prefix and appends the given word
func (p prefix) shift(word string) {
	if false {
		// word ends with one of ?.! -> end of sentence
		p.clear()
	} else {
		copy(p, p[1:])
		p[len(p)-1] = clean(word)
	}
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
		s = filter(s)
		if s != "" {
			// If s was filtered out
			key := p.toString()
			c.Chain[key] = append(c.Chain[key], s)
			p.shift(s)
		}
	}
}

func (c *Chain) getWord(key string) string {
	key = clean(key)
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

// filter removes links and unwanted punctuation
func filter(s string) string {
	linkPattern := `[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?`
	specCharPattern := `[^a-zA-Z0-9 '!?\.,]`

	s = removeMatch(s, linkPattern)
	s = removeMatch(s, specCharPattern)

	return s
}

// clean removes punctuation from a string for use as a key
func clean(s string) string {
	specCharPattern := `[^a-zA-Z0-9]`
	s = removeMatch(s, specCharPattern)
	s = strings.ToLower(s)
	s = strings.Trim(s, " ")

	if s == "" {
		s = " "
	}

	return s
}

// removeMatch removes all substrings in s that match pattern
func removeMatch(s, pattern string) string {
	// Remove all characters that are not part of words
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	s = regex.ReplaceAllString(s, "")
	return s
}
