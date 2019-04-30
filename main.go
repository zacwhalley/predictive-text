package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/zacwhalley/predictive-text/data"
)

var db data.DBClient = data.NewMongoClient(os.Getenv("PREDTEXT_MONGODB_URI"))

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

	port := ":" + os.Getenv("PREDTEXT_PORT")
	err := http.ListenAndServe(port,
		handlers.CORS(allowedOrigins, allowedMethods)(r))
	if err != nil {
		log.Fatal(err)
	}
}
