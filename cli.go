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
	"github.com/zacwhalley/reddit-simulator/util"
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
func getUserComments(comments chan<- [][]string, usernames <-chan string,
	done chan<- bool, pageLimit int) {

	var userComments [][]string
	var page []string
	api := redditAPIClient{}
	pageRef := ""
	for username := range usernames {
		fmt.Printf("Getting data for user %s\n", username)
		// pageLimit <= 0 means no limit has been specified
		for i := 0; i < pageLimit || pageLimit <= 0; i++ {
			page, pageRef = api.getUserComments(username, pageRef)
			userComments = append(userComments, page)
			if pageRef == "" {
				break
			}
		}
		fmt.Printf("Done getting data for user %s\n", username)
		comments <- userComments
		done <- true
	}
}

// getAllComments gets up to pageLimit comments for each user in users and
// passes it to the comments channel
func getAllComments(users []string, pageLimit int, comments chan<- [][]string) {
	// create an arbitrary number of workers to get the comments
	// see https://gobyexample.com/worker-pools
	usernames := make(chan string, len(users))
	numWorkers := util.MinInt(len(users), 3)
	done := make(chan bool, len(users))
	for i := 0; i < numWorkers; i++ {
		go getUserComments(comments, usernames, done, pageLimit)
	}

	for _, user := range users {
		usernames <- user
	}
	close(usernames)

	// close comments channel after data for all users has been collected
	for i := 0; i < len(users); i++ {
		<-done
	}
	fmt.Println("Done getting data for all users")
	close(comments)
}

func buildChain(users []string, pageLimit int) error {
	comments := make(chan [][]string, 100)
	getAllComments(users, pageLimit, comments)

	// build chain from comments got
	chain := markov.NewChain(2)
	for commentSet := range comments {
		for _, page := range commentSet {
			for _, comment := range page {
				reader := strings.NewReader(comment)
				chain.Build(reader)
			}
		}
	}

	// Save chain for fast lookup later
	err := db.UpsertChain(users, chain)
	if err == nil {
		fmt.Println("Save successful.")
	}
	return err
}

func write(users []string, length int, generator markov.Generator) error {
	chainResult, err := db.GetChain(users)
	if chainResult == nil {
		return fmt.Errorf(`No data found. Please run "build --pageLimit=N" ` +
			`to generate data for the users`)
	} else if err != nil {
		return err
	}

	chain := chainResult.Chain

	// Generate text from the chain
	rand.Seed(time.Now().UnixNano())
	fmt.Println(generator.Generate(chain, length))

	return nil
}
