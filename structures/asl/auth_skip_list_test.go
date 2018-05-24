package asl

import (
	"strconv"
	"testing"
)

func TestBuildSkiplistWithNoElements(t *testing.T) {
	_, err := NewSkipList([]string{})

	if err == nil {
		t.Error("Expected error for empty tree")
	}
}

func TestBuildSkiplistWithSeveralElements(t *testing.T) {
	asl, _ := NewSkipList([]string{"A", "B", "C", "D", "E", "F"})
	lengths := asl.Lengths()

	if len(lengths) != 3 {
		t.Error("Expected 3 skiplists, got " + strconv.Itoa(len(lengths)))
	}

	if lengths[0] != 6 {
		t.Error("Expected 6 elements, got " + strconv.Itoa(lengths[0]))
	}

	if lengths[1] != 3 {
		t.Error("Expected 3 elements, got " + strconv.Itoa(lengths[1]))
	}

	if lengths[2] != 2 {
		t.Error("Expected 2 element, got " + strconv.Itoa(lengths[2]))
	}
}

func TestVerifySkiplistValidTransaction(t *testing.T) {
	data := []string{"A", "B", "C"}
	sl, _ := NewSkipList(data)
	result, _, _, _ := VerifyTransaction(*sl, "B")

	if !result {
		t.Error("Expected true, got false")
	}
}

func TestVerifySkiplistInvalidTransaction(t *testing.T) {
	data := []string{"A", "B", "C"}
	sl, _ := NewSkipList(data)
	_, _, _, err := VerifyTransaction(*sl, "Z")

	if err == nil {
		t.Error("Invalid verification")
	}
}
