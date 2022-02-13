package main

import "strings"

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

func (b *BPlusTree) Print() {}

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
		siblingNode := makeEmptyLeafNode()
		tempKeys := make([]Key, 0)
		tempPointers := make([]interface{}, 0)
		tempNode := makeEmptyLeafNode()
		tempNode.keys = tempKeys
		tempNode.pointers = tempPointers
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

}
func (b *BPlusTree) split(node *Node)                                     {}
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
