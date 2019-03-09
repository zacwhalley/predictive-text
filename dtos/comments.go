package dtos

// CommentsDto is the data transfer object for holding a user's comments
type CommentsDto struct {
	Data struct {
		Children []struct {
			Data struct {
				Body string `json:"body"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}
