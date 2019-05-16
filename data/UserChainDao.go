package data

import (
	"time"

	"github.com/zacwhalley/predictivetext/markov"
)

// UserChainDao is the data access object / schema for user chain objects
type UserChainDao struct {
	Users        []string
	Chain        *markov.Chain
	LastModified time.Time
}
