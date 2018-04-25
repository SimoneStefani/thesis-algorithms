package main

import (
	"crypto/sha256"
	"encoding/base64"
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

func (L *List) Insert(tr string) {
	list := &LNode{
		next: L.head,
		tr:   tr,
	}
	if L.head != nil {
		L.head.prev = list
	}
	L.head = list

	l := L.head
	for l.next != nil {
		l = l.next
	}
	L.tail = l
}

func (l *List) Show() {
	list := l.head
	for list != nil {
		fmt.Printf("%+v -> ", list.tr[0:6])
		list = list.next
	}
	fmt.Println()
}

func (l *List) BuildList(data []string) {
	for _, tr := range data {
		tailHash := ""
		if l.tail != nil {
			tailHash = l.head.tr
		}

		l.Insert(hashTransaction(tailHash + hashTransaction(tr)))
	}
}

func (l *List) GetTail() string {
	return l.tail.tr
}

func hashTransaction(tr string) string {
	h := sha256.Sum256([]byte(tr))
	return base64.StdEncoding.EncodeToString(h[:])
}
