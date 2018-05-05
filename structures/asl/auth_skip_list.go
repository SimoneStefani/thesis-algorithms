package asl

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	//. "github.com/SimoneStefani/thesis-algorithms/structures/common"
)

type Node struct {
	prev  *Node
	next  *Node
	tr    string
	auth  string
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

	var currentIndex int
	var authBuffer string
	if sl.lists[0].tail == nil {
		currentIndex = 0
	} else {
		currentIndex = sl.lists[0].tail.index + 1
	}

	sl.lists[0] = *insert(sl.lists[0], tr, currentIndex)
	authBuffer = authBuffer + computePartialAuthenticator(*sl.lists[0].tail, 0)

	nextLevel := 1
	for {
		if (currentIndex+1)%int(math.Pow(2.0, float64(nextLevel))) != 0 {
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
		sl.lists[nextLevel] = *insert(sl.lists[nextLevel], tr, currentIndex)
		sl.lists[nextLevel].tail.down = sl.lists[nextLevel-1].tail
		authBuffer = authBuffer + computePartialAuthenticator(*sl.lists[nextLevel].tail, nextLevel)
		nextLevel = nextLevel + 1
	}

	sl.lists[0].tail.auth = authBuffer

	return &sl
}

func insert(list List, tr string, index int) *List {

	if list.head == nil {
		new := &Node{
			prev:  nil,
			next:  nil,
			tr:    tr,
			down:  nil,
			index: index,
		}
		list.head = new
		list.tail = new
	} else {
		new := &Node{
			next:  nil,
			prev:  list.tail,
			tr:    tr,
			down:  nil,
			index: index,
		}
		list.tail = new
		new.prev.next = new
	}

	return &list
}

func Includes(sl SkipList, tr string) (int, bool) {

	currentNode := sl.lists[sl.levels-1].head
	if currentNode == nil {
		return -1, false
	}

	nextNode := currentNode.next
	for {
		if tr == currentNode.tr {
			return currentNode.index, true
		}
		if nextNode == nil {
			if currentNode.down != nil {
				currentNode = currentNode.down
				nextNode = currentNode.next
				if nextNode == nil {
					for {
						if currentNode.down.next != nil {
							currentNode = currentNode.down
							nextNode = currentNode.next
							break
						}
						currentNode = currentNode.down
					}
				}
			} else {
				return -1, false
			}
		}
		if tr >= nextNode.tr {
			currentNode = nextNode
			nextNode = currentNode.next
		} else {
			if currentNode.down != nil {
				currentNode = currentNode.down
				nextNode = currentNode.next
			} else {
				nextNode = currentNode
				currentNode = currentNode.prev
			}
		}
	}
}

func isGreaterThan(x string, y string) bool {
	//Compares two strings in alpabetical order
	return true
}

func computePartialAuthenticator(node Node, level int) string {
	prevAuth := ""
	if node.prev != nil {
		prevAuth = "|" + node.prev.auth
	}
	return "{" + strconv.Itoa(node.index) + "|" + strconv.Itoa(level) + "|" + node.tr + prevAuth + "}"
}

func PrintList(sl SkipList) {
	for i := len(sl.lists) - 1; i >= 0; i-- {
		list := sl.lists[i]
		currentNode := list.head
		gaps := int(math.Pow(2.0, float64(list.level))) - 1
		fmt.Printf("Level %d: ", list.level)
		for {
			if currentNode == nil {
				break
			}
			for j := gaps; j > 0; j-- {
				fmt.Printf("-----")
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
