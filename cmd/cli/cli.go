package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
	"github.com/zacwhalley/predictivetext/markov"
	"github.com/zacwhalley/predictivetext/util"
)

func initApp(app *cli.App) {
	setInfo(app)
	setCommands(app)
}

func setInfo(app *cli.App) {
	app.Name = "Predictive text data builder"
	app.Usage = "A CLI for generating data used for predictive text generation"
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
				cli.StringFlag{
					Name:  "source",
					Value: reddit.String(),
				},
			},
			Action: func(c *cli.Context) error {
				return buildAction(c)
			},
		},
	}
}

func buildAction(c *cli.Context) error {
	source := c.String("source")
	if source == reddit.String() {
		// Generate data from scraping reddit comments
		pageLimit := c.Int("pageLimit")
		if pageLimit < 0 {
			return errors.New("pageLimit must be greater than 0")
		}
		users := readUsers()
		log.Println("Done getting user names. Please wait for data to generate.")
		if err := buildChainFromReddit(users, pageLimit); err != nil {
			return err
		}
		return nil
	} else if source == text.String() {
		if err := buildChainFromStdin(); err != nil {
			return err
		}
	} else {
		return errors.New(source + " is not a valid data source.")
	}

	return nil
}

// readUsers reads an arbitrary number of usernames from standard input
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
		log.Printf("Getting data for user %s\n", username)
		// pageLimit <= 0 means no limit has been specified
		for i := 0; i < pageLimit || pageLimit <= 0; i++ {
			page, pageRef = api.getUserComments(username, pageRef)
			userComments = append(userComments, page)
			if pageRef == "" {
				break
			}
		}
		log.Printf("Done getting data for user %s\n", username)
		comments <- userComments
		done <- true
	}
}

// getAllComments gets up to pageLimit comments for each user in users and
// passes it to the comments channel
func getAllComments(users []string, pageLimit int) <-chan [][]string {
	comments := make(chan [][]string, 100)
	go (func() {
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
		log.Println("Done getting data for all users")
		close(comments)
	})()

	return comments
}

func buildChainFromReddit(users []string, pageLimit int) error {
	chain := markov.NewChain(2)
	for commentSet := range getAllComments(users, pageLimit) {
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
		log.Println("Save successful.")
	}
	return err
}

func buildChainFromStdin() error {
	reader := bufio.NewReader(os.Stdin)
	chain := markov.NewChain(2)

	// Generate
	chain.Build(reader)
	log.Printf("Chain generated")

	// Save
	if err := db.UpsertChain([]string{}, chain); err != nil {
		return err
	}

	return nil
}
