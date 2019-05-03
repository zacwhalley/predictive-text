package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/zacwhalley/predictive-text/data"
)

var db data.DBClient = data.NewMongoClient(os.Getenv("MONGODB_URI"))

func main() {
	r := mux.NewRouter()

	// Get environment variables
	api := os.Getenv("API_URL")
	if strings.TrimSpace(api) == "" {
		log.Fatal("API_URL must be set")
	}

	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		log.Fatal("PORT must be set")
	}

	// API
	r.HandleFunc("/api/prediction", PredictionController)

	// ui

	r.HandleFunc("/", serveDemo).
		Methods(http.MethodGet)

	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))).
		Methods(http.MethodGet)

	// start server
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
