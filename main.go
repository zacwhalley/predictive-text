package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zacwhalley/reddit-simulator/markov"
)

func main() {
	//get user from input
	user, wordCount, pageLimit, err := getArgs()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Printf("How to use: reddit-simulator userName wordCount pageLimit (optional)")
	}

	// get comments from user
	comments := make(chan []string)
	go getAllComments(comments, user, pageLimit)

	totalComments := 0

	// build chain from comments got
	rand.Seed(time.Now().UnixNano())
	chain := markov.NewChain(2)
	for page := range comments {
		totalComments += len(page)
		for _, comment := range page {
			reader := strings.NewReader(comment)
			chain.Build(reader)
		}
	}

	// var db data.DBClient = data.MongoClient{}

	fmt.Println(chain.Generate(wordCount))
}

// getArgs
func getArgs() (string, int, int, error) {
	numArgs := len(os.Args)
	if numArgs < 3 {
		return "", -1, -1, fmt.Errorf("Expected %v or more arguments but received %v",
			2, len(os.Args)-1)
	}

	wordCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return "", -1, -1, fmt.Errorf("Length argument must be an integer")
	}

	var pageLimit int
	if numArgs > 3 {
		pageLimit, err = strconv.Atoi(os.Args[3])
		if err != nil {
			return "", -1, -1, fmt.Errorf("Page limit argument must be an integer")
		}
	}

	return os.Args[1], wordCount, pageLimit, nil
}

// getAllComments makes requests to all (or pageLimit) pages of comments
// and sends them to the comments channel
func getAllComments(comments chan<- []string, user string, pageLimit int) {
	api := redditAPIClient{}
	pageRef := ""
	var page []string
	// loop through all pages if no pagelimit specified
	for i := 0; i < pageLimit || pageLimit <= 0; i++ {
		page, pageRef = api.getUserComments(user, pageRef)
		comments <- page
		if pageRef == "" {
			break
		}
	}
	close(comments)
}
