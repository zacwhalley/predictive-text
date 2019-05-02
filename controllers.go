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

// PredictionController returns an array of predictions for the given
// input text
func PredictionController(w http.ResponseWriter, r *http.Request) {
	log.Print("POST /prediction")
	// Decode request body
	var predSource dto.PredictionRequestDto
	if err := json.NewDecoder(r.Body).Decode(&predSource); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create prediction
	response, err := predictionSvc(predSource.Input)
	if err != nil {
		log.Print(err)
		return
	}

	// Send response
	if err = respondWithJSON(w, http.StatusOK, response); err != nil {
		log.Fatal(err)
	}
	log.Print("Prediction returned")
}
