package asl

import (
	"errors"
	"math/rand"

	. "github.com/SimoneStefani/thesis-algorithms/structures/common"
)

type Node struct {
	prev *Node
	next *Node
	hash string
	tr   string
	down *Node
	rank int
}

type List struct {
	level int
	head  *Node
	tail  *Node
}

type SkipList struct {
	levels int
	lists  []List
}

func coinFlipIsHead() bool {
	return rand.Intn(2) == 1
}

func down(node Node) *Node {
	return node.down
}

func hopforward(node Node) *Node {
	return node.next
}

/* 	Commutative Hashing Funtion as suggested in:
*		Goodrich, M. T., & Tamassia, R. (2000).
*		Efficient authenticated dictionaries with skip lists and commutative hashing.
*		Technical Report, Johns Hopkins Information Security Institute.
*  A commutative hash function 'h(x,y)' can be constructed with a cryptographic hash funtion 'f()'
* in the following why: h(x,y) = f(min{x,y}, max{x,y})
 */
func CommutativeHash(x Node, y Node) (string, error) {
	if x.rank == y.rank {
		return "", errors.New("Error: Two nodes with equal rank.")
	} else if x.rank <= y.rank {
		return HashTransaction(x.tr + y.tr), nil
	}
	return HashTransaction(y.tr + x.tr), nil
}
