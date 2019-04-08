package main

import (
	"flag"
	"fmt"
	"strconv"
)

// getArgs returns the user's choice of username and sentence size
// with optional arguments pageLimit and refresh
func getArgs() (string, int, int, bool, error) {
	// Optional arguments
	pageLimPtr := flag.Int("limit", 0, "number of pages to get data from")
	refreshPtr := flag.Bool("refresh", false, "if generated chain is updated with new data (slow)")
	flag.Parse()

	// Required arguments
	numArgs := len(flag.Args())
	if numArgs < 2 {
		return "", -1, -1, false, fmt.Errorf("Expected %v arguments but received %v",
			2, len(flag.Args()))
	}

	user := flag.Arg(0)
	wordCount, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		return "", -1, -1, false, fmt.Errorf("Length argument must be an integer")
	}

	return user, wordCount, *pageLimPtr, *refreshPtr, nil
}
