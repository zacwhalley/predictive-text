package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli"
	"github.com/zacwhalley/reddit-simulator/markov"
)

func initApp(app *cli.App) {
	setInfo(app)
	setCommands(app)
}

func setInfo(app *cli.App) {
	app.Name = "Reddit Simulator"
	app.Usage = "A CLI for generating comments for Reddit users"
	app.Author = "Zac Whalley"
	app.Version = "1.0.0"
}

func setCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Build the data needed to generate comments for a user",
			Action: func(c *cli.Context) error {
				var err error
				if err = validateArgs(1, c); err != nil {
					return err
				}

				username := c.Args().Get(0)
				pageLimit := 0
				if c.NArg() > 1 {
					// pageLimit is optional
					pageLimit, err = strconv.Atoi(c.Args().Get(1))
					if err != nil {
						return err
					} else if pageLimit < 0 {
						return errors.New("pageLimit must be greater than 0")
					}
				}

				if err = build(username, pageLimit); err != nil {
					return err
				}

				fmt.Printf("Data for %s successfully generated.\n", username)
				return nil
			},
		},
		{
			Name:    "write",
			Aliases: []string{"w"},
			Usage:   "Write a comment for a specified user",
			Action: func(c *cli.Context) error {
				if err := validateArgs(2, c); err != nil {
					return err
				}

				username := c.Args().Get(0)
				length, err := strconv.Atoi(c.Args().Get(1))
				if err != nil {
					return err
				} else if length < 0 {
					return errors.New("length must be greater than 0")
				}

				err = write(username, length)
				return err
			},
		},
	}
}

func validateArgs(argNum int, c *cli.Context) error {
	if c.NArg() < argNum {
		return fmt.Errorf("Expected %v arguments but received %v", 2, c.NArg())
	}

	return nil
}

func build(username string, pageLimit int) error {
	// get comments from username
	comments := make(chan []string)
	go getAllComments(comments, username, pageLimit)

	// build chain from comments got
	chain := markov.NewChain(2)
	for page := range comments {
		for _, comment := range page {
			reader := strings.NewReader(comment)
			chain.Build(reader)
		}
	}

	// Save chain for fast lookup later
	err := db.UpsertChain(username, chain)
	return err
}

func write(username string, length int) error {
	chainResult, err := db.GetChain(username)
	if chainResult == nil {
		return fmt.Errorf(`User does not exist. Please `+
			`run "build %s pageLimit" to generate data for the user`, username)
	} else if err != nil {
		return err
	}

	chain := chainResult.Chain

	// Generate text from the chain
	rand.Seed(time.Now().UnixNano())
	fmt.Println(chain.Generate(length))

	return nil
}

// getAllComments makes requests to all (or pageLimit) pages of comments
// and sends them to the comments channel
func getAllComments(comments chan<- []string, user string, pageLimit int) {
	api := redditAPIClient{}
	pageRef := ""
	var page []string
	// loop through all pages if no pagelimit specified
	for i := 0; i < pageLimit || pageLimit <= 0; i++ {
		page, pageRef = api.getUserComments(user, pageRef)
		comments <- page
		if pageRef == "" {
			break
		}
	}
	close(comments)
}
