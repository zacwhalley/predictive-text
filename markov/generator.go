package markov

import (
	"strings"

	"github.com/zacwhalley/predictive-text/util"
)

// Generator is an interface for generating strings with a chain
type Generator interface {
	Generate(c *Chain, n int) string
}

// WordGenerator generates a specified number of words from a markov chain
type WordGenerator struct {
	Beginning string
}

// SentenceGenerator generates a specified number of sentences from a markov chain
type SentenceGenerator struct {
	Beginning string
}

// Generate returns a randomly generated length n sequence of words
func (g WordGenerator) Generate(c *Chain, n int) string {
	var sentences []string
	for n > 0 {
		sentence, numWritten := makeSentence(c, n, g.Beginning)
		n -= numWritten
		sentences = append(sentences, sentence)
		g.Beginning = ""
	}

	return strings.Join(sentences, " ")
}

// Generate returns a randomly generated length n sequence of sentences
func (g SentenceGenerator) Generate(c *Chain, n int) string {
	var sentences []string

	// Generate words until a sentence is produced
	for i := 0; i < n; i++ {
		sentence, _ := makeSentence(c, 0, g.Beginning)
		sentences = append(sentences, sentence)
		g.Beginning = ""
	}

	return strings.Join(sentences, " ")
}

func makeSentence(c *Chain, wordLimit int, beginning string) (string, int) {
	limit := wordLimit != 0
	var words []string

	p := make(Prefix, c.PrefixLen)

	if len(beginning) > 0 {
		initSentence(&words, &p, beginning, c)
	}

	for i := 0; i < wordLimit || !limit; i++ {
		next := c.getWord(p.ToString())
		for next == "" {

			// No more options. Shorten prefix
			p.Reduce()
			if !p.IsEmpty() {
				next = c.getWord(p.ToString())
			} else {
				// no words found be shortening - end sentence
				if len(words) > 0 {
					words[len(words)-1] += "."
				}

				return strings.Join(words, " "), len(words)
			}
		}
		words = append(words, next)
		if util.EndsSentence(next) {
			// End of a sentence has been added
			break
		}
		p.Shift(next)
	}

	return strings.Join(words, " "), len(words)
}

func initSentence(words *[]string, p *Prefix, beginning string, c *Chain) {
	// Find the last c.PrefixLen words in beginning
	beginWords := strings.Split(beginning, " ")
	beginSuffixPos := util.MaxInt(len(beginWords)-c.PrefixLen, 0)
	prefixStart := beginWords[beginSuffixPos:]

	// shift prefixStart words onto string
	for _, word := range prefixStart {
		p.Shift(word)
	}

	// reduce until a word is found
	next := c.getWord(p.ToString())
	for next == "" {
		p.Reduce()
		next = c.getWord(p.ToString())
	}

	if p.ToString() == " " {
		// beginning matches with nothing - don't print beginning
	} else {
		// Add the whole beginning to words if a single result was found
		copy(*words, beginWords)
	}
}
