package data

import (
	"github.com/zacwhalley/reddit-simulator/markov"
)

// DBClient is an interface for database access
type DBClient interface {
	GetChain(users []string) (*UserChainDao, error)
	UpsertChain(users []string, chain *markov.Chain) error
}
