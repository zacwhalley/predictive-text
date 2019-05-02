package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zacwhalley/predictive-text/data"
)

var db data.DBClient = data.NewMongoClient(os.Getenv("PREDTEXT_MONGODB_URI"))

func main() {
	r := mux.NewRouter()

	// API
	r.HandleFunc("/api/prediction", PredictionController)

	// web content
	r.PathPrefix("/").
		Handler(http.FileServer(http.Dir("./static/"))).
		Methods(http.MethodGet)

	// start server
	port := ":" + os.Getenv("PREDTEXT_PORT")
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}
