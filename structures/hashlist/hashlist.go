package hashlist

import (
	"errors"

	. "github.com/Daynex/thesis-algorithms/structures/common"
)

type Node struct {
	prev *Node
	next *Node
	tr   string
}

type List struct {
	head *Node
	tail *Node
}

type HashList struct {
	list     *List
	headHash string
}

func NewHashList(data []string) (*HashList, error) {

	hashList, err := buildHashList(data)

	if err != nil {
		return nil, err
	}

	return hashList, nil
}

func buildHashList(data []string) (*HashList, error) {

	if len(data) == 0 {
		return nil, errors.New("Error: cannot construct tree with no content.")
	}

	list := &List{
		head: nil,
		tail: nil,
	}

	for _, tr := range data {
		list = insert(*list, tr)
	}
	return nil, nil
}

func insert(list List, tr string) *List {

	if list.head == nil {
		new := &Node{
			prev: nil,
			next: nil,
			tr:   HashTransaction(HashTransaction(tr)),
		}
		list.head = new
		list.tail = new
	} else {
		new := &Node{
			next: list.head,
			prev: nil,
			tr:   HashTransaction(HashTransaction(tr) + list.head.tr),
		}
		list.head = new
		new.next.prev = new
	}

	return &list
}
