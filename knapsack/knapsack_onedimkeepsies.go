package knapsack

type OneDimArrayKeepsies struct{}

func (k OneDimArrayKeepsies) Run(transactions []Transaction, maxTime int) []Transaction {
	n := len(transactions)
	dp := make([]float64, maxTime+1)
	keep := make([][]bool, n)

	// Initialize the keep array to track which items are included
	for i := range keep {
		keep[i] = make([]bool, maxTime+1)
	}

	// Fill the dp array and keep track of decisions
	for i := 0; i < n; i++ {
		for t := maxTime; t >= transactions[i].ProcessingTime; t-- {
			if dp[t] < dp[t-transactions[i].ProcessingTime]+transactions[i].Amount {
				dp[t] = dp[t-transactions[i].ProcessingTime] + transactions[i].Amount
				keep[i][t] = true
			}
		}
	}

	//Determine which ones we actually sent by walking the timeslots back. In reverse order of tx processing
	selectedTransactions := []Transaction{}
	t := maxTime
	for i := n - 1; i >= 0; i-- {
		if keep[i][t] {
			selectedTransactions = append(selectedTransactions, transactions[i])
			t -= transactions[i].ProcessingTime
		}
	}

	return selectedTransactions
}
