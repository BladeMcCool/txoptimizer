package knapsack

type TwoDimArrayNoKeepsies struct{}

func (k TwoDimArrayNoKeepsies) Run(transactions []Transaction, maxTime int) []Transaction {
	n := len(transactions)
	dp := make([][]float64, n+1)
	for i := range dp {
		dp[i] = make([]float64, maxTime+1)
	}

	for i := 1; i <= n; i++ {
		for t := 1; t <= maxTime; t++ {
			if transactions[i-1].ProcessingTime <= t {
				// if it fits, decide if its better use of space.
				dp[i][t] = max(dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount)
			} else {
				// just carry forward prior value.
				dp[i][t] = dp[i-1][t]
			}
		}
	}

	// Find the selected transactions
	res := []Transaction{}
	t := maxTime
	for i := n; i > 0; i-- {
		if dp[i][t] != dp[i-1][t] {
			res = append(res, transactions[i-1])
			t -= transactions[i-1].ProcessingTime
		}
	}
	return res
}
