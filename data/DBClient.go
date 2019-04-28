package data

import (
	"github.com/zacwhalley/predictive-text/markov"
)

// DBClient is an interface for database access
type DBClient interface {
	GetChain(users []string) (*UserChainDao, error)
	GetPrediction(input string, n int) ([]string, error)
	UpsertChain(users []string, chain *markov.Chain) error
}
