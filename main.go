package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zacwhalley/predictive-text/data"
)

var db data.DBClient = data.NewMongoClient("mongodb://localhost:27017")

func main() {
	/*
		Legacy CLI stuff
		app := cli.NewApp()
		initApp(app)

		err := app.Run(os.Args)
		if err != nil {
			log.Fatal(err)
		}
	*/

	r := mux.NewRouter()
	r.HandleFunc("/prediction", GetPrediction).Methods("POST")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
