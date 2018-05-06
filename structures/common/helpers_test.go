package common

import (
	"strconv"
	"testing"
)

func TestHashTransactionWithEmptyString(t *testing.T) {
	plainTextTransaction := ""
	hashedTransaction := "47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU="

	result := HashTransaction(plainTextTransaction)

	if result != hashedTransaction {
		t.Error("Expected " + hashedTransaction + ", got " + result)
	}
}

func TestHashTransactionWithSampleString(t *testing.T) {
	plainTextTransaction := "The quick brown fox jumps over the lazy dog"
	hashedTransaction := "16j7swfXgJRpypq8sAguT41WUeRtPNt2LQLQvzfJ5ZI="

	result := HashTransaction(plainTextTransaction)

	if result != hashedTransaction {
		t.Error("Expected " + hashedTransaction + ", got " + result)
	}
}

func TestIncludesTransactionIsSuccessful(t *testing.T) {
	transactionsSet := []string{"A", "B", "C", "D"}
	transaction := "C"

	index, err := Includes(transaction, transactionsSet)

	if index != 2 {
		t.Error("Expected index 3, got " + strconv.Itoa(index))
	}

	if err != nil {
		t.Error("Expected error nil, got " + err.Error())
	}
}

func TestIncludesTransactionFails(t *testing.T) {
	transactionsSet := []string{"A", "B", "C", "D"}
	transaction := "Z"

	index, err := Includes(transaction, transactionsSet)

	if index != -1 {
		t.Error("Expected index -1, got " + strconv.Itoa(index))
	}

	if err == nil {
		t.Error("Expected error " + err.Error() + ", got nil")
	}
}
