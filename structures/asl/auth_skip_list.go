package asl

import (
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

/* 	Commutative Hashing Funtion as suggested in:
*		Goodrich, M. T., & Tamassia, R. (2000).
*		Efficient authenticated dictionaries with skip lists and commutative hashing.
*		Technical Report, Johns Hopkins Information Security Institute.
*  A commutative hash function 'h(x,y)' can be constructed with a cryptographic hash funtion 'f()'
* in the following why: h(x,y) = f(min{x,y}, max{x,y})
 */
func CommutativeHash(x Node, y Node) string {
	if x.rank <= y.rank {
		return HashTransaction(x.tr + y.tr)
	}
	return HashTransaction(y.tr + x.tr)
}
