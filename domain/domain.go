package domain

import (
	"io"
	"time"
)

// PredictionSvc is a service for generating predictions
type PredictionSvc interface {
	GetPrediction(input string) ([]string, error)
	SavePrediction(Prediction) error
	GeneratePredictionSet(input string) error
}

// Set counts occurrences of strings
type Set interface {
	Add(key string)
	Get(key string) (int, bool)
	IsEmpty() bool
	GetWeight(key string) (float64, bool)
	Union(set Set)
	MakePrefixSet(text string, weight int) Set
	ToPairs() []Pair
}

// SetMap is a map from a key to a set
type SetMap interface {
	Add(key, value string)
	Get(key string) (Set, bool)
	Union(other SetMap)
	ToPrimitive() map[string]map[string]int
}

// Chain is a markov chain mapping prefixes to suffixes
type Chain interface {
	GetData() SetMap
	GetPrefixLen() int
	Get(key string) (Set, bool)
	Build(r io.Reader)
}

// Prefix is a markov chain prefix of one or more words
type Prefix interface {
	ToString() string
	IsEmpty() bool
	Copy() Prefix
	Clear()
	Shift(word string)
	Reduce()
	Last() string
}

// DBClient is an interface for database access
type DBClient interface {
	GetChainByID(id string) (UserChainDao, error)
	UpsertChain(users []string, chain Chain) error
	GetPrediction(prefix, source string) (Prediction, error)
	UpsertPrediction(prediction Prediction) error
}

// PredictionRequest is the dto for requesting a prediction
type PredictionRequest struct {
	Input string `json:"input"`
}

// PredictionResponse is the Dto for returning a prediction
type PredictionResponse struct {
	Input       string   `json:"input"`
	Predictions []string `json:"predictions"`
}

// PredictionDao is the data access object / schema for a prediction
type PredictionDao struct {
	Source   string `bson:"source"`
	Prefix   string `bson:"prefix"`
	Suffixes []Pair `bson:"suffixes"`
}

// Prediction is a struct containing
type Prediction struct {
	Prefix   string
	Suffixes []Pair
}

// UserChainDao is the data access object / schema for user chain objects
type UserChainDao struct {
	Users        []string                  `bson:"users"`
	Data         map[string]map[string]int `bson:"data"`
	PrefixLen    int                       `bson:"prefixlen"`
	LastModified time.Time                 `bson:"lastmodified"`
}

// Pair is a struct containing a string and int
type Pair struct {
	Key   string
	Value int
}
