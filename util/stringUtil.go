package util

import (
	"regexp"
	"strings"
)

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
	specCharPattern := `[^a-zA-Z0-9 '!?\.,]`

	s = RemoveMatch(s, linkPattern)
	s = RemoveMatch(s, specCharPattern)

	return s
}

// Clean removes punctuation from a string for use as a key
func Clean(s string) string {
	specCharPattern := `[^a-zA-Z0-9]`
	s = RemoveMatch(s, specCharPattern)
	s = strings.Trim(s, " ")

	if s == "" {
		s = " "
	}

	return s
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
