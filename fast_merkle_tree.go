package main

import (
	"errors"
)

type FastMerkleTree struct {
	Root       *FMTNode
	merkleRoot string
	Leafs      []*FMTNode
}

type FMTNode struct {
	Parent *FMTNode
	Left   *FMTNode
	Right  *FMTNode
	leaf   bool
	dup    bool
	Hash   string
	data   string
}

func NewFastMerkleTree(data []string) (*FastMerkleTree, error) {
	root, leafs, err := buildWithContentF(data)

	if err != nil {
		return nil, err
	}

	t := &FastMerkleTree{
		Root:       root,
		merkleRoot: root.Hash,
		Leafs:      leafs,
	}

	return t, nil
}

func buildWithContentF(data []string) (*FMTNode, []*FMTNode, error) {
	if len(data) == 0 {
		return nil, nil, errors.New("Error: cannot construct tree with no content.")
	}

	var leafs []*FMTNode
	for _, tr := range data {
		leafs = append(leafs, &FMTNode{
			Hash: hashTransaction(hashTransaction(tr)),
			data: tr,
			leaf: true,
		})
	}

	if len(leafs)%2 == 1 {
		duplicate := &FMTNode{
			Hash: leafs[len(leafs)-1].Hash,
			data: leafs[len(leafs)-1].data,
			leaf: true,
			dup:  true,
		}
		leafs = append(leafs, duplicate)
	}

	root := buildIntermediateF(leafs)
	return root, leafs, nil
}

func buildIntermediateF(nl []*FMTNode) *FMTNode {
	var nodes []*FMTNode

	for i := 0; i < len(nl); i += 2 {

		var left, right int = i, i + 1
		if i+1 == len(nl) {
			right = i
		}

		n := &FMTNode{
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

	return buildIntermediateF(nodes)
}
