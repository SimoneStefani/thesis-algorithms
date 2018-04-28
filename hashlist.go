package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

type LNode struct {
	prev *LNode
	next *LNode
	tr   string
}

type List struct {
	head *LNode
	tail *LNode
}

type HashList struct {
	list     *List
	headHash string
}

func NewHashList(data []string) (*HashList, error) {

	hashList, err := BuildHashList(data)

	if err != nil {
		return nil, err
	}

	return hashList, nil
}

func BuildHashList(data []string) (*HashList, error) {

	if len(data) == 0 {
		return nil, errors.New("Error: cannot construct tree with no content.")
	}

	list := &List{
		head: nil,
		tail: nil,
	}

	for _, tr := range data {
		list = Insert(*list, tr)
	}
	return nil, nil
}

func Insert(list List, tr string) *List {

	if list.head == nil {
		new := &LNode{
			prev: nil,
			next: nil,
			tr:   hashTransaction(hashTransaction(tr)),
		}
		list.head = new
		list.tail = new
	} else {
		new := &LNode{
			next: list.head,
			prev: nil,
			tr:   hashTransaction(hashTransaction(tr) + list.head.tr),
		}
		list.head = new
		new.next.prev = new
	}

	return &list
}

func Show(list List) {
	l := list.head
	for l != nil {
		fmt.Printf("%+v -> ", l.tr[0:6])
		l = l.next
	}
	fmt.Println("nil")
}

func ShowReverse(list List) {
	l := list.tail
	for l != nil {
		fmt.Printf("%+v -> ", l.tr[0:6])
		l = l.prev
	}
	fmt.Println("nil")
}

func GetTail(hl *HashList) string {
	return hl.list.tail.tr
}

func hashTransaction(tr string) string {
	h := sha256.Sum256([]byte(tr))
	return base64.StdEncoding.EncodeToString(h[:])
}
