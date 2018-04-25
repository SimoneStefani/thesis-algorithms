package main

import (
	"errors"
)

type MerkleTree struct {
	Root       *Node
	merkleRoot string
	Leafs      []*Node
}

type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	leaf   bool
	dup    bool
	Hash   string
	data   string
}

func NewTree(data []string) (*MerkleTree, error) {
	root, leafs, err := buildWithContent(data)

	if err != nil {
		return nil, err
	}

	t := &MerkleTree{
		Root:       root,
		merkleRoot: root.Hash,
		Leafs:      leafs,
	}

	return t, nil
}

func buildWithContent(data []string) (*Node, []*Node, error) {
	if len(data) == 0 {
		return nil, nil, errors.New("Error: cannot construct tree with no content.")
	}

	var leafs []*Node
	for _, tr := range data {
		leafs = append(leafs, &Node{
			Hash: hashTransaction(hashTransaction(tr)),
			data: tr,
			leaf: true,
		})
	}

	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			data: leafs[len(leafs)-1].data,
			leaf: true,
			dup:  true,
		}
		leafs = append(leafs, duplicate)
	}

	root := buildIntermediate(leafs)
	return root, leafs, nil
}

func buildIntermediate(nl []*Node) *Node {
	var nodes []*Node

	for i := 0; i < len(nl); i += 2 {

		var left, right int = i, i + 1
		if i+1 == len(nl) {
			right = i
		}

		n := &Node{
			Left:  nl[left],
			Right: nl[right],
			Hash:  hashTransaction(hashTransaction(nl[left].Hash + nl[right].Hash)),
		}
		nodes = append(nodes, n)

		nl[left].Parent = n
		nl[right].Parent = n

		if len(nl) == 2 {
			return n
		}
	}

	return buildIntermediate(nodes)
}
