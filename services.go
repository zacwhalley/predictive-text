package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/zacwhalley/predictive-text/dto"

	"github.com/zacwhalley/predictive-text/markov"
	"github.com/zacwhalley/predictive-text/util"
)

// getUserComments makes requests to all (or pageLimit) pages of comments
// and sends them to the comments channel
func getUserComments(comments chan<- [][]string, usernames <-chan string,
	done chan<- bool, pageLimit int) {

	var userComments [][]string
	var page []string
	api := redditAPIClient{}
	pageRef := ""
	for username := range usernames {
		fmt.Printf("Getting data for user %s\n", username)
		// pageLimit <= 0 means no limit has been specified
		for i := 0; i < pageLimit || pageLimit <= 0; i++ {
			page, pageRef = api.getUserComments(username, pageRef)
			userComments = append(userComments, page)
			if pageRef == "" {
				break
			}
		}
		fmt.Printf("Done getting data for user %s\n", username)
		comments <- userComments
		done <- true
	}
}

// getAllComments gets up to pageLimit comments for each user in users and
// passes it to the comments channel
func getAllComments(users []string, pageLimit int, comments chan<- [][]string) {
	// create an arbitrary number of workers to get the comments
	// see https://gobyexample.com/worker-pools
	usernames := make(chan string, len(users))
	numWorkers := util.MinInt(len(users), 3)
	done := make(chan bool, len(users))
	for i := 0; i < numWorkers; i++ {
		go getUserComments(comments, usernames, done, pageLimit)
	}

	for _, user := range users {
		usernames <- user
	}
	close(usernames)

	// close comments channel after data for all users has been collected
	for i := 0; i < len(users); i++ {
		<-done
	}
	fmt.Println("Done getting data for all users")
	close(comments)
}

// build builds and stores the data set for users
// the implementation of the build command
func build(users []string, pageLimit int) error {
	comments := make(chan [][]string, 100)
	getAllComments(users, pageLimit, comments)

	// build chain from comments got
	chain := markov.NewChain(2)
	for commentSet := range comments {
		for _, page := range commentSet {
			for _, comment := range page {
				reader := strings.NewReader(comment)
				chain.Build(reader)
			}
		}
	}

	// Save chain for fast lookup later
	err := db.UpsertChain(users, chain)
	if err == nil {
		fmt.Println("Save successful.")
	}
	return err
}

// write prints length words/sentences generated from the data for given users
// implementation of the write command
func write(users []string, length int, generator markov.Generator) error {
	chainResult, err := db.GetChain(users)
	if chainResult == nil {
		return fmt.Errorf(`No data found. Please run "build --pageLimit=N" ` +
			`to generate data for the users`)
	} else if err != nil {
		return err
	}

	chain := chainResult.Chain

	// Generate text from the chain
	rand.Seed(time.Now().UnixNano())
	fmt.Println(generator.Generate(chain, length))

	return nil
}

// postPredictionSvc predicts the 3 most likely next n words for an input
// implementation for the print command
func postPredictionSvc(input string) (dto.PredictionResponseDto, error) {
	const numWords = 2 // default to <=2 words to predict
	result, err := db.GetPrediction(input, numWords)
	if err != nil {
		return dto.PredictionResponseDto{}, err
	}

	response := dto.PredictionResponseDto{
		Input:       input,
		Predictions: result,
	}

	return response, nil
}
