package knapsack

import "fmt"

type TwoDimArrayNoKeepsies struct{}

func (k TwoDimArrayNoKeepsies) Run(transactions []Transaction, maxTime int) []Transaction {
	n := len(transactions)
	fmt.Printf("looking at %d transactions, and a maxTime of %d ms ...\n", n, maxTime)

	dp := make([][]float64, n+1)
	for i := range dp {
		dp[i] = make([]float64, maxTime+1)
	}
	//fmt.Println("dp initially: %+v", dp)

	// initial implementation suggestion did not do a keep array, it did it
	// Initialize the keep array to track which items are included
	keep := make([][]bool, n+1)
	for i := range keep {
		keep[i] = make([]bool, maxTime+1)
	}

	for i := 1; i <= n; i++ {
		for t := 1; t <= maxTime; t++ {
			if transactions[i-1].ProcessingTime <= t {
				dp[i][t] = max(dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount)

				///not normally doing this lol
				if dp[i-1][t] < dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount {
					//fmt.Printf("KEEP i: %d and t: %d: choose %.2f vs %.2f (= %.2f + %.2f) -- we chose %.2f\n", i, t, dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount, dp[i-1][t-transactions[i-1].ProcessingTime], transactions[i-1].Amount, dp[i][t])
					keep[i][t] = true
				}

				//fmt.Printf("i: %d and t: %d: choose %.2f vs %.2f (= %.2f + %.2f) -- we chose %.2f\n", i, t, dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount, dp[i-1][t-transactions[i-1].ProcessingTime], transactions[i-1].Amount, dp[i][t])
				//fmt.Printf("fookin wot: %+v\n", dp[i-1])
				//fmt.Printf("fookin wot2: %+v\n", dp[i-1][t-transactions[i-1].ProcessingTime])
				//panic("just stop")
			} else {
				//fmt.Printf("i: %d and t: %d and we had to skip b/c processingtime for this tx is: %d\n", i, t, transactions[i-1].ProcessingTime)
				dp[i][t] = dp[i-1][t]
			}
		}
		//fmt.Printf("dp after populating for tx #%d: %+v\n\n", i, dp)
	}
	//for i, v := range dp {
	//	fmt.Printf("dp row %d after populating: %+v\n", i, v)
	//}
	//for i, tx := range transactions {
	//	fmt.Printf("txid: %s -> keep timeslots: %+v\n", tx.ID, keep[i+1])
	//}

	// Find the selected transactions
	res := []Transaction{}
	t := maxTime
	for i := n; i > 0; i-- {
		//fmt.Printf("i: %d, t: %d\n", i, t)
		if dp[i][t] != dp[i-1][t] {
			res = append(res, transactions[i-1])
			t -= transactions[i-1].ProcessingTime
			//	fmt.Printf("appended selected tx, t is now: %d\n", t)
			//} else {
			//	fmt.Printf("did nothing because %.2f == %.2f\n", dp[i][t], dp[i-1][t])
		}
	}
	return res
	//_ = res
	//
	//t = maxTime
	//altRes := []Transaction{}
	//for i := n; i >= 0; i-- {
	//	//fmt.Printf("i: %d t: %d, keep[i][t]: %t\n", i, t, keep[i][t])
	//	if keep[i][t] {
	//		altRes = append(altRes, transactions[i-1])
	//		t -= transactions[i-1].ProcessingTime
	//	}
	//}
	//panic("stop here")
	//return dp[n][maxTime], altRes
}
