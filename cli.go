package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
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
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "pageLimit",
					Value: 0,
				},
			},
			Action: func(c *cli.Context) error {
				pageLimit := c.Int("pageLimit")
				if pageLimit < 0 {
					return errors.New("pageLimit must be greater than 0")
				}

				users := readUsers()

				fmt.Println("Done getting user names. Please wait for data to generate.")
				if err := buildChain(users, pageLimit); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:    "write",
			Aliases: []string{"w"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "words",
					Usage: "Generate specified number of words instead of sentences",
				},
				cli.StringFlag{
					Name:  "startWith",
					Usage: "Begin the generation with a given phrase",
					Value: "",
				},
			},
			Usage: "Write a comment for a specified user",
			Action: func(c *cli.Context) error {
				if err := validateArgs(1, c); err != nil {
					return err
				}

				// Get required arguments
				users := readUsers()
				length, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					return err
				} else if length < 0 {
					return errors.New("length must be greater than 0")
				}

				// Optional arguments
				beginning := strings.TrimSpace(c.String("startWith"))

				var generator markov.Generator
				if c.Bool("words") {
					generator = markov.WordGenerator{Beginning: beginning}
				} else {
					generator = markov.SentenceGenerator{Beginning: beginning}
				}
				err = write(users, length, generator)
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

// readUserNames reads an arbitrary number of usernames from standard input
func readUsers() []string {
	var users []string
	var line string

	// prompt user
	fmt.Println("Please enter the names of the users you want to use.")

	// read all words from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line = scanner.Text()
		for _, word := range strings.Fields(line) {
			users = append(users, word)
		}
	}

	return users
}

// getUserComments makes requests to all (or pageLimit) pages of comments
// and sends them to the comments channel
func getUserComments(comments chan<- []string, user string, pageLimit int) {
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
}

// getAllComments gets up to pageLimit comments for each user in users and
// passes it to the comments channel
func getAllComments(users []string, pageLimit int, comments chan []string) {
	for _, user := range users {
		// get comments from user
		getUserComments(comments, user, pageLimit)
		fmt.Printf("Data for %s found.\n", user)
	}
	close(comments)
}

func buildChain(users []string, pageLimit int) error {
	comments := make(chan []string)
	go getAllComments(users, pageLimit, comments)

	// build chain from comments got
	chain := markov.NewChain(2)
	for page := range comments {
		for _, comment := range page {
			reader := strings.NewReader(comment)
			chain.Build(reader)
		}
	}

	// Save chain for fast lookup later
	err := db.UpsertChain(users, chain)
	if err == nil {
		fmt.Println("Data for all users has been generated.")
	}
	return err
}

func write(users []string, length int, generator markov.Generator) error {
	chainResult, err := db.GetChain(users)
	if chainResult == nil {
		return fmt.Errorf(`User does not exist. Please `+
			`run "build %v --pageLimit=N" to generate data for the users`, users)
	} else if err != nil {
		return err
	}

	chain := chainResult.Chain

	// Generate text from the chain
	rand.Seed(time.Now().UnixNano())
	fmt.Println(generator.Generate(chain, length))

	return nil
}
