package common

import (
	"strings"

	"github.com/zacwhalley/predictivetext/domain"
	"github.com/zacwhalley/predictivetext/util"
)

// Prefix is a markov chain prefix of one of more words
type Prefix []string

// MakePrefix creates a prefix of size prefixLen from input
func MakePrefix(input string, prefixLen int) Prefix {
	words := strings.Split(strings.TrimSpace(input), " ")
	limit := util.MaxInt(len(words)-prefixLen, 0)

	return words[limit:]
}

// ToString returns the Prefix as a string (for use as a map key)
func (p Prefix) ToString() string {
	s := strings.Join(p, " ")
	return strings.ToLower(util.Clean(s))
}

// IsEmpty returns true if ToString returns an empty value
func (p Prefix) IsEmpty() bool {
	return p.ToString() == " "
}

// Copy returns a value copy of a prefix
func (p Prefix) Copy() domain.Prefix {
	var newPrefix Prefix
	copy(newPrefix, p)
	return newPrefix
}

// Clear removes all words from the Prefix
func (p Prefix) Clear() {
	p = make([]string, len(p))
}

// Shift removes the first word from the Prefix and appends the given word
func (p Prefix) Shift(word string) {
	if util.EndsSentence(p.Last()) {
		// word ends with one of ?.! -> end of sentence
		p.Clear()
	} else {
		copy(p, p[1:])
	}
	p[len(p)-1] = util.Clean(word)
}

// Reduce removes the last non-empty word from the Prefix with ""
func (p Prefix) Reduce() {
	for i := 0; i < len(p); i++ {
		if p[i] != "" {
			p[i] = ""
			break
		}
	}
}

// Last returns the last word in the prefix
func (p Prefix) Last() string {
	return p[len(p)-1]
}
