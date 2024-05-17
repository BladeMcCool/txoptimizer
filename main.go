package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

// Transaction represents a transaction with an ID, amount, country code, and processing time
type Transaction struct {
	ID              string
	Amount          float64
	BankCountryCode string
	ProcessingTime  int
	ValuePerMs      float64
}

// LatencyMapping represents the mapping of country codes to processing times
type LatencyMapping map[string]int

func main() {
	// Paths to the input files
	csvFilePath := "challenge_source/transactions.csv"
	jsonFilePath := "challenge_source/latencies.json"

	// Read the latencies from the JSON file
	latencies, err := readLatencies(jsonFilePath)
	if err != nil {
		log.Fatalf("failed to read latencies: %v", err)
	}

	// Read the transactions from the CSV file
	transactions, err := readTransactions(csvFilePath)
	if err != nil {
		log.Fatalf("failed to read transactions: %v", err)
	}

	// Add processing time and value per ms to each transaction
	for i := range transactions {
		if latency, exists := latencies[transactions[i].BankCountryCode]; exists {
			transactions[i].ProcessingTime = latency
			transactions[i].ValuePerMs = transactions[i].Amount / float64(latency)
		} else {
			log.Fatalf("no latency found for country code: %s", transactions[i].BankCountryCode)
		}
	}

	// Array of time options in milliseconds
	timeOptions := []int{50, 60, 90, 1000}
	//timeOptions := []int{50}
	for _, maxTime := range timeOptions {
		// Run the knapsack function for each time option
		runKnapsack(transactions, maxTime)
	}
}

func runKnapsack(transactions []Transaction, maxTime int) {
	maxAmount, selectedTransactions := knapsack(transactions, maxTime)

	fmt.Printf("For max time %d ms:\n", maxTime)
	fmt.Printf("Max amount: %.2f\n", maxAmount)
	fmt.Println("Selected transactions:")
	for _, transaction := range selectedTransactions {
		fmt.Println(transaction)
	}

	maxAmountGreedy, selectedTransactionsGreedy := greedyKnapsack(transactions, maxTime)

	fmt.Printf("Greedy Max amount: %.2f\n", maxAmountGreedy)
	fmt.Println("Greedy Selected transactions:")
	for _, transaction := range selectedTransactionsGreedy {
		fmt.Println(transaction)
	}

	fmt.Println("----------------------")
}

// readLatencies reads the latency mapping from a JSON file
func readLatencies(filePath string) (LatencyMapping, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var latencies LatencyMapping
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&latencies)
	if err != nil {
		return nil, err
	}

	return latencies, nil
}

// readTransactions reads the transactions from a CSV file
func readTransactions(filePath string) ([]Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	transactions := make([]Transaction, 0, len(records)-1)
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, Transaction{
			ID:              record[0],
			Amount:          amount,
			BankCountryCode: record[2],
		})
		//if i >= 10 {
		//	break
		//} //hax just work with 10 for now
	}

	return transactions, nil
}

// Knapsack problem solution
func knapsack(transactions []Transaction, maxTime int) (float64, []Transaction) {

	n := len(transactions)
	//fmt.Printf("looking at %d transactions, and a maxTime of %d ms ...", n, maxTime)

	dp := make([][]float64, n+1)
	for i := range dp {
		dp[i] = make([]float64, maxTime+1)
	}
	//fmt.Println("dp initially: %+v", dp)

	for i := 1; i <= n; i++ {
		for t := 1; t <= maxTime; t++ {
			if transactions[i-1].ProcessingTime <= t {
				dp[i][t] = max(dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount)
				//fmt.Printf("the two choices were %f and %f -- we chose %f\n", dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount, dp[i][t])
			} else {
				//fmt.Printf("t was %d and we had to skip b/c processingtime for this tx is: %d\n", t, transactions[i-1].ProcessingTime)
				dp[i][t] = dp[i-1][t]
			}
		}
		//fmt.Printf("dp after populating for tx #%d: %+v\n\n", i, dp)
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

// greedyKnapsack implements the greedy algorithm for the knapsack problem
func greedyKnapsack(transactions []Transaction, maxTime int) (float64, []Transaction) {
	// Sort transactions by value per ms in descending order
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].ValuePerMs > transactions[j].ValuePerMs
	})

	//fmt.Printf("sorted txs: %+v", transactions)
	var totalAmount float64
	var totalTime int
	selectedTransactions := []Transaction{}

	// Include all the best transactions until we run out of time.
	for _, transaction := range transactions {
		if totalTime+transaction.ProcessingTime <= maxTime {
			selectedTransactions = append(selectedTransactions, transaction)
			totalAmount += transaction.Amount
			totalTime += transaction.ProcessingTime
		}
	}

	return totalAmount, selectedTransactions
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
