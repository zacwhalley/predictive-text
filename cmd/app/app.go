package app

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/zacwhalley/predictivetext/common"
)

// Main is the entrypoint for the predictivetext web service
func Main() {
	// Get environment variables - TODO: Move this into a config somewhere
	port := getEnv("PORT")
	mongodbURI := getEnv("MONGODB_URI")

	// Set up services and handlers
	db := common.NewMongoClient(mongodbURI)
	predictionSvc := common.PredictionSvc{db}
	predictionHandler := PredictionHandler{predictionSvc}
	demoHandler := DemoHandler{}

	r := mux.NewRouter()

	// API handling
	r.HandleFunc("/api/prediction", predictionHandler.Handle)

	// ui handling
	r.HandleFunc("/", demoHandler.Handle).
		Methods(http.MethodGet)

	wd, _ := os.Getwd()
	staticDir := filepath.Join(wd, "./cmd/app/static/")
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))).
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
