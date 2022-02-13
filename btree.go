package main

import "strings"

type BPlusTree struct {
	// 叶子结点指针最多数量, key半满条件: (leafMaxSize-1)/2
	leafMaxSize int
	// 非叶子结点指针最多数量, 指针半满条件: internalMaxSize/2，指针的数量叫做fanout
	internalMaxSize int
	// 0个或者2-n个子节点
	root *Node
}

type Node struct {
	// 是否是叶子结点
	isLeaf bool
	keys   []Key
	// 非叶子结点是*Node，叶子结点是*Record, 最后一个指针挪到了lastOrNextNode
	// 所以len(keys)=len(pointers)
	pointers []interface{}
	parent   *Node
	// 最后一个指针
	lastOrNextNode *Node
}

type Key string

func (k Key) Compare(target Key) int {
	return strings.Compare(k.toString(), target.toString())
}

func (k Key) toString() string {
	return string(k)
}

type Record struct {
	value []byte
}

func StartNewTree(leafMaxSize, internalMaxSize int) *BPlusTree {
	return &BPlusTree{
		leafMaxSize:     leafMaxSize,
		internalMaxSize: internalMaxSize,
	}
}

// 功能接口
func (b *BPlusTree) Insert(key, value string) {}
func (b *BPlusTree) Remove(key string)        {}
func (b *BPlusTree) Find(targetKey string) (string, bool) {
	tKey := Key(targetKey)
	currentNode := b.root
	for !currentNode.isLeaf {
		number := -1
		for i, key := range currentNode.keys {
			if tKey.Compare(key) == 0 {
				number = i + 1
				break
			} else if tKey.Compare(key) < 1 {
				number = i
				break
			}
		}
		var ok bool
		if number == -1 || number == len(currentNode.keys) {
			currentNode = currentNode.lastOrNextNode
		} else {
			currentNode, ok = currentNode.pointers[number].(*Node)
			if !ok {
				panic("should be *node")
			}
		}
	}
	if !currentNode.isLeaf {
		panic("should be leaf node")
	}
	return currentNode.findRecord(tKey)
}
func (b *BPlusTree) FindRange(start, end string) []string {
	return []string{}
}

func (b *BPlusTree) Print() {}

// 内部方法
func (b *BPlusTree) insertIntoLeaf(key, value string)                        {}
func (b *BPlusTree) insertIntoParent(oldNode, newNode *Node)                 {}
func (b *BPlusTree) split(node *Node)                                        {}
func (b *BPlusTree) delete(key string, pointer interface{})                  {}
func (b *BPlusTree) deleteEntry(node *Node, key string, pointer interface{}) {}

func (n *Node) findRecord(targetKey Key) (string, bool) {
	if !n.isLeaf {
		panic("not leaf node")
	}
	if len(n.keys) == 0 {
		return "", false
	}
	for i, key := range n.keys {
		if key.Compare(targetKey) == 0 {
			record := n.pointers[i]
			if value, ok := record.(*Record); ok {
				return string(value.value), true
			}
			panic("should be record")
		}
	}
	return "", false
}
