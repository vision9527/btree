package main

import (
	"fmt"
	"strings"
)

type BPlusTree struct {
	// 叶子结点key最多数量, key半满条件: leafMaxSize/2
	leafMaxSize int
	// 非叶子结点key最多数量, key半满条件: internalMaxSize/2，指针的数量叫做fanout
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

func (r *Record) ToValue() string {
	return string(r.value)
}

func StartNewTree(leafMaxSize, internalMaxSize int) *BPlusTree {
	return &BPlusTree{
		leafMaxSize:     leafMaxSize,
		internalMaxSize: internalMaxSize,
	}
}

func makeEmptyLeafNode() *Node {
	return &Node{
		isLeaf:   true,
		keys:     make([]Key, 0),
		pointers: make([]interface{}, 0),
	}
}

func makeEmptyInternalNode() *Node {
	return &Node{
		isLeaf:   false,
		keys:     make([]Key, 0),
		pointers: make([]interface{}, 0),
	}
}

// 功能接口
func (b *BPlusTree) Insert(key, value string) {
	pointer := &Record{
		value: []byte(value),
	}
	b.insert(Key(key), pointer)
}
func (b *BPlusTree) Remove(key string) {}
func (b *BPlusTree) Find(targetKey string) (string, bool) {
	leafNode := b.findLeafNode(Key(targetKey))
	return leafNode.findRecord(Key(targetKey))
}
func (b *BPlusTree) FindRange(start, end string) []string {
	return []string{}
}

func (b *BPlusTree) Print() {
	b.print(b.root)
}

func (b *BPlusTree) print(node *Node) {
	fmt.Println("node keys:", node.keys)
	if len(node.pointers) != 0 {
		for _, n := range node.pointers {
			if p, ok := n.(*Node); ok {
				b.print(p)
			}
			if p, ok := n.(*Record); ok {
				fmt.Println("record: ", string(p.value))
			}
		}
	}
	if node.lastOrNextNode != nil && !node.isLeaf {
		b.print(node.lastOrNextNode)
	}
}

// 内部方法
func (b *BPlusTree) insert(targetKey Key, record *Record) {
	var leafNode *Node
	if b.root == nil {
		leafNode = makeEmptyLeafNode()
		b.root = leafNode
	} else {
		leafNode = b.findLeafNode(targetKey)
	}
	if leafNode == nil {
		panic("should find leaf node")
	}
	if leafNode.updateRecord(targetKey, record) {
		return
	}

	if len(leafNode.keys) < b.leafMaxSize {
		b.insertIntoLeaf(leafNode, targetKey, record)
	} else {
		// split
		siblingNode := makeEmptyLeafNode()
		tempNode := makeEmptyLeafNode()
		tempNode.keys = append(tempNode.keys, leafNode.keys...)
		tempNode.pointers = append(tempNode.pointers, leafNode.pointers...)

		b.insertIntoLeaf(tempNode, targetKey, record)
		siblingNode.lastOrNextNode = leafNode.lastOrNextNode
		leafNode.lastOrNextNode = siblingNode
		leafNode.keys = make([]Key, 0)
		leafNode.pointers = make([]interface{}, 0)
		leafNode.keys = append(leafNode.keys, tempNode.keys[0:b.leafMaxSize/2]...)
		leafNode.pointers = append(leafNode.pointers, tempNode.pointers[0:b.leafMaxSize/2]...)

		siblingNode.keys = append(siblingNode.keys, tempNode.keys[b.leafMaxSize/2:]...)
		siblingNode.pointers = append(siblingNode.pointers, tempNode.pointers[b.leafMaxSize/2:]...)
		childKey := siblingNode.keys[0]
		b.insertIntoParent(leafNode, siblingNode, childKey)
	}
}

func (b *BPlusTree) findLeafNode(targetKey Key) *Node {
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
	return currentNode
}
func (b *BPlusTree) insertIntoLeaf(leafNode *Node, targetKey Key, value *Record) {
	number := -1
	for i, key := range leafNode.keys {
		if key.Compare(targetKey) == 1 {
			number = i
			break
		}
	}
	if number == -1 {
		leafNode.keys = append(leafNode.keys, targetKey)
		leafNode.pointers = append(leafNode.pointers, value)
		return
	}
	tempKeys := make([]Key, 0)
	tempPointers := make([]interface{}, 0)
	tempKeys = append(tempKeys, leafNode.keys[0:number]...)
	tempKeys = append(tempKeys, targetKey)
	tempKeys = append(tempKeys, leafNode.keys[number:]...)
	tempPointers = append(tempPointers, leafNode.pointers[0:number]...)
	tempPointers = append(tempPointers, value)
	tempPointers = append(tempPointers, leafNode.pointers[number:]...)
	leafNode.keys = tempKeys
	leafNode.pointers = tempPointers
}
func (b *BPlusTree) insertIntoParent(oldNode, newNode *Node, childKey Key) {
	fmt.Println("insertIntoParent", childKey)
	if oldNode.parent == nil {
		newRoot := makeEmptyInternalNode()
		newRoot.keys = append(newRoot.keys, childKey)
		newRoot.pointers = append(newRoot.pointers, oldNode)
		newRoot.lastOrNextNode = newNode
		oldNode.parent = newRoot
		newRoot.parent = newRoot
		b.root = newRoot
		return
	}
	parentNode := oldNode.parent
	if len(parentNode.keys) < b.internalMaxSize {
		fmt.Println("split len(parentNode.keys):", len(parentNode.keys), b.internalMaxSize)
		// insert (childKey, newNode) to parentNode after oldNode
		parentNode.insertNextAfterPrev(childKey, oldNode, newNode)
		return
	} else {
		// split
		fmt.Println("split")
		tempNode := parentNode.copy()
		tempNode.insertNextAfterPrev(childKey, oldNode, newNode)
		parentNode.keys = make([]Key, 0)
		parentNode.pointers = make([]interface{}, 0)
		parentNodeLastOrNextNodePointer := parentNode.lastOrNextNode
		siblingNode := makeEmptyInternalNode()
		parentNode.keys = append(parentNode.keys, tempNode.keys[0:b.internalMaxSize/2]...)
		parentNode.pointers = append(parentNode.pointers, tempNode.pointers[0:b.internalMaxSize/2]...)
		lst, ok := tempNode.pointers[b.internalMaxSize/2].(*Node)
		if !ok {
			panic("should be *Node")
		}
		parentNode.lastOrNextNode = lst
		siblingNode.keys = append(siblingNode.keys, tempNode.keys[b.internalMaxSize/2+1:]...)
		siblingNode.pointers = append(siblingNode.pointers, tempNode.pointers[b.internalMaxSize/2:]...)
		siblingNode.lastOrNextNode = parentNodeLastOrNextNodePointer
		siblingNode.parent = parentNode.parent
		childKeyTwo := tempNode.keys[b.internalMaxSize/2+1]
		b.insertIntoParent(parentNode, siblingNode, childKeyTwo)
	}

}

func (b *BPlusTree) delete(key Key, pointer interface{})                  {}
func (b *BPlusTree) deleteEntry(node *Node, key Key, pointer interface{}) {}

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

func (n *Node) updateRecord(targetKey Key, record *Record) bool {
	// 如果值已经存在则更新
	for i, k := range n.keys {
		if targetKey.Compare(k) == 0 {
			r, ok := n.pointers[i].(*Record)
			if !ok {
				panic("should be *record")
			}
			r.value = record.value
			n.pointers[i] = r
			return true
		}
	}
	return false
}

func (n *Node) insertNextAfterPrev(childKey Key, prev, next *Node) {
	number := -1
	for i, p := range n.pointers {
		if p == prev {
			number = i
			break
		}
	}

	if number == -1 {
		n.keys = append(n.keys, childKey)
		n.pointers = append(n.pointers, prev)
		n.lastOrNextNode = next
	} else {
		tempKeys := make([]Key, 0)
		tempKeys = append(tempKeys, n.keys...)
		tempPointers := make([]interface{}, 0)
		tempPointers = append(tempPointers, n.pointers...)
		if len(n.keys) < number+2 {
			n.keys = append(n.keys, Key(""))
		}
		n.keys[number+1] = childKey
		copy(n.keys[number+2:], tempKeys[number+1:])
		if len(n.pointers) < number+2 {
			n.pointers = append(n.pointers, nil)
		}
		n.pointers[number+1] = next
		copy(n.pointers[number+2:], tempPointers[number+1:])
	}

}

func (n *Node) copy() *Node {
	node := &Node{
		isLeaf:         n.isLeaf,
		keys:           make([]Key, 0),
		pointers:       make([]interface{}, 0),
		parent:         n.parent,
		lastOrNextNode: n.lastOrNextNode,
	}
	node.keys = append(node.keys, n.keys...)
	node.pointers = append(node.pointers, n.pointers...)
	return node
}
