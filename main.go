package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Transaction represents a transaction with an ID, amount, country code, and processing time
type Transaction struct {
	ID              string
	Amount          float64
	BankCountryCode string
	ProcessingTime  int
}

func main() {
	// Sample CSV data
	csvData := `
id,amount,bank_country_code
dde3165e-a7e9-4dac-984b-4aa5f32a45e2,6.44,tr
95c6e393-e3dd-499e-b6bf-d075fc2abbc5,9.31,uk
f14819c0-145b-408b-8455-8a819ef7015a,15.21,ma
`

	// Parsing CSV data
	r := csv.NewReader(strings.NewReader(csvData))
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Convert CSV data to transactions
	transactions := make([]Transaction, 0, len(records)-1)
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, Transaction{
			ID:              record[0],
			Amount:          amount,
			BankCountryCode: record[2],
		})
	}

	// Sample country code to processing time mapping
	processingTimes := map[string]int{
		"tr": 30, // 30 milliseconds
		"uk": 20, // 20 milliseconds
		"ma": 50, // 50 milliseconds
	}

	// Add processing time to each transaction
	for i := range transactions {
		transactions[i].ProcessingTime = processingTimes[transactions[i].BankCountryCode]
	}

	// Example usage
	maxTime := 50 // in milliseconds
	maxAmount, selectedTransactions := knapsack(transactions, maxTime)

	fmt.Printf("Max amount: %.2f\n", maxAmount)
	fmt.Println("Selected transactions:")
	for _, transaction := range selectedTransactions {
		fmt.Println(transaction)
	}
}

// Knapsack problem solution
func knapsack(transactions []Transaction, maxTime int) (float64, []Transaction) {
	n := len(transactions)
	dp := make([][]float64, n+1)
	for i := range dp {
		dp[i] = make([]float64, maxTime+1)
	}

	for i := 1; i <= n; i++ {
		for t := 1; t <= maxTime; t++ {
			if transactions[i-1].ProcessingTime <= t {
				dp[i][t] = max(dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount)
			} else {
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

	return dp[n][maxTime], res
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
