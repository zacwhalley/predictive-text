package common

import (
	"context"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/zacwhalley/predictivetext/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
	options := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		log.Fatal(err)
	}

	return &MongoClient{client}
}

// GetChainByID gets the chain associated with a specified id
func (m *MongoClient) GetChainByID(id string) (domain.UserChainDao, error) {
	if m.client == nil {
		return domain.UserChainDao{}, errors.New("No connection to MongoDB")
	}

	chains := m.client.Database("predtext").Collection("chain")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.UserChainDao{}, err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	options := &options.FindOneOptions{}
	result := &domain.UserChainDao{}

	findResult := chains.FindOne(context.TODO(), filter, options)
	if err := findResult.Err(); err != nil {
		return domain.UserChainDao{}, err
	}

	err = findResult.Decode(result)
	if err != nil {
		// No document was found
		return domain.UserChainDao{}, err
	}

	return *result, nil
}

// UpsertChain upserts the chain for a set of users
func (m MongoClient) UpsertChain(users []string, chain domain.Chain) error {
	if m.client == nil {
		return errors.New("No connection to MongoDB")
	}

	sort.Strings(users)

	log.Printf("Saving data for %v\n", users)

	data := chain.GetData().ToPrimitive()

	// Get chain collection from redditSim db
	chains := m.client.Database("predtext").Collection("chain")
	userChain := domain.UserChainDao{
		Users:        users,
		Data:         data,
		LastModified: time.Now(),
		PrefixLen:    chain.GetPrefixLen(),
	}

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

// GetPrediction returns a prediction for the given prefix and source
func (m MongoClient) GetPrediction(prefix, source string) (domain.Prediction, error) {
	if m.client == nil {
		return domain.Prediction{}, errors.New("No connection to MongoDB")
	}

	predictions := m.client.Database("predtext").Collection("predictions")
	filter := bson.D{
		{Key: "prefix", Value: prefix},
		{Key: "source", Value: source},
	}
	options := &options.FindOneOptions{}
	result := &domain.PredictionDao{}

	findResult := predictions.FindOne(context.TODO(), filter, options)
	if err := findResult.Err(); err != nil {
		return domain.Prediction{}, err
	}

	err := findResult.Decode(result)
	if err != nil {
		// No document was found
		return domain.Prediction{}, err
	}

	predictionResult := domain.Prediction{}
	if findResult != nil {
		predictionResult.Prefix = result.Prefix
		predictionResult.Suffixes = result.Suffixes
	}

	return predictionResult, nil
}

// UpsertPrediction upserts a prediction in the prediction collection
// using the prefix as a key
func (m MongoClient) UpsertPrediction(prediction domain.Prediction) error {
	if m.client == nil {
		return errors.New("No connection to MongoDB")
	}

	// Get chain collection from redditSim db
	predictions := m.client.Database("predtext").Collection("predictions")
	document := domain.PredictionDao{
		Source:   "",
		Prefix:   prediction.Prefix,
		Suffixes: prediction.Suffixes,
	}

	// Insert chain as new document
	filter := bson.D{{Key: "prefix", Value: prediction.Prefix}}
	update := bson.D{{Key: "$set", Value: document}}
	isUpsert := true
	options := &options.UpdateOptions{Upsert: &isUpsert}

	_, err := predictions.UpdateOne(context.TODO(), filter, update, options)
	if err != nil {
		return err
	}

	return nil
}
