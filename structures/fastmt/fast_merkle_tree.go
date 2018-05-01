package fastmt

import (
	"errors"

	. "github.com/SimoneStefani/thesis-algorithms/structures/common"
)

type FastMerkleTree struct {
	Root       *Node
	merkleRoot string
	Leaves     []*Node
}

type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	hash   string
	data   string
}

func NewFastMerkleTree(data []string) (*FastMerkleTree, error) {
	root, leaves, err := buildWithContent(data)

	if err != nil {
		return nil, err
	}

	t := &FastMerkleTree{
		Root:       root,
		merkleRoot: root.hash,
		Leaves:     leaves,
	}

	return t, nil
}

func buildWithContent(data []string) (*Node, []*Node, error) {
	if len(data) == 0 {
		return nil, nil, errors.New("Error: cannot construct tree with no content.")
	}

	var leaves []*Node
	for _, tr := range data {
		leaves = append(leaves, &Node{
			hash: HashTransaction(HashTransaction(tr)),
		})
	}

	if len(leaves)%2 == 1 {
		duplicate := &Node{
			hash: leaves[len(leaves)-1].hash,
			data: leaves[len(leaves)-1].data,
		}
		leaves = append(leaves, duplicate)
	}

	root := buildIntermediate(leaves)
	return root, leaves, nil
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
			hash:  HashTransaction(nl[left].hash + nl[right].hash),
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
