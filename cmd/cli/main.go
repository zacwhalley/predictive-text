package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/zacwhalley/predictivetext/common"
	"github.com/zacwhalley/predictivetext/domain"
)

var db domain.DBClient = common.NewMongoClient("mongodb://localhost:27017")

func main() {
	app := cli.NewApp()
	initApp(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
