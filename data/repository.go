package data

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/zacwhalley/reddit-simulator/markov"
	"go.mongodb.org/mongo-driver/mongo"
)

// DBClient is an interface for database access
type DBClient interface {
	GetChain(userName string) markov.Chain
	PostChain(userName string, chain markov.Chain)
}

// MongoClient is a client for mongoDB
type MongoClient struct {
	client *mongo.Client
}

// NewMongoClient creates a new client and establishes a connection to a mongo database
func NewMongoClient() *MongoClient {
	newClient := &MongoClient{}

	options := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	newClient.client, err = mongo.Connect(context.TODO(), options)
	if err != nil {
		log.Fatal(err)
	}

	return newClient
}

// GetChain gets the chain for the user userName
func (m MongoClient) GetChain(userName string) markov.Chain {
	if m.client == nil {
		log.Fatal("No connection to MongoDB")
	}

	return *markov.NewChain(3)
}

// PostChain saves the chain
func (m MongoClient) PostChain(userName string, chain markov.Chain) {
	if m.client == nil {
		log.Fatal("No connection to MongoDB")
	}
}
