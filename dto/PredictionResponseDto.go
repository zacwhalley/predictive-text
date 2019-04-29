package dto

// PredictionResponseDto is the Dto for returning a prediction
type PredictionResponseDto struct {
	Input       string   `json:"input"`
	Predictions []string `json:"predictions"`
}
