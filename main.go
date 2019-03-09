package main

import (
	"fmt"
	"os"
)

func main() {
	user, action, err := getArgs()
	if err != nil {
		fmt.Printf("How to use: ")
	}

	fmt.Printf("User: /u/%s, Action: %s\n", user, action)
	apiClient := redditAPIClient{}
	apiClient.getUserComments(user)
}

func getArgs() (string, string, error) {
	if len(os.Args) < 2 {
		return "", "", fmt.Errorf("Expected %v arguments but received %v", 2, len(os.Args))
	}
	return os.Args[1], os.Args[2], nil
}
