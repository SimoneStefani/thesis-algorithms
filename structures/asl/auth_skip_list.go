package asl

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	. "github.com/SimoneStefani/thesis-algorithms/structures/common"
)

type Node struct {
	prev  *Node
	next  *Node
	tr    string
	auth  string
	down  *Node
	up    *Node
	index int
}

type List struct {
	level  int
	length int
	head   *Node
	tail   *Node
}

type SkipList struct {
	levels int
	auth   string
	lists  []List
}

type ProofComponent struct {
	tr            string
	authenticator string
}

// Example SL with transactions "a"-"j"
// Level 3: ------------------------------------> h
// Level 2: ----------------> d ----------------> h
// Level 1: ------> b ------> d ------> f ------> h ------> j
// Level 0: -> a -> b -> c -> d -> e -> f -> g -> h -> i -> j
func NewSkipList(data []string) (*SkipList, error) {
	skiplist, err := buildSkipList(data)

	if err != nil {
		return nil, err
	}

	return skiplist, nil
}

// A membership claim has the form “Data element 'tr' occupies the i-th position
// of the AASL whose n-th authenticator is known to the verifier,” and is denoted by ⟨i,n,d⟩.
// i = pos(tr)
// n = SkipList.auth
// d = tr
func VerifyTransaction(sl SkipList, tr string) (string, []string, bool, error) {

	_, nodePointer, exists := Lookup(sl, tr)
	if !exists {
		return "", nil, false, errors.New("error: not part of skip list")
	}
	computeMembershipProof(*nodePointer, tr, sl)

	// Incomplete
	return "", nil, false, nil
}

// computeMembershipProof compoutes the membership for a node given a skip list
func computeMembershipProof(node Node, tr string, sl SkipList) ([]ProofComponent, error) {
	var membershipProof []ProofComponent
	var tempProofComponent ProofComponent
	var tempNode = node
	index := node.index
	var singleHopTraversalLevel int

	for {
		if index > sl.lists[0].length {
			return membershipProof, nil
		}
		tempProofComponent = computeProofComponent(tempNode)
		membershipProof = append(membershipProof, tempProofComponent)
		singleHopTraversalLevel = SingleHopTraversalLevel(index, sl.lists[0].length-1)
		tempNode = singleHopTraversal(tempNode, singleHopTraversalLevel)
		index = index + int(math.Pow(2.0, float64(singleHopTraversalLevel)))
	}
}

// computeProofComponent takes a node and returns a proof component C for the
// AASL element in position j.
func computeProofComponent(node Node) ProofComponent {

	proofComponent := &ProofComponent{
		tr:            node.tr,
		authenticator: "",
	}
	tempNode := node
	for {
		if tempNode.up == nil {
			return *proofComponent
		}
		proofComponent.authenticator = proofComponent.authenticator + dropToBaseElement(*tempNode.prev).auth
		tempNode = *tempNode.up
	}
}

// dropToBaseElement taks a node at a certain List level and returns its corresponding
// equal in the base list
func dropToBaseElement(node Node) Node {
	tempNode := node
	for {
		if tempNode.down == nil {
			return tempNode
		}
		tempNode = *tempNode.down
	}
}

func singleHopTraversal(startNode Node, level int) Node {
	tempNode := startNode
	for i := 0; i <= level; i++ {
		tempNode = *tempNode.up
	}
	return *tempNode.next
}

// SingleHopTraversalLevel returns the highest linked list level l that must be followed in the Skip List
// in order to travel from element at 'start' to element at 'end'
func SingleHopTraversalLevel(start int, end int) int {
	currentLevel := 0
	highestLevel := 0
	var temp int

	for {
		temp = int(math.Pow(2.0, float64(currentLevel)))
		if start%temp != 0 {
			break
		}
		if start+temp <= end {
			highestLevel = currentLevel
		} else {
			return highestLevel
		}
		currentLevel = currentLevel + 1
	}
	return highestLevel
}

func buildSkipList(data []string) (*SkipList, error) {

	if len(data) == 0 {
		return nil, errors.New("Error: cannot construct skip list with no content.")
	}

	firstNode := &Node{
		prev:  nil,
		next:  nil,
		tr:    data[0],
		auth:  HashTransaction(HashTransaction(data[0])),
		down:  nil,
		up:    nil,
		index: 0,
	}

	list := &List{
		level:  0,
		head:   firstNode,
		tail:   firstNode,
		length: 0,
	}
	list.head.next = list.tail
	list.tail.prev = list.head

	sl := &SkipList{
		lists:  []List{*list},
		levels: 0,
	}

	for i := 1; i < len(data); i++ {
		sl = appendToSkipList(*sl, data[i])
	}
	// TO DO: check, is always top list?
	sl.auth = sl.lists[len(sl.lists)-1].tail.auth

	return sl, nil
}

func appendToSkipList(sl SkipList, tr string) *SkipList {

	var currentIndex int
	var authBuffer string
	currentIndex = sl.lists[0].tail.index + 1

	// Insert to base list
	sl.lists[0] = *insert(sl.lists[0], tr, currentIndex)
	// Authentication buffer is used to computer the Authenticator of the base Node
	authBuffer = computePartialAuthenticator(sl.lists[0], *sl.lists[0].tail, 0)

	// Insert to all upper lists that the element belongs to
	nextLevel := 1
	for {
		if (currentIndex)%int(math.Pow(2.0, float64(nextLevel))) != 0 {
			break
		}
		if nextLevel > sl.levels {
			sl.levels = nextLevel
			newList := &List{
				level:  nextLevel,
				head:   nil,
				tail:   nil,
				length: 0,
			}
			sl.lists = append(sl.lists, *newList)

			// Since new List, Insert the Head of the Baselist to the new list
			sl.lists[nextLevel] = *insert(sl.lists[nextLevel], sl.lists[0].head.tr, 0)
			sl.lists[nextLevel].head.down = sl.lists[nextLevel-1].head
			sl.lists[nextLevel-1].head.up = sl.lists[nextLevel].head

		}
		sl.lists[nextLevel] = *insert(sl.lists[nextLevel], tr, currentIndex)
		sl.lists[nextLevel].tail.down = sl.lists[nextLevel-1].tail
		sl.lists[nextLevel-1].tail.up = sl.lists[nextLevel].tail
		authBuffer = authBuffer + sl.lists[nextLevel-1].tail.auth
		nextLevel = nextLevel + 1
	}

	sl.lists[0].tail.auth = HashTransaction(authBuffer)

	return &sl
}

func insert(list List, tr string, index int) *List {

	if list.head == nil {
		new := &Node{
			prev:  nil,
			next:  nil,
			tr:    tr,
			down:  nil,
			up:    nil,
			index: index,
			auth:  HashTransaction(HashTransaction(tr)),
		}
		list.head = new
		list.tail = new
	} else {
		new := &Node{
			next:  nil,
			prev:  list.tail,
			tr:    tr,
			down:  nil,
			up:    nil,
			index: index,
		}
		list.tail = new
		new.prev.next = new
		new.auth = computePartialAuthenticator(list, *list.tail, list.level)
		list.length = list.length + 1
	}

	return &list
}

func Lookup(sl SkipList, tr string) (int, *Node, bool) {

	// fmt.Print("\n")
	// fmt.Printf("Searching Transaction ---> %s\n", tr)

	currentLevel := sl.levels - 1
	nextNode := sl.lists[currentLevel].head
	currentNode := nextNode
	if nextNode == nil {
		return -1, nil, false
	}

	// Find list to start from
	for {
		if tr >= nextNode.tr {
			currentNode = nextNode
			nextNode = currentNode.next
			// fmt.Printf("Starts on Level %d with CurrentNode: %s\n\n", currentLevel, currentNode.tr)
			break
		} else {
			currentLevel = currentLevel - 1
			if currentLevel < 0 {
				return -1, nil, false
			}
			nextNode = sl.lists[currentLevel].head
		}
	}

	for {
		//Check existence of current nodes
		for {
			if nextNode == nil {
				// fmt.Print("LOOP -> ")
				// fmt.Printf("NextNode is Null, CurrentNode is: %s\n", currentNode.tr)
				if currentNode.down == nil {
					if currentNode.tr == tr {
						return currentNode.index, currentNode, true
					}
					return -1, nil, false
				}
				currentNode = currentNode.down
				nextNode = currentNode.next
			} else {
				//fmt.Printf("CurrentNode is: %s\n", currentNode.tr)
				break
			}
		}

		if tr == currentNode.tr {
			return currentNode.index, currentNode, true
		}
		if tr >= nextNode.tr {
			currentNode = nextNode
			nextNode = currentNode.next
		} else {
			if currentNode.down != nil {
				currentNode = currentNode.down
				nextNode = currentNode.next
			} else {
				currentNode = nextNode
				nextNode = currentNode.next
			}
		}
	}
}

func computePartialAuthenticator(list List, node Node, level int) string {
	prevAuth := ""

	if node.prev != nil {
		tempNode := node.prev
		for {
			if tempNode.down == nil {
				prevAuth = tempNode.auth
				break
			}
			tempNode = tempNode.down
		}
	}
	return HashTransaction(strconv.Itoa(node.index) + strconv.Itoa(level) + node.tr + prevAuth)
}

//For debugging purposes only
func dummyComputePartialAuthenticator(node Node, level int) string {
	prevAuth := ""

	if node.prev != nil {
		tempNode := node.prev
		for {
			if tempNode.down == nil {
				prevAuth = "|" + tempNode.auth
				break
			}
			tempNode = tempNode.down
		}
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
			if currentNode != list.head {
				for j := gaps; j > 0; j-- {
					fmt.Printf("-----")
				}
			}
			fmt.Printf("-> %s ", currentNode.tr)
			currentNode = currentNode.next
		}
		fmt.Print("\n")
	}
}

func PrintListAuthenticators(sl SkipList) {
	for i := len(sl.lists) - 1; i >= 0; i-- {
		list := sl.lists[i]
		currentNode := list.head
		gaps := int(math.Pow(2.0, float64(list.level))) - 1
		fmt.Printf("Level %d: ", list.level)
		for {
			if currentNode == nil {
				break
			}
			if currentNode != list.head {
				for j := gaps; j > 0; j-- {
					fmt.Printf("---------")
				}
			}
			fmt.Printf("-> %s ", currentNode.auth[0:5])
			currentNode = currentNode.next
		}
		fmt.Print("\n")
	}
}

func PrintSkipListHeadsAndTails(sl SkipList) {
	for j := len(sl.lists) - 1; j >= 0; j-- {
		list := sl.lists[j]
		fmt.Printf("List %d --> HEAD: %s , TAIL: %s\n", j, list.head.tr, list.tail.tr)
	}
	fmt.Print("\n")
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
