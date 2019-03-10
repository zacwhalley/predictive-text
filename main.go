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
	apiClient := redditAPIClient{}
	comments := apiClient.getUserComments(user)

	// generate and print text
	rand.Seed(time.Now().UnixNano())
	chain := newChain(1)
	for _, comment := range comments {
		reader := strings.NewReader(comment)
		chain.build(reader)
	}

	fmt.Println(chain.generate(wordCount))
}

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
