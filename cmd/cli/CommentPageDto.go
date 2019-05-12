package main

// CommentsPageDto is the data transfer object for holding a user's comments
type CommentsPageDto struct {
	Data struct {
		Children []struct {
			Data struct {
				Body string `json:"body"`
			} `json:"data"`
		} `json:"children"`
		After string `json:"after"`
	} `json:"data"`
}
