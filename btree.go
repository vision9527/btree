package main

import (
	"errors"
	"strings"
)

const (
	SizeInternalNode = 14
	SizeLeafNode     = 10
)

type InternalNode struct {
	Pairs      [SizeInternalNode]Element // element.value: internalNode, leafNode
	Size       int
	LastNode   interface{} // lastNode:  internalNode, leafNode
	ParentNode *InternalNode
}

type LeafNode struct {
	Pairs      [SizeLeafNode]Element // element.value: string
	Size       int
	NextNode   *LeafNode
	ParentNode *InternalNode
}

type Element struct {
	Key   string
	Value interface{}
}

type Tree struct {
	Root interface{}
}

func NewTree() *Tree {
	return new(Tree)
}

func (t *Tree) Find(key string) (string, bool) {
	leafNode, err := t.SearchLeafNode(key)
	if err != nil {
		panic(err)
	}
	if leafNode == nil {
		panic("not found leafnode")
	}

	return "", false
}

func (t *Tree) SearchLeafNode(key string) (*LeafNode, error) {
	if t.Root == nil {
		return nil, nil
	}
	if leafNode, ok := t.Root.(*LeafNode); ok {
		return leafNode, nil
	}
	currentNode := t.Root
	for {
	StartCurrentNode:
		if node, ok := currentNode.(*LeafNode); ok {
			return node, nil
		}
		node, ok := currentNode.(*InternalNode)
		if !ok {
			return nil, errors.New("wrong node type")
		}
		for _, pair := range node.Pairs {
			if strings.Compare(key, pair.Key) < 0 {
				currentNode = pair.Value
				goto StartCurrentNode
			}
		}
		currentNode = node.Pairs[SizeInternalNode-1].Value
		break
	}
	if leafNode, ok := currentNode.(*LeafNode); ok {
		return leafNode, nil
	}
	return nil, errors.New("not found leaf node")
}

func (t *Tree) Insert(key, value string) error {
	return nil
}

func (t *Tree) InsertInternal(key, value string) error {
	return nil
}

func (t *Tree) InsertLeaf(key, value string) error {
	return nil
}
