package btree

import (
	"errors"
	"fmt"
	"strings"
)

const defaultLeafMaxSize = 200
const defaultInternalMaxSize = 100

type BPlusTree struct {
	// 叶子节点key最多数量, key半满条件: leafMaxSize/2，指针的数量: leafMaxSize+1, 最后一个指针在node的lastOrNextNode
	leafMaxSize int
	// 非叶子节点key最多数量, key半满条件: internalMaxSize/2，指针的数量: internalMaxSize+1, 最后一个指针在node的lastOrNextNode
	internalMaxSize int
	// 0个或者2-n个子节点
	root *Node
	// 测试使用
	*Stat
}

type Node struct {
	// 是否是叶子节点
	isLeaf bool
	keys   []Key
	// 非叶子节点是*Node，叶子节点是*Entry, 最后一个指针挪到了lastOrNextNode
	// 所以len(keys)=len(pointers)
	pointers []interface{}
	parent   *Node
	// 最后一个指针
	lastOrNextNode *Node
}

// 查询统计（测试使用）
type Stat struct {
	// 查询遍历到的节点数
	count int64
	// 树的节点总数
	nodeCount int64
	// 数的高度
	level int64
}

func (b *Stat) incrCount() {
	b.count++
}

func (b *Stat) resetCount() {
	b.count = 0
}

func (b *Stat) GetCount() int64 {
	return b.count
}

func (b *Stat) GetLevel() int64 {
	return b.level
}

func (b *Stat) GetNodeCount() int64 {
	return b.nodeCount
}

type Key string

func (k Key) compare(target Key) int {
	return strings.Compare(k.toString(), target.toString())
}

func (k Key) toString() string {
	return string(k)
}

type Entry struct {
	value []byte
}

func (r *Entry) toValue() string {
	return string(r.value)
}

func StartNewTree(leafMaxSize, internalMaxSize int) (*BPlusTree, error) {
	if leafMaxSize < 3 || internalMaxSize < 3 {
		return nil, errors.New("need more than 2")
	}
	return &BPlusTree{
		leafMaxSize:     leafMaxSize,
		internalMaxSize: internalMaxSize,
		Stat:            new(Stat),
	}, nil
}

func StartDefaultNewTree() (*BPlusTree, error) {
	return &BPlusTree{
		leafMaxSize:     defaultLeafMaxSize,
		internalMaxSize: defaultInternalMaxSize,
		Stat:            new(Stat),
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

// 功能接口
func (b *BPlusTree) Insert(key, value string) {
	pointer := &Entry{
		value: []byte(value),
	}
	b.insert(Key(key), pointer)
}

func (b *BPlusTree) InsertByte(key string, value []byte) {
	pointer := &Entry{
		value: value,
	}
	b.insert(Key(key), pointer)
}

func (b *BPlusTree) Delete(key string) (value string, err error) {
	panic("wait for implement")
}

func (b *BPlusTree) DeleteByte(key string) (value []byte, err error) {
	panic("wait for implement")
}

func (b *BPlusTree) Find(targetKey string) (string, bool) {
	b.resetCount()
	leafNode := b.findLeafNode(Key(targetKey))
	value, ok := leafNode.findRecord(Key(targetKey))
	b.incrCount()
	return string(value), ok
}

func (b *BPlusTree) FindByte(targetKey string) ([]byte, bool) {
	b.resetCount()
	leafNode := b.findLeafNode(Key(targetKey))
	value, ok := leafNode.findRecord(Key(targetKey))
	b.incrCount()
	return value, ok
}

func (b *BPlusTree) FindRange(start, end string) []string {
	result := make([]string, 0)
	byteResult := b.FindRangeByte(start, end)
	if len(byteResult) == 0 {
		return result
	}
	for _, i := range byteResult {
		result = append(result, string(i))
	}
	return result
}

func (b *BPlusTree) FindRangeByte(start, end string) [][]byte {
	b.resetCount()
	result := make([][]byte, 0)
	startKey := Key(start)
	endKey := Key(end)
	if startKey.compare(endKey) == 1 {
		return result
	}
	leafNode := b.findLeafNode(Key(start))
	currentNode := leafNode
	for currentNode != nil {
		b.incrCount()
		for i, key := range currentNode.keys {
			if key.compare(startKey) >= 0 && key.compare(endKey) <= 0 {
				entry, ok := currentNode.pointers[i].(*Entry)
				if !ok {
					panic("should be *Entry")
				}
				result = append(result, entry.value)
			}
			if key.compare(endKey) == 1 {
				return result
			}
		}
		currentNode = currentNode.lastOrNextNode
	}
	return result
}

// 统计节点数量
func (b *BPlusTree) CountNode() {
	queue := make([]interface{}, 0)
	queue = append(queue, b.root)
	b.level = 0
	b.nodeCount = 0
	for len(queue) != 0 {
		b.level++
		size := len(queue)
		for i := 0; i < size; i++ {
			b.nodeCount++
			nodeI := queue[i]
			if nodeI == nil {
				continue
			}
			node, ok := nodeI.(*Node)
			if !ok {
				panic("should node")
			}
			if node.isLeaf {
				continue
			} else {
				if len(node.pointers) != 0 && !node.isLeaf {
					queue = append(queue, node.pointers...)
					queue = append(queue, node.lastOrNextNode)
				}
			}
		}

		if len(queue) > size {
			queue = queue[size:]
		} else {
			break
		}
	}
}

func (b *BPlusTree) Print() {
	fmt.Println("----------------------------------------------------------------------------------------------------start print tree")
	queue := make([]interface{}, 0)
	queue = append(queue, b.root)
	level := 0
	for len(queue) != 0 {
		level++
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
					entry := node.pointers[j]
					if p, ok := entry.(*Entry); ok {
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
		if len(queue) > size {
			queue = queue[size:]
		} else {
			break
		}
	}
	fmt.Println("----------------------------------------------------------------------------------------------------end print tree")
}

// 内部方法
func (b *BPlusTree) insert(targetKey Key, entry *Entry) {
	var leafNode *Node
	if b.root == nil {
		leafNode = makeEmptyLeafNode()
		b.root = leafNode
		leafNode.keys = append(leafNode.keys, targetKey)
		leafNode.pointers = append(leafNode.pointers, entry)
		return
	} else {
		leafNode = b.findLeafNode(targetKey)
	}
	if leafNode == nil {
		panic("should find leaf node")
	}
	if leafNode.updateRecord(targetKey, entry) {
		return
	}
	if len(leafNode.keys) < b.leafMaxSize {
		b.insertIntoLeaf(leafNode, targetKey, entry)
	} else {
		// split
		siblingNode := makeEmptyLeafNode()
		tempNode := makeEmptyLeafNode()
		tempNode.keys = append(tempNode.keys, leafNode.keys...)
		tempNode.pointers = append(tempNode.pointers, leafNode.pointers...)
		b.insertIntoLeaf(tempNode, targetKey, entry)
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
		b.incrCount()
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
func (b *BPlusTree) insertIntoLeaf(leafNode *Node, targetKey Key, value *Entry) {
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

func (b *BPlusTree) delete(key Key, pointer interface{})                  {}
func (b *BPlusTree) deleteEntry(node *Node, key Key, pointer interface{}) {}

func (n *Node) findRecord(targetKey Key) ([]byte, bool) {
	if !n.isLeaf {
		panic("should be leaf node")
	}
	if len(n.keys) == 0 {
		return []byte{}, false
	}
	// 可使用二分查找，待优化
	for i, key := range n.keys {
		if key.compare(targetKey) == 0 {
			entry := n.pointers[i]
			if value, ok := entry.(*Entry); ok {
				return value.value, true
			}
			panic("should be entry")
		}
	}
	return []byte{}, false
}

func (n *Node) updateRecord(targetKey Key, entry *Entry) bool {
	// 如果值已经存在则更新
	for i, k := range n.keys {
		if targetKey.compare(k) == 0 {
			r, ok := n.pointers[i].(*Entry)
			if !ok {
				panic("should be *entry")
			}
			r.value = entry.value
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
