package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zacwhalley/predictive-text/dto"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

	return nil
}

// GetPrediction returns an array of predictions for the given
// input text
func GetPrediction(w http.ResponseWriter, r *http.Request) {
	log.Print("POST /prediction")
	// Decode request body
	var predSource dto.PredictionRequestDto
	if err := json.NewDecoder(r.Body).Decode(&predSource); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create prediction
	const numWords = 2 // default number of extra words to predict
	predictions, err := predict(predSource.Input, numWords)
	if err != nil {
		log.Print(err)
		return
	}

	// Send response
	response := dto.PredictionResponseDto{
		Input:       predSource.Input,
		Predictions: predictions,
	}
	if err = respondWithJSON(w, http.StatusOK, response); err != nil {
		log.Fatal(err)
	}
	log.Print("Prediction returned")
}
