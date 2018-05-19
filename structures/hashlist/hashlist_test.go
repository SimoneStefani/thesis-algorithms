package hashlist

import (
	"strconv"
	"testing"
)

func TestBuildHashlistWithNoElements(t *testing.T) {
	_, err := NewHashList([]string{})

	if err == nil {
		t.Error("Expected error for empty list")
	}
}

func TestBuildHashlistWithOneElement(t *testing.T) {
	hashlist, _ := NewHashList([]string{"A"})
	length := hashlist.Length()

	if length != 1 {
		t.Error("Expected length 1, got " + strconv.Itoa(length))
	}
}

func TestBuildHashlistWithSeveralElements(t *testing.T) {
	hashlist, _ := NewHashList([]string{"A", "B", "C"})
	length := hashlist.Length()

	if length != 3 {
		t.Error("Expected length 3, got " + strconv.Itoa(length))
	}
}

func TestVerifyHashlistValidTransaction(t *testing.T) {
	data := []string{"A", "B", "C"}
	_, _, result, _ := VerifyTransaction("B", data)

	if !result {
		t.Error("Expected true, got false")
	}
}

func TestVerifyHashlistValidTransactionInFirstElement(t *testing.T) {
	data := []string{"A", "B", "C"}
	_, _, result, _ := VerifyTransaction("A", data)

	if !result {
		t.Error("Expected true, got false")
	}
}

func TestVerifyHashlistInvalidTransaction(t *testing.T) {
	data := []string{"A", "B", "C"}
	_, _, _, err := VerifyTransaction("Z", data)

	if err == nil {
		t.Error("Invalid verification")
	}
}
