package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ex request for comments: GET http://www.reddit.com/user/USERNAME/comments.json
// can get text from comment.data.body
// no auth required

type redditApiClient struct {
}

func (r redditApiClient) getUserComments(username string) {
	url := fmt.Sprintf("https://www.reddit.com/user/%s/comments.json", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "zw-reddit-simulator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
