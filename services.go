package main

import (
	"github.com/zacwhalley/predictive-text/dto"
)

// predictionSvc predicts the 3 most likely next n words for an input
// implementation for the print command
func predictionSvc(input string) (dto.PredictionResponseDto, error) {
	const numWords = 2 // default to <=2 words to predict
	result, err := db.GetPrediction(input, numWords)
	if err != nil {
		return dto.PredictionResponseDto{}, err
	}

	response := dto.PredictionResponseDto{
		Input:       input,
		Predictions: result,
	}

	return response, nil
}
