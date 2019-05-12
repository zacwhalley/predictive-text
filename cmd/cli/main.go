package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/zacwhalley/predictivetext/data"
)

var db data.DBClient = data.NewMongoClient("mongodb://localhost:27017")

func main() {
	app := cli.NewApp()
	initApp(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
