package data

import (
	"context"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/zacwhalley/predictivetext/common"
	"github.com/zacwhalley/predictivetext/markov"
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

// UpsertChain upserts the chain for a set of users
func (m MongoClient) UpsertChain(users []string, chain *markov.Chain) error {
	if m.client == nil {
		return errors.New("No connection to MongoDB")
	}

	sort.Strings(users)

	log.Printf("Saving data for %v\n", users)
	// Get chain collection from redditSim db
	chains := m.client.Database("redditSim").Collection("chain")
	userChain := UserChainDao{users, chain, time.Now()}

	// Insert chain as new document
	filter := bson.D{{Key: "users", Value: users}}
	update := bson.D{{Key: "$set", Value: userChain}}
	isUpsert := true
	options := &options.UpdateOptions{Upsert: &isUpsert}

	result, err := chains.UpdateOne(context.TODO(), filter, update, options)
	if err != nil {
		return err
	}

	log.Printf("ID: %v", result.UpsertedID)

	return nil
}

// GetPredictionMap returns a map of specified depth containing all words
// that may occur after the given input
func (m MongoClient) GetPredictionMap(p markov.Prefix, depth int) (common.SetMap, error) {
	var newPrefix markov.Prefix = make([]string, len(p))
	copy(newPrefix, p)
	predictionMap, err := m.getPredGraph(newPrefix, depth)
	if err != nil {
		return nil, err
	}
	return predictionMap, nil
}

// getPredGraph recursively returns a map of specified depth given a depth
func (m MongoClient) getPredGraph(p markov.Prefix, depth int) (common.SetMap, error) {
	// recursive base case
	if depth <= 0 {
		return nil, nil
	}

	// Search for current prefix
	chains := m.client.Database("redditSim").Collection("chain")
	filter := bson.D{{Key: "users", Value: bson.A{}}}
	options := &options.FindOneOptions{}
	options.Projection = bson.D{{Key: "chain.chain." + p.ToString(), Value: 1}}
	result := &UserChainDao{}

	findResult := chains.FindOne(context.TODO(), filter, options)
	if err := findResult.Err(); err != nil {
		return nil, err
	}
	err := findResult.Decode(result)
	if err != nil {
		return nil, err
	}

	// gets the list of suffixes for the given prefix
	if result.Chain == nil ||
		result.Chain.Chain == nil ||
		len(result.Chain.Chain) <= 0 {
		// no suffixes exist
		return nil, nil
	}

	// Recursively shift each returned value onto the current prefix
	resultMap := result.Chain.Chain
	suffixes := resultMap[p.ToString()]

	// Merge each with the current map after it
	for suffix := range suffixes {
		var newP markov.Prefix = make([]string, len(p))
		copy(newP, p)
		newP.Shift(suffix)

		newMap, err := m.getPredGraph(newP, depth-1)
		if err != nil {
			return nil, err
		}
		resultMap.Union(newMap)
	}

	return resultMap, nil
}
