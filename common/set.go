package common

// Set represents a map[string]int used for counting the occurance of strings
type Set map[string]int

// NewSet initializes a new set
func NewSet() Set {
	return make(map[string]int)
}

// AddString adds a string to the set. Incrementing by 1 if it already exists
func (s Set) AddString(key string) {
	if _, ok := s[key]; ok {
		s[key]++
	} else {
		s[key] = 1
	}
}

// AddSet adds all values from one set into another
func (s Set) AddSet(set Set) {
	for key, count := range set {
		if _, ok := s[key]; ok {
			s[key] += count
		} else {
			s[key] = 1
		}
	}
}
