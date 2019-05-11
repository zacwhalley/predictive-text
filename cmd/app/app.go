package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/zacwhalley/predictivetext/data"
)

func main() {
	// Get environment variables - TODO: Move this into a config somewhere
	port := getEnv("PORT")
	mongodbURI := getEnv("MONGODB_URI")

	// Set up services and handlers
	db := data.NewMongoClient(mongodbURI)
	predictionSvc := PredictionSvc{db}
	predictionHandler := PredictionHandler{predictionSvc}
	demoHandler := DemoHandler{}

	r := mux.NewRouter()

	// API handling
	r.HandleFunc("/api/prediction", predictionHandler.Handle)

	// ui handling
	r.HandleFunc("/", demoHandler.Handle).
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

func getEnv(varname string) string {
	result := strings.TrimSpace(os.Getenv(varname))
	if result == "" {
		log.Fatalf("%s must be set", varname)
	}

	return result
}
