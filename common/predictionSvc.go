package common

import (
	"log"
	"sort"

	"github.com/zacwhalley/predictivetext/domain"
	"github.com/zacwhalley/predictivetext/util"
)

// PredictionSvc is an implementation of the PredictionSvc interface
type PredictionSvc struct {
	DB domain.DBClient
}

// GetPrediction predicts the most likely next words for an input
func (svc PredictionSvc) GetPrediction(input string) ([]string, error) {
	prediction, err := svc.DB.GetPrediction(input, "")
	if err != nil {
		return nil, err
	}

	suffixes := make([]string, 0)
	for _, suffix := range prediction.Suffixes {
		suffixes = append(suffixes, suffix.Key)
	}

	return suffixes, nil
}

// SavePrediction saves a prediction to the db
func (svc PredictionSvc) SavePrediction(prediction domain.Prediction) error {
	err := svc.DB.UpsertPrediction(prediction)
	return err
}

// GeneratePredictionSet builds the prediction set for a markov chain
func (svc PredictionSvc) GeneratePredictionSet(id string) error {
	chaindao, err := svc.DB.GetChainByID(id)
	if err != nil {
		return err
	}
	chainData := MakeSetMap(chaindao.Data)
	chain := Chain{
		data:      chainData,
		prefixLen: chaindao.PrefixLen,
	}

	const depth = 2 // arbitrary
	for prefix := range chainData {
		prediction := predictionFromChain(prefix, chain, depth)
		if err := svc.DB.UpsertPrediction(prediction); err != nil {
			return err
		}
	}

	log.Printf("All predictions saved for id %s", id)
	return nil
}

func predictionFromChain(key string, chain domain.Chain, depth int) domain.Prediction {
	prefixLen := chain.GetPrefixLen()
	prefix := MakePrefix(key, prefixLen)
	suffixes := getFollowSet(prefix, chain.GetData(), depth).ToPairs()

	// sort in descending order + return top 3
	sort.Slice(suffixes, func(i, j int) bool {
		return suffixes[i].Value > suffixes[j].Value
	})

	prediction := domain.Prediction{
		Prefix:   prefix.ToString(),
		Suffixes: suffixes[:util.MinInt(3, len(suffixes))],
	}

	return prediction
}

func getFollowSet(prefix domain.Prefix, chain domain.SetMap, depth int) domain.Set {
	results := make(Set)
	suffixSet, _ := chain.Get(prefix.ToString())
	suffixMap, _ := suffixSet.(Set)
	if len(suffixMap) <= 0 || depth == 0 {
		// base case
		return results
	}

	for suffix := range suffixMap {
		newPrefix := prefix.Copy()
		newPrefix.Shift(suffix)

		next := getFollowSet(newPrefix, chain, depth-1)

		results.Union(next)
	}

	return results
}
