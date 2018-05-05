package asl

import (
	"errors"
	"fmt"
	//. "github.com/SimoneStefani/thesis-algorithms/structures/common"
)

type Node struct {
	prev  *Node
	next  *Node
	tr    string
	down  *Node
	index int
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

func NewSkipList(data []string) (*SkipList, error) {
	skiplist, err := buildSkipList(data)

	if err != nil {
		return nil, err
	}

	return skiplist, nil
}

func buildSkipList(data []string) (*SkipList, error) {

	if len(data) == 0 {
		return nil, errors.New("Error: cannot construct skip list with no content.")
	}

	list := &List{
		level: 0,
		head:  nil,
		tail:  nil,
	}
	sl := &SkipList{
		lists:  []List{*list},
		levels: 0,
	}

	for _, tr := range data {
		sl = appendToSkipList(*sl, tr)
	}

	return sl, nil
}

func appendToSkipList(sl SkipList, tr string) *SkipList {

	sl.lists[0] = *insert(sl.lists[0], tr)

	nextLevel := 1
	for {
		if sl.lists[nextLevel-1].tail.index%2 == 0 {
			break
		}
		if nextLevel > sl.levels {
			sl.levels = nextLevel
			newList := &List{
				level: nextLevel,
				head:  nil,
				tail:  nil,
			}
			sl.lists = append(sl.lists, *newList)
		}
		sl.lists[nextLevel] = *insert(sl.lists[nextLevel], tr)
		sl.lists[nextLevel].tail.down = sl.lists[nextLevel-1].tail
		nextLevel = nextLevel + 1
	}

	return &sl
}

func insert(list List, tr string) *List {

	if list.head == nil {
		new := &Node{
			prev: nil,
			next: nil,
			tr:   tr,
			//tr:   HashTransaction(HashTransaction(tr)),
			down:  nil,
			index: 0,
		}
		list.head = new
		list.tail = new
	} else {
		new := &Node{
			next: nil,
			prev: list.tail,
			tr:   tr,
			//tr:   HashTransaction(list.tail.tr + HashTransaction(tr)),
			down: nil,
		}
		list.tail = new
		new.prev.next = new
		new.index = new.prev.index + 1
	}

	return &list
}

func down(node Node) *Node {
	return node.down
}

func hopforward(node Node) *Node {
	return node.next
}

func PrintList(sl SkipList) {
	for _, list := range sl.lists {
		currentNode := list.head
		fmt.Printf("Level %d: ", list.level)
		for {
			if currentNode == nil {
				break
			}
			fmt.Printf("-> %s ", currentNode.tr)
			currentNode = currentNode.next
		}
		fmt.Print("\n")
	}
}

/* 	Commutative Hashing Funtion as suggested in:
*		Goodrich, M. T., & Tamassia, R. (2000).
*		Efficient authenticated dictionaries with skip lists and commutative hashing.
*		Technical Report, Johns Hopkins Information Security Institute.
*  A commutative hash function 'h(x,y)' can be constructed with a cryptographic hash funtion 'f()'
* in the following why: h(x,y) = f(min{x,y}, max{x,y})
 */

// func CommutativeHash(x Node, y Node) (string, error) {
// 	if x.rank == y.rank {
// 		return "", errors.New("Error: Two nodes with equal rank.")
// 	} else if x.rank <= y.rank {
// 		return HashTransaction(x.tr + y.tr), nil
// 	}
// 	return HashTransaction(y.tr + x.tr), nil
// }
