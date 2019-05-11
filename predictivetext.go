package predictivetext

// PredictionSvc is a service for generating predictions
type PredictionSvc interface {
	Predict(input string) ([]string, error)
}

// PredictionRequest is the dto for requesting a prediction
type PredictionRequest struct {
	Input string `json:"input"`
}

// PredictionResponse is the Dto for returning a prediction
type PredictionResponse struct {
	Input       string   `json:"input"`
	Predictions []string `json:"predictions"`
}
