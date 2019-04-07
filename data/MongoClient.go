package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zacwhalley/reddit-simulator/markov"
	"go.mongodb.org/mongo-driver/bson"
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

	chains := m.client.Database("redditSim").Collection("chain")
	filter := bson.D{{Key: "user", Value: userName}}
	options := &options.FindOneOptions{}
	result := &UserChainDao{}

	findResult := chains.FindOne(context.TODO(), filter, options)
	if err := findResult.Err(); err != nil {
		log.Fatal(err)
	}

	err := findResult.Decode(result)
	if err != nil {
		// No document was found
		fmt.Printf("Returning nothing from find chain")
		return nil
	}

	return result
}

// UpsertChain upserts the chain for a user
func (m MongoClient) UpsertChain(userName string, chain *markov.Chain) {
	if m.client == nil {
		log.Fatal("No connection to MongoDB")
	}

	// Get chain collection from redditSim db
	chains := m.client.Database("redditSim").Collection("chain")
	userChain := UserChainDao{userName, chain, time.Now()}

	// Insert chain as new document
	filter := bson.D{{Key: "user", Value: userName}}
	update := bson.D{{Key: "$set", Value: userChain}}
	isUpsert := true
	options := &options.UpdateOptions{Upsert: &isUpsert}

	_, err := chains.UpdateOne(context.TODO(), filter, update, options)
	if err != nil {
		log.Fatal(err)
	}
}
