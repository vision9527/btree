package btree

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
)

type BPlusTree struct {
	// 叶子结点key最多数量, key半满条件: leafMaxSize/2，指针的数量: leafMaxSize+1, 最后一个指针在node的lastOrNextNode
	leafMaxSize int
	// 非叶子结点key最多数量, key半满条件: internalMaxSize/2，指针的数量: internalMaxSize+1, 最后一个指针在node的lastOrNextNode
	internalMaxSize int
	// 0个或者2-n个子节点
	root *Node
	// 测试使用
	stat *Stat
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

// 查询统计（测试使用）
type Stat struct {
	// 查询遍历的节点数
	Count int64
}

func (s *Stat) incrCount() {
	atomic.AddInt64(&s.Count, 1)
}

func (s *Stat) resetCount() {
	atomic.StoreInt64(&s.Count, 0)
}

type Key string

func (k Key) compare(target Key) int {
	return strings.Compare(k.toString(), target.toString())
}

func (k Key) toString() string {
	return string(k)
}

type Record struct {
	value []byte
}

func (r *Record) toValue() string {
	return string(r.value)
}

func StartNewTree(leafMaxSize, internalMaxSize int) (*BPlusTree, error) {
	if leafMaxSize < 3 || internalMaxSize < 3 {
		return nil, errors.New("need more than 2")
	}
	return &BPlusTree{
		leafMaxSize:     leafMaxSize,
		internalMaxSize: internalMaxSize,
	}, nil
}

func StartDefaultNewTree() (*BPlusTree, error) {
	return &BPlusTree{
		leafMaxSize:     200,
		internalMaxSize: 100,
	}, nil
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

func (b *BPlusTree) SetStat(stat *Stat) {
	b.stat = new(Stat)
}

// 功能接口
func (b *BPlusTree) Insert(key, value string) {
	pointer := &Record{
		value: []byte(value),
	}
	b.insert(Key(key), pointer)
}

func (b *BPlusTree) InsertByte(key string, value []byte) {
	pointer := &Record{
		value: value,
	}
	b.insert(Key(key), pointer)
}

func (b *BPlusTree) Delete(key string) (value string, err error) {
	panic("need implement")
}

func (b *BPlusTree) DeleteByte(key string) (value []byte, err error) {
	panic("need implement")
}

func (b *BPlusTree) Find(targetKey string) (string, bool) {
	leafNode := b.findLeafNode(Key(targetKey))
	value, ok := leafNode.findRecord(Key(targetKey))
	b.IncrCount()
	return string(value), ok
}

func (b *BPlusTree) FindByte(targetKey string) ([]byte, bool) {
	leafNode := b.findLeafNode(Key(targetKey))
	value, ok := leafNode.findRecord(Key(targetKey))
	b.IncrCount()
	return value, ok
}

func (b *BPlusTree) FindRange(start, end string) []string {
	startKey := Key(start)
	endKey := Key(end)
	leafNode := b.findLeafNode(Key(start))
	currentNode := leafNode
	result := make([]string, 0)
	for currentNode != nil {
		b.IncrCount()
		for i, key := range currentNode.keys {
			if key.compare(startKey) >= 0 && key.compare(endKey) <= 0 {
				record, ok := currentNode.pointers[i].(*Record)
				if !ok {
					panic("should be *Record")
				}
				result = append(result, record.toValue())
			}
			if key.compare(endKey) == 1 {
				return result
			}
		}
		currentNode = currentNode.lastOrNextNode
	}
	return result
}

func (b *BPlusTree) Print() {
	fmt.Println("----------------------------------------------------------------------------------------------------start print tree")
	queue := make([]interface{}, 0)
	queue = append(queue, b.root)
	level := 1
	for len(queue) != 0 {
		size := len(queue)
		str := ""
		for i := 0; i < size; i++ {
			nodeI := queue[i]
			if nodeI == nil {
				str = strings.Trim(str, " &&")
				str = str + " --- "
				continue
			}
			node, ok := nodeI.(*Node)
			if !ok {
				panic("should node")
			}
			if !node.isLeaf {
				if str == "" {
					str = fmt.Sprintf("%v", node.keys)
				} else {
					str = str + "," + fmt.Sprintf("%v", node.keys)
				}
			} else {
				str = str + "("
				for j := 0; j < len(node.keys); j++ {
					key := node.keys[j]
					record := node.pointers[j]
					if p, ok := record.(*Record); ok {
						if j == 0 {
							str = str + fmt.Sprintf("%s|%s", key, string(p.value))
						} else {
							str = str + "," + fmt.Sprintf("%s|%s", key, string(p.value))
						}
					}
				}
				str = str + ") && "
			}
			if len(node.pointers) != 0 && !node.isLeaf {
				queue = append(queue, node.pointers...)
				queue = append(queue, node.lastOrNextNode)
				queue = append(queue, nil)
			}
		}
		str = strings.Trim(str, " &&")
		str = strings.Trim(str, "---")
		fmt.Printf("level %d: %s\n", level, str)
		level++
		if len(queue) > size {
			queue = queue[size:]
		} else {
			break
		}
	}
	fmt.Println("----------------------------------------------------------------------------------------------------end print tree")
}

// 内部方法
func (b *BPlusTree) insert(targetKey Key, record *Record) {
	var leafNode *Node
	if b.root == nil {
		leafNode = makeEmptyLeafNode()
		b.root = leafNode
		leafNode.keys = append(leafNode.keys, targetKey)
		leafNode.pointers = append(leafNode.pointers, record)
		return
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
		leafNode.keys = append(leafNode.keys, tempNode.keys[0:b.leafMaxSize/2+1]...)
		leafNode.pointers = append(leafNode.pointers, tempNode.pointers[0:b.leafMaxSize/2+1]...)
		siblingNode.keys = append(siblingNode.keys, tempNode.keys[b.leafMaxSize/2+1:]...)
		siblingNode.pointers = append(siblingNode.pointers, tempNode.pointers[b.leafMaxSize/2+1:]...)

		childKey := siblingNode.keys[0]
		b.insertIntoParent(leafNode, siblingNode, childKey)
	}
}

func (b *BPlusTree) findLeafNode(targetKey Key) *Node {
	tKey := Key(targetKey)
	currentNode := b.root
	for !currentNode.isLeaf {
		b.IncrCount()
		number := -1
		for i, key := range currentNode.keys {
			if tKey.compare(key) == 0 {
				number = i + 1
				break
			} else if tKey.compare(key) < 1 {
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
		if key.compare(targetKey) == 1 {
			number = i
			break
		}
	}
	if number == -1 {
		leafNode.keys = append(leafNode.keys, targetKey)
		leafNode.pointers = append(leafNode.pointers, value)
		return
	}
	leafNode.keys = append(leafNode.keys[:number], append([]Key{targetKey}, leafNode.keys[number:]...)...)
	leafNode.pointers = append(leafNode.pointers[:number], append([]interface{}{value}, leafNode.pointers[number:]...)...)
}
func (b *BPlusTree) insertIntoParent(oldNode, newNode *Node, childKey Key) {
	if oldNode.parent == nil {
		newRoot := makeEmptyInternalNode()
		newRoot.keys = append(newRoot.keys, childKey)
		newRoot.pointers = append(newRoot.pointers, oldNode)
		newRoot.lastOrNextNode = newNode
		oldNode.parent = newRoot
		newNode.parent = newRoot
		b.root = newRoot
		return
	}
	parentNode := oldNode.parent
	if len(parentNode.keys) < b.internalMaxSize {
		// insert (childKey, newNode) to parentNode after oldNode
		parentNode.insertNextAfterPrev(childKey, oldNode, newNode)
		newNode.parent = parentNode
		return
	} else {
		// split
		tempNode := parentNode.copy()
		tempNode.insertNextAfterPrev(childKey, oldNode, newNode)
		tempKeys := tempNode.keys
		tempPointers := tempNode.pointers
		tempPointers = append(tempPointers, tempNode.lastOrNextNode)
		parentNode.keys = make([]Key, 0)
		parentNode.pointers = make([]interface{}, 0)
		siblingParentNode := makeEmptyInternalNode()
		parentNode.keys = append(parentNode.keys, tempKeys[0:b.internalMaxSize/2]...)
		// parentNode.pointers = append(parentNode.pointers, tempPointers[0:b.internalMaxSize/2]...)
		for i := 0; i < b.internalMaxSize/2; i++ {
			childPointer := tempPointers[i]
			parentNode.pointers = append(parentNode.pointers, childPointer)
			if childPointer == newNode {
				newNode.parent = parentNode
			}
			if childPointer == oldNode {
				oldNode.parent = parentNode
			}
			p, ok := childPointer.(*Node)
			if !ok {
				panic("should be *node")
			}
			p.parent = parentNode
		}
		lst, ok := tempPointers[b.internalMaxSize/2].(*Node)
		if !ok {
			panic("should be *Node")
		}
		parentNode.lastOrNextNode = lst
		parentNode.lastOrNextNode.parent = parentNode
		if parentNode.lastOrNextNode == newNode {
			newNode.parent = parentNode
		}
		if parentNode.lastOrNextNode == oldNode {
			oldNode.parent = parentNode
		}

		siblingParentNode.keys = append(siblingParentNode.keys, tempKeys[b.internalMaxSize/2+1:]...)
		// siblingParentNode.pointers = append(siblingParentNode.pointers, tempPointers[b.internalMaxSize/2+1:b.internalMaxSize+1]...)
		for i := b.internalMaxSize/2 + 1; i < b.internalMaxSize+1; i++ {
			childPointer := tempPointers[i]
			siblingParentNode.pointers = append(siblingParentNode.pointers, childPointer)
			if childPointer == newNode {
				newNode.parent = siblingParentNode
			}
			if childPointer == oldNode {
				oldNode.parent = siblingParentNode
			}
			p, ok := childPointer.(*Node)
			if !ok {
				panic("should be *node")
			}
			p.parent = siblingParentNode
		}
		lst, ok = tempPointers[b.internalMaxSize+1].(*Node)
		if !ok {
			panic("should be *Node")
		}
		siblingParentNode.lastOrNextNode = lst
		siblingParentNode.lastOrNextNode.parent = siblingParentNode
		if siblingParentNode.lastOrNextNode == newNode {
			newNode.parent = siblingParentNode
		}
		if siblingParentNode.lastOrNextNode == oldNode {
			oldNode.parent = siblingParentNode
		}

		childKeyTwo := tempKeys[b.internalMaxSize/2]
		b.insertIntoParent(parentNode, siblingParentNode, childKeyTwo)
	}

}

func (b *BPlusTree) IncrCount() {
	if b.stat != nil {
		b.stat.incrCount()
	}
}

func (b *BPlusTree) ResetCount() {
	if b.stat != nil {
		b.stat.resetCount()
	}
}

func (b *BPlusTree) GetCount() int64 {
	if b.stat != nil {
		return b.stat.Count
	}
	return 0
}

func (b *BPlusTree) delete(key Key, pointer interface{})                  {}
func (b *BPlusTree) deleteEntry(node *Node, key Key, pointer interface{}) {}

func (n *Node) findRecord(targetKey Key) ([]byte, bool) {
	if !n.isLeaf {
		panic("not leaf node")
	}
	if len(n.keys) == 0 {
		return []byte{}, false
	}
	for i, key := range n.keys {
		if key.compare(targetKey) == 0 {
			record := n.pointers[i]
			if value, ok := record.(*Record); ok {
				return value.value, true
			}
			panic("should be record")
		}
	}
	return []byte{}, false
}

func (n *Node) updateRecord(targetKey Key, record *Record) bool {
	// 如果值已经存在则更新
	for i, k := range n.keys {
		if targetKey.compare(k) == 0 {
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
		// prev 在lastOrNextNode
		n.keys = append(n.keys, childKey)
		n.pointers = append(n.pointers, prev)
		n.lastOrNextNode = next
	} else {
		n.keys = append(n.keys[:number], append([]Key{childKey}, n.keys[number:]...)...)
		n.pointers = append(n.pointers[:number+1], append([]interface{}{next}, n.pointers[number+1:]...)...)
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
