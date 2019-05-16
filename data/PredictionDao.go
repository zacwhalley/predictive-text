package data

import (
	"github.com/zacwhalley/predictivetext/common"
)

// PredictionDao is the data access object / schema for a prediction
type PredictionDao struct {
	Source   string
	Prefix   string
	Suffixes []common.Pair
}
