package data

import (
	"context"
	"log"
	"time"

	"github.com/zacwhalley/reddit-simulator/markov"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
func (m MongoClient) GetChain(userName string) *UserChainDao {
	if m.client == nil {
		log.Fatal("No connection to MongoDB")
	}

	return nil
}

// UpsertChain upserts the chain for a user
func (m MongoClient) UpsertChain(userName string, chain *markov.Chain) {
	if m.client == nil {
		log.Fatal("No connection to MongoDB")
	}

	// Get chain collection from redditSim db
	chains := m.client.Database("redditSim").Collection("chain")

	// Insert chain as new document
	userChain := UserChainDao{userName, chain, time.Now()}
	_, err := chains.InsertOne(context.TODO(), userChain)
	if err != nil {
		log.Fatal(err)
	}
}
