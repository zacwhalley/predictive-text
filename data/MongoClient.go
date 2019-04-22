package data

import (
	"context"
	"errors"
	"log"
	"sort"
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
func NewMongoClient(uri string) *MongoClient {
	newClient := &MongoClient{}
	options := options.Client().ApplyURI(uri)

	var err error
	newClient.client, err = mongo.Connect(context.TODO(), options)
	if err != nil {
		log.Fatal(err)
	}

	return newClient
}

// GetChain gets the chain for the user userName
func (m MongoClient) GetChain(users []string) (*UserChainDao, error) {
	if m.client == nil {
		return nil, errors.New("No connection to MongoDB")
	}

	sort.Strings(users)

	chains := m.client.Database("redditSim").Collection("chain")
	filter := bson.D{{Key: "users", Value: users}}
	options := &options.FindOneOptions{}
	result := &UserChainDao{}

	findResult := chains.FindOne(context.TODO(), filter, options)
	if err := findResult.Err(); err != nil {
		return nil, err
	}

	err := findResult.Decode(result)
	if err != nil {
		// No document was found
		return nil, err
	}

	return result, nil
}

// UpsertChain upserts the chain for a user
func (m MongoClient) UpsertChain(users []string, chain *markov.Chain) error {
	if m.client == nil {
		return errors.New("No connection to MongoDB")
	}

	sort.Strings(users)

	// Get chain collection from redditSim db
	chains := m.client.Database("redditSim").Collection("chain")
	userChain := UserChainDao{users, chain, time.Now()}

	// Insert chain as new document
	filter := bson.D{{Key: "users", Value: users}}
	update := bson.D{{Key: "$set", Value: userChain}}
	isUpsert := true
	options := &options.UpdateOptions{Upsert: &isUpsert}

	_, err := chains.UpdateOne(context.TODO(), filter, update, options)
	if err != nil {
		return err
	}

	return nil
}
