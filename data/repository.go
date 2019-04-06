package data

import "github.com/zacwhalley/reddit-simulator/markov"

// DBClient is an interface for database access
type DBClient interface {
	Init()
	GetChain(userName string) markov.Chain
	PostChain(userName string, chain markov.Chain)
}

// MongoClient is a client for mongoDB
type MongoClient struct {
}

func (m MongoClient) Init() {

}

// GetChain gets the chain for the user userName
func (m MongoClient) GetChain(userName string) markov.Chain {
	return *markov.NewChain(3)
}

// PostChain saves the chain
func (m MongoClient) PostChain(userName string, chain markov.Chain) {

}
