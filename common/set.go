package common

import "strings"

// Pair is a struct containing a string and int
type Pair struct {
	Key   string
	Value int
}

// Set represents a map[string]int used for counting the occurance of strings
type Set map[string]int

// Add adds a string to the set. Incrementing by 1 if it already exists
func (s Set) Add(key string) {
	// awful workaround for bad db design - remove prefix $
	if strings.HasPrefix(key, "$") {
		key = strings.Replace(key, "$", "", 1)
	}
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

// GetWeight returns the weight of the value divided by the total
// number of items in the set
func (s Set) GetWeight(key string) (float64, bool) {
	if weight, ok := s[key]; ok {
		return float64(weight) / float64(len(s)), true
	}
	// no value exists for key
	return -1.0, false
}

// ToPairs converts the key-value sets to a list of pairs
func (s Set) ToPairs() []Pair {
	list := make([]Pair, len(s))
	for key, value := range s {
		list = append(list, Pair{Key: key, Value: value})
	}

	return list
}

// SetMap is a map from a key to a set
type SetMap map[string]Set

// Add adds value to the set associated with key
func (sm SetMap) Add(key, value string) {
	if _, ok := sm[key]; !ok {
		sm[key] = make(Set)
	}
	sm[key].Add(value)
}

// Union merges two SetMaps, adding weights
func (sm SetMap) Union(other SetMap) {
	for key, value := range other {
		if _, ok := sm[key]; ok {
			sm[key].AddSet(value)
		} else {
			sm[key] = value
		}
	}
}
