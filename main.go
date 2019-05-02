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

	// API
	r.HandleFunc("/api/prediction", PredictionController)

	// web content
	r.PathPrefix("/").
		Handler(http.FileServer(http.Dir("./static/"))).
		Methods(http.MethodGet)

	// start server
	port := os.Getenv("PORT")
	if strings.TrimSpace(port) == "" {
		log.Fatal("Port must be set")
	}
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
