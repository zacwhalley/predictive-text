package app

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zacwhalley/predictivetext/domain"
)

// PredictionHandler handles requests for predictions
type PredictionHandler struct {
	PredictionSvc domain.PredictionSvc
}

// DemoHandler handles requests for the demo page
type DemoHandler struct{}

// Handle handles requests for predictions
func (handler PredictionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var predSource domain.PredictionRequest
	if err := json.NewDecoder(r.Body).Decode(&predSource); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create predictions
	predictions, err := handler.PredictionSvc.Predict(predSource.Input)
	if err != nil {
		log.Print(err)
		http.Error(w, "Could not create prediction", http.StatusInternalServerError)
		return
	}

	response := domain.PredictionResponse{
		Input:       predSource.Input,
		Predictions: predictions,
	}

	// Send response
	if err = respondWithJSON(w, http.StatusOK, response); err != nil {
		log.Print(err)
		http.Error(w, "Error returning prediction", http.StatusInternalServerError)
	}
}

// Handle handles requests for the demo page
func (handler DemoHandler) Handle(w http.ResponseWriter, r *http.Request) {
	wd, _ := os.Getwd()
	fileName := filepath.Join(wd, "./cmd/app/templates/demo.html")
	t, err := template.ParseFiles(fileName)
	if err != nil {
		log.Print(err)
		respondWithJSON(w, http.StatusInternalServerError, nil)
		return
	}
	data := struct{ APIUrl string }{os.Getenv("API_URL")}

	t.Execute(w, data)
}

// respondWithJSON converts the payload data to JSON and returns it
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
