package data

import (
	"github.com/zacwhalley/predictivetext/common"
	"github.com/zacwhalley/predictivetext/markov"
)

// DBClient is an interface for database access
type DBClient interface {
	GetChain(users []string) (*UserChainDao, error)
	UpsertChain(users []string, chain *markov.Chain) error
	GetPrediction(prefix, source string) ([]string, error)
	UpsertPrediction(prefix string, suffixes []common.Pair) error
	GetPredictionMap(p markov.Prefix, depth int) (common.SetMap, error)
}
