package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	SizeInternalNode int = 3 // 包含尾部的指针
	SizeLeafNode     int = 3 // 不包含下一个叶子节点的指针，下一个叶子节点保存在NextNode
)

type InternalNode struct {
	Pairs      [SizeInternalNode]*Element // element.value: internalNode, leafNode
	Size       int
	ParentNode *InternalNode
}

type LeafNode struct {
	Pairs      [SizeLeafNode]*Element // element.value: string
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

func (t *Tree) Print() {
	data, _ := json.Marshal(t)
	fmt.Println(string(data))
}

func (t *Tree) Delete(key string) error {
	// TODO
	return nil
}

func (t *Tree) Find(key string) (string, bool) {
	leafNode, err := t.SearchLeafNode(key)
	if err != nil {
		panic(err)
	}
	if leafNode == nil {
		panic("not found leafnode")
	}

	for i := 0; i < leafNode.Size; i++ {
		pair := leafNode.Pairs[i]
		if pair.Key == key {
			result, ok := pair.Value.(string)
			if !ok {
				panic("wrong leaf value")
			}
			return result, true
		}
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
		for i := 0; i < node.Size-1; i++ {
			pair := node.Pairs[i]
			if strings.Compare(key, pair.Key) < 0 {
				currentNode = pair.Value
				goto StartCurrentNode
			}
		}
		currentNode = node.Pairs[node.Size-1].Value
		break
	}
	if leafNode, ok := currentNode.(*LeafNode); ok {
		return leafNode, nil
	}
	return nil, errors.New("not found leaf node")
}

func (t *Tree) Insert(key, value string) error {
	if t.Root == nil {
		root := &LeafNode{
			Pairs: [SizeLeafNode]*Element{
				{Key: key, Value: value},
			},
			Size: 1,
		}
		t.Root = root
		return nil
	}
	leafNode, err := t.SearchLeafNode(key)
	if err != nil {
		return err
	}
	if leafNode == nil {
		return errors.New("Insert not found leaf node")
	}
	if leafNode.Size < SizeLeafNode {
		return leafNode.InsertToLeaf(key, value)
	} else {
		tempPair := t.MakeTempPairs(key, value, leafNode)
		m := SizeLeafNode / 2
		mKey := tempPair[m].Key
		rNode := new(LeafNode)
		copy(leafNode.Pairs[0:], tempPair[0:m])
		copy(rNode.Pairs[0:], tempPair[m:SizeLeafNode])
		rNode.NextNode = leafNode.NextNode
		leafNode.NextNode = rNode
		rNode.ParentNode = leafNode.ParentNode
		if leafNode.ParentNode == nil {
			root := &InternalNode{
				Pairs: [SizeInternalNode]*Element{
					{Key: mKey, Value: leafNode},
					{Key: "", Value: rNode},
				},
				Size: 2,
			}
			t.Root = root
			return nil
		} else {
			return t.InsertInternal(mKey, rNode, leafNode.ParentNode)
		}

	}
}

func (t *Tree) InsertInternal(key string, ptr interface{}, interN *InternalNode) error {
	if interN.Size < SizeInternalNode {
		return interN.InsertToInternal(key, ptr)
	} else {
		tempPair := t.MakeTempPairs(key, ptr, interN)
		m := SizeInternalNode / 2
		mKey := tempPair[m].Key
		rNode := new(InternalNode)
		copy(interN.Pairs[0:], tempPair[0:m])
		copy(rNode.Pairs[0:], tempPair[m:SizeInternalNode])
		rNode.ParentNode = interN.ParentNode
		if interN.ParentNode == nil {
			root := &InternalNode{
				Pairs: [SizeInternalNode]*Element{
					{Key: mKey, Value: interN},
					{Key: "", Value: rNode},
				},
				Size: 2,
			}
			t.Root = root
			return nil
		} else {
			return t.InsertInternal(mKey, rNode, interN.ParentNode)
		}
	}
}

func (t *Tree) MakeTempPairs(key string, value, node interface{}) []*Element {
	n, ok := node.(*LeafNode)
	temp := make([]*Element, 0, 100)
	if ok {

		var flag bool
		for i := 0; i < n.Size; i++ {
			ele := n.Pairs[i]
			if !flag && strings.Compare(key, ele.Key) < 0 {
				temp = append(temp, &Element{Key: key, Value: value})
				flag = true
			}
			temp = append(temp, ele)
		}
		return temp
	}
	n2, ok2 := node.(*InternalNode)
	if ok2 {
		var flag bool
		for i := 0; i < n2.Size; i++ {
			ele := n2.Pairs[i]
			if !flag && strings.Compare(key, ele.Key) < 0 {
				temp = append(temp, &Element{Key: key, Value: value})
				flag = true
			}
			temp = append(temp, ele)
		}
		return temp
	}
	return temp

}

func (leaf *LeafNode) InsertToLeaf(key, value string) error {
	temp := make([]*Element, 0, SizeLeafNode)
	var flag bool
	for i := 0; i < leaf.Size; i++ {
		ele := leaf.Pairs[i]
		if !flag && strings.Compare(key, ele.Key) < 0 {
			temp = append(temp, &Element{Key: key, Value: value})
			flag = true
		}
		temp = append(temp, ele)
	}
	if !flag {
		temp = append(temp, &Element{Key: key, Value: value})
	}
	for i := 0; i < leaf.Size+1; i++ {
		leaf.Pairs[i] = temp[i]
	}
	leaf.Size++
	if leaf.Size > SizeLeafNode {
		return errors.New("wrong leaf node count")
	}
	return nil
}

func (interN *InternalNode) InsertToInternal(key string, Ptr interface{}) error {
	temp := make([]*Element, 0, SizeInternalNode)
	var flag bool
	for i := 0; i < interN.Size; i++ {
		ele := interN.Pairs[i]
		if !flag && strings.Compare(key, ele.Key) < 0 {
			temp = append(temp, &Element{Key: key, Value: Ptr})
			flag = true
		}
		temp = append(temp, ele)
	}
	for i := 0; i < interN.Size+1; i++ {
		interN.Pairs[i] = temp[i]
	}
	interN.Size++
	if interN.Size > SizeInternalNode {
		return errors.New("wrong internal node count")
	}
	return nil
}
