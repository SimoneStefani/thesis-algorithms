package hashlist

import (
	"errors"

	. "github.com/SimoneStefani/thesis-algorithms/structures/common"
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

func VerifyTransaction(tr string, list []string) (string, []string, bool, error) {

	pos, err := Includes(tr, list)

	if err != nil {
		return "", nil, false, err
	}
	path, hl := computePath(pos, list)

	return hl.headHash, path, CheckPath(tr, hl.headHash, path), nil
}

func CheckPath(tr string, headHash string, path []string) bool {

	var hash string
	if HashTransaction(HashTransaction(tr)) == path[0] {
		hash = path[0]
	} else {
		hash = HashTransaction(HashTransaction(tr) + path[0])
	}

	for i := 1; i < len(path); i++ {
		hash = HashTransaction(path[i] + hash)
	}

	return hash == headHash
}

func (hl *HashList) Length() int {
	if hl.list.head == hl.list.tail {
		return 1
	}

	count := 0
	current := hl.list.head
	for {
		count = count + 1
		if current.next == nil {
			break
		}
		current = current.next
	}

	return count
}

func computePath(pos int, list []string) ([]string, *HashList) {

	var path []string
	temp := &List{
		head: nil,
		tail: nil,
	}

	for i, tr := range list {
		temp = insert(*temp, tr)
		if pos < i {
			path = append(path, HashTransaction(tr))
		} else if pos == i+1 || pos == 0 {
			path = append(path, temp.head.tr)
		}
	}

	hl := &HashList{
		headHash: temp.head.tr,
		list:     temp,
	}

	return path, hl
}

func buildHashList(data []string) (*HashList, error) {

	if len(data) == 0 {
		return nil, errors.New("Error: cannot construct hashlist with no content.")
	}

	list := &List{
		head: nil,
		tail: nil,
	}

	for _, tr := range data {
		list = insert(*list, tr)
	}

	hl := &HashList{
		list:     list,
		headHash: list.head.tr,
	}

	return hl, nil
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
