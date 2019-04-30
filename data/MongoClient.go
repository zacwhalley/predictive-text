package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/zacwhalley/predictive-text/markov"
	"github.com/zacwhalley/predictive-text/util"
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

// GetPrediction predicts the 3 most likely next n words for an input
func (m MongoClient) GetPrediction(input string, n int) ([]string, error) {
	const numPredictions = 1
	const prefixLen = 2
	predictions := make([]string, numPredictions)
	var p markov.Prefix = make([]string, prefixLen)

	inputWords := strings.Split(input, " ")

	// put last prefixLen words of input into prefix
	copy(p, inputWords[util.MaxInt(0, len(inputWords)-prefixLen):])

	for i := 0; i < numPredictions; i++ {
		var newPrefix markov.Prefix = make([]string, len(p))
		copy(newPrefix, p)
		prediction, err := m.getMostCommon(newPrefix, n)
		if err != nil {
			return nil, err
		}
		predictions[i] = strings.TrimSpace(prediction)
		if err != nil {
			return nil, err
		}
	}

	return predictions, nil
}

// getMostCommon recursively returns the most common string in the set of strings p maps to
func (m MongoClient) getMostCommon(p markov.Prefix, n int) (string, error) {
	// recursive base case
	if n == 0 {
		return "", nil
	}

	chains := m.client.Database("redditSim").Collection("chain")
	filter := bson.D{{Key: "users", Value: []string{"predict"}}}
	options := &options.FindOneOptions{}
	options.Projection = bson.D{{Key: "chain.chain." + p.ToString(), Value: 1}}
	result := &UserChainDao{}

	findResult := chains.FindOne(context.TODO(), filter, options)
	if err := findResult.Err(); err != nil {
		return "", err
	}
	err := findResult.Decode(result)
	if err != nil {
		return "", err
	}

	// gets the list of suffixes for the given prefix
	if result.Chain == nil ||
		result.Chain.Chain == nil ||
		len(result.Chain.Chain) < 1 {
		// prefix exists but no suffix available - return early
		return "", nil
	}

	nextOptions := result.Chain.Chain[p.ToString()]
	current := nextOptions[0]
	p.Shift(current)

	// recursively find next most common
	next, err := m.getMostCommon(p, n-1)
	if err != nil {
		return "", err
	}

	return current + " " + next, nil
}

// UpsertChain upserts the chain for a user
func (m MongoClient) UpsertChain(users []string, chain *markov.Chain) error {
	if m.client == nil {
		return errors.New("No connection to MongoDB")
	}

	sort.Strings(users)

	fmt.Printf("Saving data for %v\n", users)
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
