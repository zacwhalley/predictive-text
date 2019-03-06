package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ex request for comments: GET http://www.reddit.com/user/USERNAME/comments.json
// can get text from comment.data.body
// no auth required

// Create class for getting all comments

func sampleRequest() {
	url := "http://www.reddit.com/user/gallowboob/comments.json"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
