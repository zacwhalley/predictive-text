package data

// DBClient is an interface for database access
type DBClient interface {
	GetChain(users []string) (*UserChainDao, error)
	GetPrediction(input string, n int) ([]string, error)
}
