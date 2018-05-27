package mt

import (
	"strconv"
	"testing"
)

func TestBuildMerkleTreeWithNoElements(t *testing.T) {
	_, err := NewTree([]string{})

	if err == nil {
		t.Error("Expected error for empty tree")
	}
}

func TestBuildMerkleTreeWithSeveralElements(t *testing.T) {
	mt, _ := NewTree([]string{"A", "B", "C", "D"})
	depth := mt.Root.Depth()

	if depth != 3 {
		t.Error("Expected depth 3, got " + strconv.Itoa(depth))
	}
}

func TestBuildUnbalancedMerkleTreeWithSeveralElements(t *testing.T) {
	mt, _ := NewTree([]string{"A", "B", "C"})
	depth := mt.Root.Depth()

	if depth != 3 {
		t.Error("Expected depth 3, got " + strconv.Itoa(depth))
	}
}

func TestVerifyMerkleTreeValidTransaction(t *testing.T) {
	data := []string{"A", "B", "C"}
	tree, _ := NewTree(data)
	_, _, result, _ := VerifyTransaction("B", data, tree)

	if !result {
		t.Error("Expected true, got false")
	}
}

func TestVerifyMerkleTreeInvalidTransaction(t *testing.T) {
	data := []string{"A", "B", "C"}
	tree, _ := NewTree(data)
	_, _, _, err := VerifyTransaction("Z", data, tree)

	if err == nil {
		t.Error("Invalid verification")
	}
}
