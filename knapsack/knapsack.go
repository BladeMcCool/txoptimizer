package knapsack

// Transaction represents a transaction with an ID, amount, country code, and processing time
type Transaction struct {
	ID              string
	Amount          float64
	BankCountryCode string
	ProcessingTime  int
	TxNum           int
}

type KnapsackAlgorithm interface {
	Run(transactions []Transaction, maxTime int) []Transaction
}

// KnapsackPrioritizer prioritizes transactions using a given knapsack algorithm implementation
type KnapsackPrioritizer struct {
	Algorithm KnapsackAlgorithm
}

func (kp KnapsackPrioritizer) Run(transactions []Transaction, maxTime int) []Transaction {
	// Use the algorithm to run the knapsack
	return kp.Algorithm.Run(transactions, maxTime)
}
