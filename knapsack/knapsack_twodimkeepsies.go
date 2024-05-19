package knapsack

type TwoDimArrayKeepsies struct{}

func (k TwoDimArrayKeepsies) Run(transactions []Transaction, maxTime int) []Transaction {
	n := len(transactions)
	dp := make([][]float64, n+1)
	for i := range dp {
		dp[i] = make([]float64, maxTime+1)
	}

	// Initialize the keep array to track which items are included
	keep := make([][]bool, n+1)
	for i := range keep {
		keep[i] = make([]bool, maxTime+1)
	}

	for i := 1; i <= n; i++ {
		for t := 1; t <= maxTime; t++ {
			// assume we'll just carry forward prior value.
			dp[i][t] = dp[i-1][t]
			if transactions[i-1].ProcessingTime <= t && dp[i-1][t] < dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount {
				//but if it fits, and it gives better value  ...
				dp[i][t] = dp[i-1][t-transactions[i-1].ProcessingTime] + transactions[i-1].Amount
				keep[i][t] = true
			}
		}
	}

	// Find the selected transactions
	res := []Transaction{}
	t := maxTime
	for i := n; i >= 0; i-- {
		if keep[i][t] {
			res = append(res, transactions[i-1])
			t -= transactions[i-1].ProcessingTime
		}
	}
	return res
}
