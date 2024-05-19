package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"txoptimizer/knapsack"
)

// LatencyMapping represents the mapping of country codes to processing times
type LatencyMapping map[string]int

func main() {
	// Define command line arguments with default values
	csvFilePath := flag.String("csvFilePath", "challenge_source/transactions.csv", "Path to the input transactions CSV file")
	jsonFilePath := flag.String("jsonFilePath", "challenge_source/latencies.json", "Path to the latency mapping JSON file")
	totalTime := flag.Int("totalTime", 1000, "Total processing time available, in milliseconds")

	// Parse command line arguments
	flag.Parse()

	// Read the latencies from the JSON file
	latencies, err := readLatencies(*jsonFilePath)
	if err != nil {
		log.Fatalf("failed to read latencies: %v", err)
	}

	// Read the transactions from the CSV file
	transactions, err := readTransactions(*csvFilePath)
	if err != nil {
		log.Fatalf("failed to read transactions: %v", err)
	}

	// Update tx with latency predictions
	updateTxWithLatency(transactions, latencies)

	// Get the selected transactions
	selectedTransactions := prioritize(transactions, *totalTime)

	// Calculate the total value
	totalValue := 0.0
	totalMsUsed := 0
	for _, transaction := range selectedTransactions {
		totalValue += transaction.Amount
		totalMsUsed += transaction.ProcessingTime
	}

	// Output the results
	fmt.Printf("Max USD value: %.2f\n", totalValue)
	fmt.Printf("consuming total MS time: %d\n", totalMsUsed)
	fmt.Println("Selected transactions:")
	for _, transaction := range selectedTransactions {
		fmt.Printf("ID: %s, Amount: %.2f, Processing Time: %d\n", transaction.ID, transaction.Amount, transaction.ProcessingTime)
	}
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
func readTransactions(filePath string) ([]knapsack.Transaction, error) {
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

	transactions := make([]knapsack.Transaction, 0, len(records)-1)
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}
		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, knapsack.Transaction{
			ID:              record[0],
			Amount:          amount,
			BankCountryCode: record[2],
		})
	}

	return transactions, nil
}

func updateTxWithLatency(transactions []knapsack.Transaction, latencies LatencyMapping) {
	// Add processing time and value per ms to each transaction
	for i := range transactions {
		if latency, exists := latencies[transactions[i].BankCountryCode]; exists {
			transactions[i].TxNum = i + 1 //for debug easier to look at than uuid
			transactions[i].ProcessingTime = latency
		} else {
			log.Fatalf("no latency found for country code: %s", transactions[i].BankCountryCode)
		}
	}
}

func prioritize(transaction []knapsack.Transaction, totalTime int) []knapsack.Transaction {
	// Create a knapsack prioritizer with the desired implementation, then run it.
	prioritizer := knapsack.KnapsackPrioritizer{Algorithm: knapsack.TwoDimArrayKeepsies{}}
	return prioritizer.Run(transaction, totalTime)
}
