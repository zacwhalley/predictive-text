package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
	"github.com/zacwhalley/predictive-text/markov"
)

// initApp initializes the app with commands and info
func initApp(app *cli.App) {
	setInfo(app)
	setCommands(app)
}

// setInfo sets basic info for use in the help command
func setInfo(app *cli.App) {
	app.Name = "Reddit Simulator"
	app.Usage = "A CLI for generating comments for Reddit users"
	app.Author = "Zac Whalley"
	app.Version = "1.0.0"
}

// setCommands defines the commands build, write, and predict
// and their options
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
				if err := build(users, pageLimit); err != nil {
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
		{
			Name:    "predict",
			Aliases: []string{"p"},
			Usage:   "Output the input text combined with the n words most likely to come next",
			Action: func(c *cli.Context) error {
				// get input from stdin
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter text: ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				input = strings.TrimSpace(input)

				predicted, err := postPredictionSvc(input)
				if err != nil {
					return err
				}

				for _, words := range predicted.Predictions {
					fmt.Printf("%s %s\n", input, words)
				}

				return nil
			},
		},
	}
}

// validateArgs checks that the correct number of arguments has been given
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
