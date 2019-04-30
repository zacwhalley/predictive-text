package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"

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
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{
		http.MethodPost,
	})

	r.HandleFunc("/prediction", PostPredictionController).Methods(http.MethodPost)

	err := http.ListenAndServe(":8080",
		handlers.CORS(allowedOrigins, allowedMethods)(r))
	if err != nil {
		log.Fatal(err)
	}
}
