package data

import (
	"github.com/zacwhalley/reddit-simulator/markov"
)

// DBClient is an interface for database access
type DBClient interface {
	GetChain(userName string) (*UserChainDao, error)
	UpsertChain(userName string, chain *markov.Chain) error
}
