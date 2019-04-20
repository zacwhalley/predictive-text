package markov

import (
	"strings"

	"github.com/zacwhalley/reddit-simulator/util"
)

// Generator is an interface for generating strings with a chain
type Generator interface {
	Generate(c *Chain, n int) string
}

// WordGenerator generates a specified number of words from a markov chain
type WordGenerator struct{}

// SentenceGenerator generates a specified number of sentences from a markov chain
type SentenceGenerator struct{}

// Generate returns a randomly generated length n sequence of words
func (g WordGenerator) Generate(c *Chain, n int) string {
	var sentences []string
	for n > 0 {
		sentence, numWritten := makeSentence(c, n)
		n -= numWritten
		sentences = append(sentences, sentence)
	}

	return strings.Join(sentences, " ")
}

// Generate returns a randomly generated length n sequence of sentences
func (g SentenceGenerator) Generate(c *Chain, n int) string {
	var sentences []string

	// Generate words until a sentence is produced
	for i := 0; i < n; i++ {
		sentence, _ := makeSentence(c, 0)
		sentences = append(sentences, sentence)
	}

	return strings.Join(sentences, " ")
}

func makeSentence(c *Chain, wordLimit int) (string, int) {
	limit := wordLimit != 0
	p := make(prefix, c.PrefixLen)
	var words []string
	for i := 0; i < wordLimit || !limit; i++ {
		next := c.getWord(p.toString())
		for next == "" {
			// No more options. Shorten prefix
			p.reduce()
			next = c.getWord(p.toString())
		}
		words = append(words, next)
		if util.EndsSentence(next) {
			// End of a sentence has been added
			break
		}
		p.shift(next)
	}

	// Capitalize first word
	if len(words) != 0 {
		words[0] = strings.Title(words[0])
	}

	return strings.Join(words, " "), len(words)
}
