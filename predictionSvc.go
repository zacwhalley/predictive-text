package main

import (
	"sort"
	"strings"

	"github.com/zacwhalley/predictive-text/dto"
	"github.com/zacwhalley/predictive-text/markov"
	"github.com/zacwhalley/predictive-text/util"
)

type prediction struct {
	text   string
	weight int
}

// predictionSvc predicts the most likely next words for an input
func predictionSvc(input string) (dto.PredictionResponseDto, error) {
	// create prefix for search
	const prefixLen = 2

	// put last prefixLen words of input into prefix
	inputWords := strings.Split(input, " ")
	prefixWords := inputWords[util.MaxInt(0, len(inputWords)-prefixLen):]

	const depth = 2 // default to search for next 2 most common
	result, err := db.GetPredictionMap(prefixWords, depth)
	if err != nil {
		return dto.PredictionResponseDto{}, err
	}

	freqSet := buildFreqSet(result, depth, prefixWords, "")
	predictions := makePredictions(freqSet)

	response := dto.PredictionResponseDto{
		Input:       input,
		Predictions: predictions,
	}

	return response, nil
}

func buildFreqSet(chain map[string][]string, depth int, p markov.Prefix, result string) Set {
	var resultSet Set = make(map[string]int)
	suffixes := chain[p.ToString()]

	if depth == 0 || len(suffixes) == 0 {
		// base case - add result to set and return
		trimmedResult := strings.TrimSpace(result)
		if trimmedResult != "" {
			resultSet.AddString(trimmedResult)
		}
		return resultSet
	}

	for _, suffix := range suffixes {
		// copy prefix
		var newP markov.Prefix = make([]string, len(p))
		copy(newP, p)
		newP.Shift(suffix)

		// recurse, get map with depth reduced by 1 - merge with current map
		nextSet := buildFreqSet(chain, depth-1, newP, result+" "+suffix)
		resultSet.AddSet(nextSet)
	}

	return resultSet
}

func makePredictions(freqSet Set) []string {
	// sort by weight & return top n
	predictions := []prediction{}
	for text, weight := range freqSet {
		predictions = append(predictions, prediction{text: text, weight: weight})
	}
	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].weight > predictions[j].weight
	})

	// create slice of top numPredictions results
	numPredictions := util.MinInt(3, len(predictions))
	predictionRes := make([]string, numPredictions)
	for i, pred := range predictions[:numPredictions] {
		predictionRes[i] = pred.text
	}

	return predictionRes
}
