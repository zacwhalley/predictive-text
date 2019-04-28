package util

import (
	"regexp"
	"strings"
)

// EndsSentence returns true if s ends with a ./!/? and is not
// a common word like Mr. or Dr.
func EndsSentence(s string) bool {
	// check for matches against common abbreviations
	abbrevs := []string{"mr.", "mrs.", "ms.", "etc.", "jr.", "sr.", "dr."}
	for _, val := range abbrevs {
		if val == strings.ToLower(s) {
			return false
		}
	}

	// check for matches against cases that would end a sentence
	return DoesEndWith(s, []string{".", "!", "?"})
}

// DoesEndWith returns true if s has any string from match as a suffix
func DoesEndWith(s string, match []string) bool {
	for _, pattern := range match {
		if strings.HasSuffix(s, pattern) {
			return true
		}
	}

	return false
}

// Filter removes links and unwanted punctuation
func Filter(s string) string {
	linkPattern := `[-a-zA-Z0-9@:%_\+.~#?&//=]{2,256}\.[a-z]{2,4}\b(\/[-a-zA-Z0-9@:%_\+.~#?&//=]*)?`
	miscMarkdownPattern := `[&[a-zA-Z]+;]`

	s = RemoveMatch(s, linkPattern)
	s = RemoveMatch(s, miscMarkdownPattern)

	return s
}

// Clean removes punctuation from a string for use as a key
func Clean(s string) string {
	specCharPattern := `[^a-zA-Z0-9 ]`
	s = RemoveMatch(s, specCharPattern)
	s = strings.Trim(s, " ")

	if s == "" {
		s = " "
	}

	return strings.ToLower(s)
}

// RemoveMatch removes all substrings in s that match pattern
func RemoveMatch(s, pattern string) string {
	// Remove all characters that are not part of words
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	s = regex.ReplaceAllString(s, "")
	return s
}
