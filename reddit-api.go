package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zacwhalley/predictive-text/dto"
)

type redditAPIClient struct {
}

// getUserComments returns an array of the comments by username on page, and a reference to the next page
func (r redditAPIClient) getUserComments(username string, pageRef string) ([]string, string) {
	// make request to /u/username's comments
	url := fmt.Sprintf("https://www.reddit.com/user/%s/comments.json", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "zw-predictive-text")
	// add requested page to query params of url
	if pageRef != "" {
		params := req.URL.Query()
		params.Add("after", pageRef)
		req.URL.RawQuery = params.Encode()
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Could not get user %s\n", username)
	}

	// decode response and convert json objects to simple array of comments
	var page dto.CommentsPageDto
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		panic(err)
	}

	comments := make([]string, len(page.Data.Children))
	for i, comment := range page.Data.Children {
		comments[i] = comment.Data.Body
	}

	return comments, page.Data.After
}
