package fastmt

import (
	"errors"

	. "github.com/Daynex/thesis-algorithms/structures/common"
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

type VerificationNode struct {
	hash   string
	isLeft bool
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

func VerifyTransaction(tr string, list []string) (string, []VerificationNode, bool, error) {

	pos, err := Includes(tr, list)

	if err != nil {
		return "", nil, false, err
	}

	tree, err := NewFastMerkleTree(list)
	path := computeMerklePath(pos, tree)

	return tree.merkleRoot, path, CheckPath(tr, tree.merkleRoot, path), err
}

func CheckPath(tr string, roothash string, path []VerificationNode) bool {

	hash := HashTransaction(HashTransaction(tr))

	for _, node := range path {
		if node.isLeft {
			hash = HashTransaction(node.hash + hash)
		} else {
			hash = HashTransaction(hash + node.hash)
		}
	}

	return hash == roothash
}

func computeMerklePath(pos int, tree *FastMerkleTree) []VerificationNode {

	node := tree.Leaves[pos]

	var path []VerificationNode
	var temp VerificationNode

	for {
		if node.Parent == nil {
			break
		}
		if isLeftChild(node) {
			temp = VerificationNode{
				hash:   node.Parent.Right.hash,
				isLeft: false,
			}
		} else {
			temp = VerificationNode{
				hash:   node.Parent.Left.hash,
				isLeft: true,
			}
		}
		path = append(path, temp)
		node = node.Parent
	}

	return path
}

func isLeftChild(node *Node) bool {
	return node.Parent.Left.hash == node.hash
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
