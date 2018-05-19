package fastmt

import (
	"strconv"
	"testing"
)

func TestBuildFastMerkleTreeWithNoElements(t *testing.T) {
	_, err := NewFastMerkleTree([]string{})

	if err == nil {
		t.Error("Expected error for empty tree")
	}
}

func TestBuildFastMerkleTreeWithSeveralElements(t *testing.T) {
	mt, _ := NewFastMerkleTree([]string{"A", "B", "C", "D"})
	depth := mt.Root.Depth()

	if depth != 3 {
		t.Error("Expected depth 3, got " + strconv.Itoa(depth))
	}
}

func TestBuildUnbalancedFastMerkleTreeWithSeveralElements(t *testing.T) {
	mt, _ := NewFastMerkleTree([]string{"A", "B", "C"})
	depth := mt.Root.Depth()

	if depth != 3 {
		t.Error("Expected depth 3, got " + strconv.Itoa(depth))
	}
}

func TestVerifyFastMerkleTreeValidTransaction(t *testing.T) {
	data := []string{"A", "B", "C"}
	_, _, result, _ := VerifyTransaction("B", data)

	if !result {
		t.Error("Expected true, got false")
	}
}

func TestVerifyFastMerkleTreeInvalidTransaction(t *testing.T) {
	data := []string{"A", "B", "C"}
	_, _, _, err := VerifyTransaction("Z", data)

	if err == nil {
		t.Error("Invalid verification")
	}
}
