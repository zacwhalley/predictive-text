package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/zacwhalley/reddit-simulator/dtos"
)

// ex request for comments: GET http://www.reddit.com/user/USERNAME/comments.json
// can get text from comment.data.body
// no auth required

type redditAPIClient struct {
}

func (r redditAPIClient) getUserComments(username string) []string {
	// make request to /u/username's comments
	url := fmt.Sprintf("https://www.reddit.com/user/%s/comments.json", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "zw-reddit-simulator")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// decode response and convert json objects to simple array of comments
	var comments dtos.CommentsPageDto
	if err := json.NewDecoder(res.Body).Decode(&comments); err != nil {
		panic(err)
	}

	content := make([]string, len(comments.Data.Children))
	for i, comment := range comments.Data.Children {
		content[i] = comment.Data.Body
	}

	return content
}
