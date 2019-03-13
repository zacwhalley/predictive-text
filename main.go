package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//get user from input
	user, wordCount, err := getArgs()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Printf("How to use: reddit-simulator userName wordCount")
	}

	// get comments from user
	comments := make(chan []string)
	go getAllComments(comments, user)

	totalComments := 0

	// generate and print text
	rand.Seed(time.Now().UnixNano())
	chain := newChain(2)
	for page := range comments {
		totalComments += len(page)
		for _, comment := range page {
			reader := strings.NewReader(comment)
			chain.build(reader)
		}
	}

	fmt.Println(chain.generate(wordCount))
}

// getArgs
func getArgs() (string, int, error) {
	if len(os.Args) != 3 {
		return "", -1, fmt.Errorf("Expected %v arguments but received %v", 2, len(os.Args)-1)
	}

	wordCount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return "", -1, fmt.Errorf("Second argument must be an integer")
	}

	return os.Args[1], wordCount, nil
}

// getAllComments makes requests to all pages of comments and sends them to the comments channel
func getAllComments(comments chan<- []string, user string) {
	api := redditAPIClient{}
	pageRef := ""
	var page []string
	for {
		page, pageRef = api.getUserComments(user, pageRef)
		comments <- page
		if pageRef == "" {
			break
		}
	}
	close(comments)
}
