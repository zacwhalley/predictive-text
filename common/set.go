package common

import (
	"strings"

	"github.com/zacwhalley/predictivetext/domain"
)

// Set represents a map[string]int used for counting the occurance of strings
type Set map[string]int

// Add adds a string to the set. Incrementing by 1 if it already exists
func (s Set) Add(key string) {
	// awful workaround for bad db design - remove prefix $
	if strings.HasPrefix(key, "$") {
		key = strings.Replace(key, "$", "", 1)
	}
	if strings.TrimSpace(key) == "" {
		return
	}
	if _, ok := s[key]; ok {
		s[key]++
	} else {
		s[key] = 1
	}
}

// IsEmpty returne true if the set is empty, false otherwise
func (s Set) IsEmpty() bool {
	return len(s) <= 0
}

// Union adds all values from one set into another, merging duplicates
func (s Set) Union(set domain.Set) {
	if set.IsEmpty() {
		return
	}

	rangeSet, _ := set.(Set) // type check - only works when implemented as rangeable
	for key, count := range rangeSet {
		if strings.TrimSpace(key) == "" {
			continue
		}
		if _, ok := s.Get(key); ok {
			s[key] += count
		} else {
			s[key] = count
		}
	}
}

// Get returns the value mapped to key
func (s Set) Get(key string) (int, bool) {
	value, ok := s[key]
	return value, ok
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

// MakePrefixSet creates a set from s where all keys are prepended by text
func (s Set) MakePrefixSet(text string, weight int) domain.Set {
	newSet := make(Set)
	for key, value := range s {
		newSet[strings.Join([]string{text, key}, " ")] = value + weight
	}

	return newSet
}

// ToPairs converts the key-value sets to a list of pairs
func (s Set) ToPairs() []domain.Pair {
	list := make([]domain.Pair, len(s))
	for key, value := range s {
		list = append(list, domain.Pair{Key: key, Value: value})
	}

	return list
}

// SetMap is a map from a key to a set
type SetMap map[string]Set

// MakeSetMap converts a primitive map into a set map
func MakeSetMap(data map[string]map[string]int) SetMap {
	newSetMap := make(SetMap)
	for key, value := range data {
		newSetMap[key] = value
	}

	return newSetMap
}

// Add adds value to the set associated with key
func (sm SetMap) Add(key, value string) {
	if _, ok := sm[key]; !ok {
		sm[key] = make(Set)
	}
	sm[key].Add(value)
}

// ToPrimitive returns the SetMap as a primitive map type w/ no methods
func (sm SetMap) ToPrimitive() map[string](map[string]int) {
	primitiveMap := make(map[string](map[string]int))
	for key, value := range sm {
		primitiveMap[key] = value
	}

	return primitiveMap
}

// Get returns the set for a given key
func (sm SetMap) Get(key string) (domain.Set, bool) {
	set, ok := sm[key]
	return set, ok
}

// Union merges two SetMaps, adding weights
func (sm SetMap) Union(other domain.SetMap) {
	rangeSetMap, _ := other.(SetMap) // typecheck - only works if range works
	for key, value := range rangeSetMap {
		if _, ok := sm[key]; ok {
			sm[key].Union(value)
		} else {
			sm[key] = value
		}
	}
}
