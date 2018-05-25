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
	down  *Node
	up    *Node
	tr    string
	auth  string
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
	authenticator []string
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
func VerifyTransaction(sl SkipList, tr string) (bool, []ProofComponent, *Node, error) {

	_, nodePointer, exists := Lookup(sl, tr)
	if !exists {
		return false, nil, nil, errors.New("error: not part of skip list")
	}
	proof, err := computeMembershipProof(*nodePointer, tr, sl)

	if err != nil {
		return false, nil, nil, err
	}

	verifactionResult := VerifyMembershipProof(*nodePointer, sl, proof)

	// Incomplete
	return verifactionResult, proof, nodePointer, err
}

// ProcessMembershipProof (i,n,d,T,E) return true or false.
// Processes the membership proof E of the membership claim ⟨i, n, d⟩ against authenticator T .
// In this function 'node' holdes the 'index', and 'datum'
// 'sl' holds the authenticator T which alwas is the digest of the Skip List
// 'proof' holds the E
func VerifyMembershipProof(node Node, sl SkipList, proof []ProofComponent) bool {

	currentAuth := processProofComponent(node.index, proof[0])
	prevAuth := currentAuth
	level := SingleHopTraversalLevel(node.index, sl.lists[0].length)
	index := node.index + int(math.Pow(2.0, float64(level)))
	componentCounter := 2

	for i := 1; i < len(proof); i++ {
		if index >= sl.lists[0].length {
			break
		}
		currentAuth = processProofComponent(index, proof[i])
		if proof[i].authenticator[level] != prevAuth {
			return false
		}
		prevAuth = currentAuth
		level = SingleHopTraversalLevel(index, sl.lists[0].length)
		index = index + int(math.Pow(2.0, float64(level)))
		componentCounter++
	}

	return true
}

func (sls *SkipList) Lengths() []int {
	lengths := []int{}

	for _, ls := range sls.lists {
		lengths = append(lengths, ls.length+1)
	}

	return lengths
}

// Processes a single Proof Component --> Calculates Ti
func processProofComponent(index int, component ProofComponent) string {
	buffer := ""
	// The first element is a special case, see paper for more info
	if index == 0 {
		return component.authenticator[0]
	}
	// fmt.Printf("For Node %s --> ", component.tr)
	// fmt.Printf("Datum: %s | ", component.tr)
	// for _, el := range component.authenticator {
	// 	fmt.Printf("%s | ", el[0:5])
	// }
	// fmt.Println()

	for level, auth := range component.authenticator {
		buffer = buffer + HashTransaction(strconv.Itoa(index)+strconv.Itoa(level)+component.tr+auth)
	}
	return HashTransaction(buffer)
}

// computeMembershipProof compoutes the membership for a node given a skip list
func computeMembershipProof(node Node, tr string, sl SkipList) ([]ProofComponent, error) {
	var membershipProof []ProofComponent
	var tempProofComponent ProofComponent
	var tempNode = node
	index := node.index
	var singleHopTraversalLevel int

	for {
		if index >= sl.lists[0].length {
			tempProofComponent = computeProofComponent(tempNode)
			membershipProof = append(membershipProof, tempProofComponent)
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

	if node.prev == nil {
		proofComponent := &ProofComponent{
			tr:            node.tr,
			authenticator: []string{node.auth},
		}
		// fmt.Printf("For Node %s --> ", node.tr)
		// fmt.Printf("Datum: %s | ", proofComponent.tr)
		// fmt.Printf("%s | ", proofComponent.authenticator[0][0:5])
		// fmt.Println()
		return *proofComponent
	}

	proofComponent := &ProofComponent{
		tr: node.tr,
		//authenticator: []string{node.prev.auth},
	}
	tempNode := node
	for {
		proofComponent.authenticator = append(proofComponent.authenticator, dropToBaseElement(*tempNode.prev).auth)
		if tempNode.up == nil {
			// fmt.Printf("For Node %s --> ", node.tr)
			// fmt.Printf("Datum: %s | ", proofComponent.tr)
			// for _, el := range proofComponent.authenticator {
			// 	fmt.Printf("%s | ", el[0:5])
			// }
			// fmt.Println()
			return *proofComponent
		}
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
	tempNode := dropToBaseElement(startNode)
	for i := 0; i < level; i++ {
		tempNode = *tempNode.up
	}
	return dropToBaseElement(*tempNode.next)
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
		authBuffer = authBuffer + sl.lists[nextLevel].tail.auth
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

	currentNode := sl.lists[sl.levels].head
	for {
		// If the transactions match drop to the baselevel and return the node
		if currentNode.tr == tr {
			for {
				if currentNode.down == nil {
					return currentNode.index, currentNode, true
				}
				currentNode = currentNode.down
			}
		}
		for {
			if currentNode.next != nil {
				if tr >= currentNode.next.tr {
					currentNode = currentNode.next
					break
				}
				if currentNode.down == nil {
					return -1, nil, false
				}
				currentNode = currentNode.down
			} else {
				if currentNode.down == nil {
					return -1, nil, false
				}
				currentNode = currentNode.down
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
		print("\n")
	}
	print("\n")
	LevelTester(sl)
	print("\n")
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

// Prints the number of levels of each element in a Skip List
func LevelTester(sl SkipList) {

	list := sl.lists[0]
	baseNode := list.head
	tempNode := baseNode
	levelCounter := 1
	fmt.Print("Levels:  -> ")
	for {
		if tempNode.up == nil {
			if baseNode.next == nil {
				fmt.Printf("%d\n", levelCounter)
				return
			}
			fmt.Printf("%d -> ", levelCounter)
			baseNode = baseNode.next
			tempNode = baseNode
			levelCounter = 1
			continue
		}
		tempNode = tempNode.up
		levelCounter++
	}
}
