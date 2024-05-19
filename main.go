package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	//"github.com/mattschofield/go-knapsack"
	"log"
	"os"
	"strconv"
	"txoptimizer/knapsack"
)

//// Implement the Packable interface for Transaction (for use with the library I'm trying out)
//func (t Transaction) Weight() int64 {
//	return int64(t.ProcessingTime)
//}
//func (t Transaction) Value() int64 {
//	return int64(t.Amount)
//}

// TODO
//*  make the function prioritize(transaction []Transaction, totalTime=1000): []Transaction {} function be what is called from main
//* it will load the resources and call the one dimensional version
// * main will report the max USD Value of this as well.
// * main should have cli arg to set total time. run it with all the options see the output looks right.
// * make a test file
// * make a interface for this just the useLibraryKnapsack(transactions []Transaction, maxTime int) ([]Transaction)
// * do the max value calc outside of the implementation, loop over the tx to do that in maybe a knapsackprioritizer that can take the tx and latency things and can be given a particular interface implementation of knapsackalgo and then .Run it
// * make a small in memory fixture of test data with expected results for a unit test
// * use the fixture on 3 different implementations in the test file, each one just a same thing with knapsackprioritizer getting a different knapsackalgo passed to it
// make sure everything is clean and then
// make variable names consistent eg totalTime
// * determine if i really does need to return the updated tx after adding latency or can i edit them in place?
// invite the ppl to the repo

// LatencyMapping represents the mapping of country codes to processing times
type LatencyMapping map[string]int

func main() {
	// Define command line arguments with default values
	csvFilePath := flag.String("csvFilePath", "challenge_source/transactions.csv", "Path to the input transactions CSV file")
	jsonFilePath := flag.String("jsonFilePath", "challenge_source/latencies.json", "Path to the latency mapping JSON file")
	totalTime := flag.Int("totalTime", 1000, "Total processing time available, in milliseconds")

	// Parse command line arguments
	flag.Parse()

	//get a number of ms, tx and latency files from cli args, with some defaults

	// Paths to the input files
	//csvFilePath := "challenge_source/transactions.csv"
	////csvFilePath := "challenge_source/transactions_mini.csv"
	//jsonFilePath := "challenge_source/latencies.json"
	//totalTime := 1000

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
	//////////////////////////
	//
	//// Array of time options in milliseconds
	//timeOptions := []int{50, 60, 90, 1000}
	////timeOptions := []int{50}
	//for _, maxTime := range timeOptions {
	//	// Run the knapsack function for each time option
	//	runKnapsack(transactions, maxTime)
	//}
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
		//if i >= 10 {
		//	break
		//} //hax just work with 10 for now
	}

	return transactions, nil
}

//func updateTxWithLatency(transactions []knapsack.Transaction, latencies LatencyMapping) []knapsack.Transaction {
//	// Add processing time and value per ms to each transaction
//	for i := range transactions {
//		if latency, exists := latencies[transactions[i].BankCountryCode]; exists {
//			transactions[i].TxNum = i + 1 //for debug easier to look at than uuid
//			transactions[i].ProcessingTime = latency
//		} else {
//			log.Fatalf("no latency found for country code: %s", transactions[i].BankCountryCode)
//		}
//	}
//	return transactions
//}

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

	// Create a knapsack prioritizer with the TwoDimArrayNoKeepsies implementation
	//prioritizer := knapsack.KnapsackPrioritizer{Algorithm: knapsack.TwoDimArrayNoKeepsies{}}
	//prioritizer := knapsack.KnapsackPrioritizer{Algorithm: knapsack.TwoDimArrayKeepsies{}}
	prioritizer := knapsack.KnapsackPrioritizer{Algorithm: knapsack.OneDimArrayKeepsies{}}
	return prioritizer.Run(transaction, totalTime)
}

//
//func runKnapsack(transactions []knapsack.Transaction, maxTime int) {
//	maxAmount, selectedTransactions := myknapsack(transactions, maxTime)
//
//	fmt.Printf("For max time %d ms:\n", maxTime)
//	fmt.Printf("Max amount: %.2f\n", maxAmount)
//	fmt.Println("Selected transactions:")
//	for _, transaction := range selectedTransactions {
//		fmt.Println(transaction)
//	}
//
//	//maxAmountGreedy, selectedTransactionsGreedy := greedyKnapsack(transactions, maxTime)
//	//
//	//fmt.Printf("Greedy Max amount: %.2f\n", maxAmountGreedy)
//	//fmt.Println("Greedy Selected transactions:")
//	//for _, transaction := range selectedTransactionsGreedy {
//	//	fmt.Println(transaction)
//	//}
//
//	maxAmountOneD, selectedTransactionsOneD := onedimensionalknapsack(transactions, maxTime)
//
//	fmt.Printf("OneD Max amount: %.2f\n", maxAmountOneD)
//	fmt.Printf("OneD Selected transactions (count: %d):\n", len(selectedTransactionsOneD))
//	for _, transaction := range selectedTransactionsOneD {
//		fmt.Println(transaction)
//	}
//
//	//maxAmtLib, selectedTxLib := useLibraryKnapsack(transactions, maxTime)
//	//
//	//fmt.Printf("Library Knapsack Max amount: %.2f\n", maxAmtLib)
//	//fmt.Printf("Library Knapsack transactions (count: %d):\n", len(selectedTxLib))
//	//for _, transaction := range selectedTxLib {
//	//	fmt.Println(transaction)
//	//}
//
//	fmt.Println("----------------------")
//}

//
//// Knapsack problem solution
//func myknapsack(transactions []knapsack.Transaction, maxTime int) (float64, []knapsack.Transaction) {
//
//	n := len(transactions)
//	fmt.Printf("looking at %d transactions, and a maxTime of %d ms ...\n", n, maxTime)
//
//	dp := make([][]float64, n+1)
//	for i := range dp {
//		dp[i] = make([]float64, maxTime+1)
//	}
//	//fmt.Println("dp initially: %+v", dp)
//
//	// initial implementation suggestion did not do a keep array, it did it
//	// Initialize the keep array to track which items are included
//	keep := make([][]bool, n+1)
//	for i := range keep {
//		keep[i] = make([]bool, maxTime+1)
//	}
//
//	for i := 1; i <= n; i++ {
//		for t := 1; t <= maxTime; t++ {
//			if transactions[i-1].ProcessingTime <= t {
//				dp[i][t] = max(dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount)
//
//				///not normally doing this lol
//				if dp[i-1][t] < dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount {
//					//fmt.Printf("KEEP i: %d and t: %d: choose %.2f vs %.2f (= %.2f + %.2f) -- we chose %.2f\n", i, t, dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount, dp[i-1][t-transactions[i-1].ProcessingTime], transactions[i-1].Amount, dp[i][t])
//					keep[i][t] = true
//				}
//
//				//fmt.Printf("i: %d and t: %d: choose %.2f vs %.2f (= %.2f + %.2f) -- we chose %.2f\n", i, t, dp[i-1][t], dp[i-1][t-transactions[i-1].ProcessingTime]+transactions[i-1].Amount, dp[i-1][t-transactions[i-1].ProcessingTime], transactions[i-1].Amount, dp[i][t])
//				//fmt.Printf("fookin wot: %+v\n", dp[i-1])
//				//fmt.Printf("fookin wot2: %+v\n", dp[i-1][t-transactions[i-1].ProcessingTime])
//				//panic("just stop")
//			} else {
//				//fmt.Printf("i: %d and t: %d and we had to skip b/c processingtime for this tx is: %d\n", i, t, transactions[i-1].ProcessingTime)
//				dp[i][t] = dp[i-1][t]
//			}
//		}
//		//fmt.Printf("dp after populating for tx #%d: %+v\n\n", i, dp)
//	}
//	//for i, v := range dp {
//	//	fmt.Printf("dp row %d after populating: %+v\n", i, v)
//	//}
//	//for i, tx := range transactions {
//	//	fmt.Printf("txid: %s -> keep timeslots: %+v\n", tx.ID, keep[i+1])
//	//}
//
//	// Find the selected transactions
//	res := []knapsack.Transaction{}
//	t := maxTime
//	for i := n; i > 0; i-- {
//		//fmt.Printf("i: %d, t: %d\n", i, t)
//		if dp[i][t] != dp[i-1][t] {
//			res = append(res, transactions[i-1])
//			t -= transactions[i-1].ProcessingTime
//			//	fmt.Printf("appended selected tx, t is now: %d\n", t)
//			//} else {
//			//	fmt.Printf("did nothing because %.2f == %.2f\n", dp[i][t], dp[i-1][t])
//		}
//	}
//	_ = res
//
//	t = maxTime
//	altRes := []knapsack.Transaction{}
//	for i := n; i >= 0; i-- {
//		//fmt.Printf("i: %d t: %d, keep[i][t]: %t\n", i, t, keep[i][t])
//		if keep[i][t] {
//			altRes = append(altRes, transactions[i-1])
//			t -= transactions[i-1].ProcessingTime
//		}
//	}
//	//panic("stop here")
//	return dp[n][maxTime], altRes
//}
//
//func onedimensionalknapsack(transactions []knapsack.Transaction, maxTime int) (float64, []knapsack.Transaction) {
//	//n := len(transactions)
//	//dp := make([]float64, maxTime+1)
//	//
//	//// Fill the dp array
//	//for i := 0; i < n; i++ {
//	//	fmt.Printf("tx: %+v\n", transactions[i])
//	//	for t := maxTime; t >= transactions[i].ProcessingTime; t-- {
//	//		fmt.Printf("t: @%dms, dp[t]: $%.2f, t-transactions[i].ProcessingTime: @%dms, dp[t-transactions[i].ProcessingTime]: $%.2f, transactions[i].Amount: $%.2f (dp[t-transactions[i].ProcessingTime] + transactions[i].Amount: $%.2f)\n", t, dp[t], t-transactions[i].ProcessingTime, dp[t-transactions[i].ProcessingTime], transactions[i].Amount, dp[t-transactions[i].ProcessingTime]+transactions[i].Amount)
//	//		dp[t] = max(dp[t], dp[t-transactions[i].ProcessingTime]+transactions[i].Amount)
//	//	}
//	//	fmt.Printf("after tx %d, dp looks like:%+v\n", i, dp)
//	//}
//	//
//	//// Traceback to find the selected transactions
//	//totalAmount := dp[maxTime]
//	//selectedTransactions := []Transaction{}
//	//t := maxTime
//	//for i := n - 1; i >= 0 && t > 0; i-- {
//	//	if t >= transactions[i].ProcessingTime && dp[t] == dp[t-transactions[i].ProcessingTime]+transactions[i].Amount {
//	//		selectedTransactions = append(selectedTransactions, transactions[i])
//	//		t -= transactions[i].ProcessingTime
//	//	}
//	//}
//	//
//	//return totalAmount, selectedTransactions
//
//	n := len(transactions)
//	dp := make([]float64, maxTime+1)
//	keep := make([][]bool, n)
//
//	// Initialize the keep array to track which items are included
//	for i := range keep {
//		keep[i] = make([]bool, maxTime+1)
//	}
//
//	// Fill the dp array and keep track of decisions
//	for i := 0; i < n; i++ {
//		for t := maxTime; t >= transactions[i].ProcessingTime; t-- {
//			if dp[t] < dp[t-transactions[i].ProcessingTime]+transactions[i].Amount {
//				dp[t] = dp[t-transactions[i].ProcessingTime] + transactions[i].Amount
//				//fmt.Printf("keep: i: %d t: %d tx: %+v\n", i, t, transactions[i])
//
//				keep[i][t] = true
//			}
//		}
//	}
//
//	//for tx, v := range keep {
//	//	fmt.Printf("txid: %s -> keep timeslots: %+v\n", transactions[tx].ID, v)
//	//}
//
//	totalAmount := dp[maxTime]
//	selectedTransactions := []knapsack.Transaction{}
//
//	//Determine which ones we actually sent by walking the timeslots back. In reverse order of tx processing
//	t := maxTime
//	for i := n - 1; i >= 0; i-- {
//		//fmt.Printf("i: %d t: %d, keep[i][t]: %t\n", i, t, keep[i][t])
//		if keep[i][t] {
//			selectedTransactions = append(selectedTransactions, transactions[i])
//			t -= transactions[i].ProcessingTime
//		}
//	}
//
//	return totalAmount, selectedTransactions
//}

// Function to use the library Knapsack function
//func useLibraryKnapsack(transactions []knapsack.Transaction, maxTime int) (float64, []knapsack.Transaction) {
//	// Convert transactions to Packable items
//	items := make([]knapsack.Packable, len(transactions))
//	for i, tx := range transactions {
//		items[i] = tx
//	}
//
//	// Call the Knapsack function from the library
//	//selectedIndices := knapsack.Knapsack(items, int64(maxTime))
//	selectedIndices := xxKnapsack(items, int64(maxTime))
//
//	// Create the result transactions array and calculate the total value
//	selectedTransactions := []Transaction{}
//	totalValue := 0.0
//	for _, index := range selectedIndices {
//		tx := transactions[index]
//		selectedTransactions = append(selectedTransactions, tx)
//		totalValue += tx.Amount
//	}
//
//	return totalValue, selectedTransactions
//}

//// greedyKnapsack implements a greedy algorithm for the knapsack problem -- which doesnt give the best results so, lets not bother with it!
//func greedyKnapsack(transactions []Transaction, maxTime int) (float64, []Transaction) {
//	// Sort transactions by value per ms in descending order
//	sort.Slice(transactions, func(i, j int) bool {
//		return transactions[i].ValuePerMs > transactions[j].ValuePerMs
//	})
//
//	//fmt.Printf("sorted txs: %+v", transactions)
//	var totalAmount float64
//	var totalTime int
//	selectedTransactions := []Transaction{}
//
//	// Include all the best transactions until we run out of time.
//	for _, transaction := range transactions {
//		if totalTime+transaction.ProcessingTime <= maxTime {
//			selectedTransactions = append(selectedTransactions, transaction)
//			totalAmount += transaction.Amount
//			totalTime += transaction.ProcessingTime
//		}
//	}
//
//	return totalAmount, selectedTransactions
//}

//func max(a, b float64) float64 {
//	if a > b {
//		return a
//	}
//	return b
//}

//
//// Knapsack uses a dynamic programming pattern to calculate the maximum value
//// to be gained from an array of items whilst keeping the total weight of items
//// less than or equal to a capacity. It will return the indices of the items
//// to pack.
//// For a very good guide to the 0/1 Knapsack Problem, see: https://www.youtube.com/watch?v=EH6h7WA7sDw
//func xxKnapsack(items []knapsack.Packable, capacity int64) []int64 {
//
//	// We store our working solutions in matrices of N+1 x M+1, where N is the number
//	// of items and M is the capacity. We add 1 so we can index from 0.
//	// `values` stores the sum of a set of items' values.
//	values := make([][]int64, len(items)+1)
//	for i := range values {
//		values[i] = make([]int64, capacity+1)
//	}
//
//	// `keep` stores a matrix of bits, 1 meaning we want to keep the item in this
//	// combination, 0 means we'll leave it.
//	keep := make([][]int, len(items)+1)
//	for i := range keep {
//		keep[i] = make([]int, capacity+1)
//	}
//
//	// Initially, we'll set all combinations in both `values` and `keep` to 0.
//	for i := int64(0); i < capacity+1; i++ {
//		values[0][i] = 0
//		keep[0][i] = 0
//	}
//
//	//for i := 0; i < len(items)+1; i++ {
//	//	values[i][0] = 0
//	//	keep[i][0] = 0
//	//}
//
//	// Simply put, for every item in `items` we want to know whether it will
//	// fit in our sack for every capacity from 0 to `capacity`.
//	// We know that with 0 items or 0 capacity, no outcome is possible, so start
//	// from item 1 and capacity of 1.
//	for i := 1; i <= len(items); i++ {
//		for c := int64(1); c <= capacity; c++ {
//
//			// Does the item fit at this capacity?
//			itemFits := (items[i-1].Weight() <= c)
//			if !itemFits {
//				values[i][c] = values[i-1][c] //still need to set a value if we skip, use previous
//				continue                      // skip this iteration
//			}
//
//			// Is the value of the item, plus the (previously calculated) value of
//			// any remaining space after the addition of this item, greater than the
//			// value gained from the previous item?
//			maxValueAtThisCapacity := items[i-1].Value() + values[i-1][c-items[i-1].Weight()]
//			//fmt.Println(maxValueAtThisCapacity)
//			previousValueAtThisCapacity := values[i-1][c]
//
//			// If the max value to be gained by using this item at this level of
//			// capacity is greater than the value to be gained from using the previous
//			// item at this capacity, then we want to use this item and keep it.
//			// Otherwise, we'll just use the previous item's combination.
//			if maxValueAtThisCapacity > previousValueAtThisCapacity {
//				//fmt.Printf("keep for i %d and c %d\n", i, c)
//				values[i][c] = maxValueAtThisCapacity
//				//fmt.Printf("KEEP i: %d and t: %d: choose %.2f vs %.2f (= %.2f + %.2f) -- we chose %.2f\n", i, c, previousValueAtThisCapacity, maxValueAtThisCapacity, items[i-1].Value(), values[i-1][c-items[i-1].Weight()], maxValueAtThisCapacity)
//				//fmt.Printf("KEEP i: %d and t: %d: choose %d vs %d (= %d + %d) -- we chose %d\n", i, c, previousValueAtThisCapacity, maxValueAtThisCapacity, items[i-1].Value(), values[i-1][c-items[i-1].Weight()], maxValueAtThisCapacity)
//				keep[i][c] = 1
//			} else {
//				//fmt.Printf("previous for i %d and c %d\n", i, c)
//				values[i][c] = previousValueAtThisCapacity
//				keep[i][c] = 0
//			}
//		}
//	}
//
//	//for i, item := range items {
//	//	fmt.Printf("tx ms %d value: %d -> keep timeslots: %+v\n", item.Weight(), item.Value(), keep[i+1])
//	//}
//
//	// We've now calculated the maximum value to be gained from a combination of
//	// items. The maximum value will live at `values[len(items)][capacity]`
//	// We now want to loop through our `keep` array and return the indices that
//	// point to the specific items to pack into our Knapsack.
//	n := len(items)
//	c := capacity
//	var indices []int64
//
//	for n > 0 {
//		if keep[n][c] == 1 {
//			indices = append(indices, int64(n-1))
//			c -= items[n-1].Weight()
//		}
//		n--
//	}
//
//	return indices
//}
