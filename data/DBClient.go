package data

import "github.com/zacwhalley/predictive-text/markov"

// DBClient is an interface for database access
type DBClient interface {
	GetChain(users []string) (*UserChainDao, error)
	GetPredictionMap(p markov.Prefix, depth int) (map[string][]string, error)
}
