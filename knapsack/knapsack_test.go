package knapsack

import (
	"testing"
)

// Test fixture setup
func getTestTransactions() []Transaction {
	return []Transaction{
		{ID: "1", Amount: 60, ProcessingTime: 10},
		{ID: "2", Amount: 100, ProcessingTime: 20},
		{ID: "3", Amount: 120, ProcessingTime: 30},
		{ID: "4", Amount: 60, ProcessingTime: 30},
	}
}

// Helper function to check the result
func checkResult(t *testing.T, result []Transaction) {
	expectedIDs := map[string]struct{}{
		"2": {},
		"3": {},
	}

	resultIDs := make(map[string]struct{})
	for _, tx := range result {
		resultIDs[tx.ID] = struct{}{}
	}

	for id := range expectedIDs {
		if _, found := resultIDs[id]; !found {
			t.Errorf("missing expected transaction ID: %s", id)
		}
	}

	for id := range resultIDs {
		if _, expected := expectedIDs[id]; !expected {
			t.Errorf("unexpected transaction ID found: %s", id)
		}
	}
}

func TestOneDimKeepsies(t *testing.T) {
	transactions := getTestTransactions()
	checkResult(t, OneDimArrayKeepsies{}.Run(transactions, 50))
}

func TestTwoDimNoKeepsies(t *testing.T) {
	transactions := getTestTransactions()
	checkResult(t, TwoDimArrayKeepsies{}.Run(transactions, 50))
}

func TestTwoDimKeepsies(t *testing.T) {
	transactions := getTestTransactions()
	checkResult(t, TwoDimArrayNoKeepsies{}.Run(transactions, 50))
}
