package data

import (
	"time"

	"github.com/zacwhalley/reddit-simulator/markov"
)

// UserChainDao is the data access object for user chain objects
type UserChainDao struct {
	User         string
	Chain        *markov.Chain
	LastModified time.Time
}