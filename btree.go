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
	root *node
	// 测试使用
	*stat
}

func StartNewTree(leafMaxSize, internalMaxSize int) (*BPlusTree, error) {
	if leafMaxSize < 3 || internalMaxSize < 3 {
		return nil, errors.New("need more than 2")
	}
	return &BPlusTree{
		leafMaxSize:     leafMaxSize,
		internalMaxSize: internalMaxSize,
		stat:            new(stat),
	}, nil
}

func StartDefaultNewTree() (*BPlusTree, error) {
	return &BPlusTree{
		leafMaxSize:     defaultLeafMaxSize,
		internalMaxSize: defaultInternalMaxSize,
		stat:            new(stat),
	}, nil
}

// 功能接口
func (b *BPlusTree) Insert(ky string, value interface{}) {
	pointer := &entry{
		value: value,
	}
	b.insert(key(ky), pointer)
}

func (b *BPlusTree) Delete(ky string) (value interface{}, err error) {
	panic("wait for implement")
}

func (b *BPlusTree) Find(targetKey string) (interface{}, bool) {
	b.resetCount()
	leafNode := b.findLeafNode(key(targetKey))
	value, ok := leafNode.findRecord(key(targetKey))
	b.incrCount()
	return value, ok
}

func (b *BPlusTree) FindRange(start, end string) []interface{} {
	b.resetCount()
	result := make([]interface{}, 0)
	startKey := key(start)
	endKey := key(end)
	if startKey.compare(endKey) == 1 {
		return result
	}
	leafNode := b.findLeafNode(key(start))
	currentNode := leafNode
	for currentNode != nil {
		b.incrCount()
		for i, ky := range currentNode.keys {
			if ky.compare(startKey) >= 0 && ky.compare(endKey) <= 0 {
				et, ok := currentNode.pointers[i].(*entry)
				if !ok {
					panic("should be *entry")
				}
				result = append(result, et.value)
			}
			if ky.compare(endKey) == 1 {
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
			nd, ok := nodeI.(*node)
			if !ok {
				panic("should node")
			}
			if nd.isLeaf {
				continue
			} else {
				if len(nd.pointers) != 0 && !nd.isLeaf {
					queue = append(queue, nd.pointers...)
					queue = append(queue, nd.lastOrNextNode)
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
	fmt.Printf("InternalMaxSize:%d LeafMaxSize: %d\n", b.internalMaxSize, b.leafMaxSize)
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
			nd, ok := nodeI.(*node)
			if !ok {
				panic("should node")
			}
			if !nd.isLeaf {
				if str == "" {
					str = fmt.Sprintf("%v", nd.keys)
				} else {
					if strings.HasSuffix(str, " --- ") {
						str = str + fmt.Sprintf("%v", nd.keys)
					} else {
						str = str + "," + fmt.Sprintf("%v", nd.keys)
					}

				}
			} else {
				str = str + "("
				for j := 0; j < len(nd.keys); j++ {
					ky := nd.keys[j]
					et := nd.pointers[j]
					if _, ok := et.(*entry); ok {
						if j == 0 {
							str = str + ky.toString()
						} else {
							str = str + "," + ky.toString()
						}
					}
				}
				str = str + ") && "
			}
			if len(nd.pointers) != 0 && !nd.isLeaf {
				queue = append(queue, nd.pointers...)
				queue = append(queue, nd.lastOrNextNode)
				queue = append(queue, nil)
			}
		}
		str = strings.Trim(str, " &&")
		str = strings.Trim(str, "---")
		fmt.Printf("Level %d: %s\n", level, str)
		if len(queue) > size {
			queue = queue[size:]
		} else {
			break
		}
	}
	fmt.Println("----------------------------------------------------------------------------------------------------end print tree")
}

// 内部方法
func (b *BPlusTree) makeEmptyLeafNode() *node {
	return &node{
		isLeaf:   true,
		keys:     make([]key, 0, b.leafMaxSize),
		pointers: make([]interface{}, 0, b.internalMaxSize),
	}
}

func (b *BPlusTree) makeEmptyInternalNode() *node {
	return &node{
		isLeaf:   false,
		keys:     make([]key, 0, b.leafMaxSize),
		pointers: make([]interface{}, 0, b.internalMaxSize),
	}
}

func (b *BPlusTree) insert(targetKey key, et *entry) {
	var leafNode *node
	if b.root == nil {
		leafNode = b.makeEmptyLeafNode()
		b.root = leafNode
		leafNode.keys = append(leafNode.keys, targetKey)
		leafNode.pointers = append(leafNode.pointers, et)
		return
	} else {
		leafNode = b.findLeafNode(targetKey)
	}
	if leafNode == nil {
		panic("should find leaf node")
	}
	if leafNode.updateRecord(targetKey, et) {
		return
	}
	if len(leafNode.keys) < b.leafMaxSize {
		b.insertIntoLeaf(leafNode, targetKey, et)
	} else {
		// split
		siblingNode := b.makeEmptyLeafNode()
		tempNode := b.makeEmptyLeafNode()
		tempNode.keys = append(tempNode.keys, leafNode.keys...)
		tempNode.pointers = append(tempNode.pointers, leafNode.pointers...)
		b.insertIntoLeaf(tempNode, targetKey, et)
		siblingNode.lastOrNextNode = leafNode.lastOrNextNode
		leafNode.lastOrNextNode = siblingNode
		leafNode.keys = make([]key, 0)
		leafNode.pointers = make([]interface{}, 0)
		leafNode.keys = append(leafNode.keys, tempNode.keys[0:b.leafMaxSize/2+1]...)
		leafNode.pointers = append(leafNode.pointers, tempNode.pointers[0:b.leafMaxSize/2+1]...)
		siblingNode.keys = append(siblingNode.keys, tempNode.keys[b.leafMaxSize/2+1:]...)
		siblingNode.pointers = append(siblingNode.pointers, tempNode.pointers[b.leafMaxSize/2+1:]...)

		childKey := siblingNode.keys[0]
		b.insertIntoParent(leafNode, siblingNode, childKey)
	}
}

func (b *BPlusTree) findFirstLeafNode() *node {
	currentNode := b.root
	for currentNode != nil {
		if currentNode.isLeaf {
			break
		}
		pointer := currentNode.pointers[0]
		nd, ok := pointer.(*node)
		if !ok {
			panic("should be *node")
		}
		currentNode = nd
	}
	return currentNode
}

func (b *BPlusTree) findLeafNode(targetKey key) *node {
	tKey := key(targetKey)
	currentNode := b.root
	for !currentNode.isLeaf {
		b.incrCount()
		number := -1
		for i, ky := range currentNode.keys {
			if tKey.compare(ky) == 0 {
				number = i + 1
				break
			} else if tKey.compare(ky) < 1 {
				number = i
				break
			}
		}
		var ok bool
		if number == -1 || number == len(currentNode.keys) {
			currentNode = currentNode.lastOrNextNode
		} else {
			currentNode, ok = currentNode.pointers[number].(*node)
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
func (b *BPlusTree) insertIntoLeaf(leafNode *node, targetKey key, value *entry) {
	number := -1
	for i, ky := range leafNode.keys {
		if ky.compare(targetKey) == 1 {
			number = i
			break
		}
	}
	if number == -1 {
		leafNode.keys = append(leafNode.keys, targetKey)
		leafNode.pointers = append(leafNode.pointers, value)
		return
	}
	leafNode.keys = append(leafNode.keys[:number], append([]key{targetKey}, leafNode.keys[number:]...)...)
	leafNode.pointers = append(leafNode.pointers[:number], append([]interface{}{value}, leafNode.pointers[number:]...)...)
}
func (b *BPlusTree) insertIntoParent(oldNode, newNode *node, childKey key) {
	if oldNode.parent == nil {
		newRoot := b.makeEmptyInternalNode()
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
		parentNode.keys = make([]key, 0)
		parentNode.pointers = make([]interface{}, 0)
		siblingParentNode := b.makeEmptyInternalNode()
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
			p, ok := childPointer.(*node)
			if !ok {
				panic("should be *node")
			}
			p.parent = parentNode
		}
		lst, ok := tempPointers[b.internalMaxSize/2].(*node)
		if !ok {
			panic("should be *node")
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
			p, ok := childPointer.(*node)
			if !ok {
				panic("should be *node")
			}
			p.parent = siblingParentNode
		}
		lst, ok = tempPointers[b.internalMaxSize+1].(*node)
		if !ok {
			panic("should be *node")
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

func (b *BPlusTree) delete(ky key, pointer interface{})                {}
func (b *BPlusTree) deleteEntry(nd *node, ky key, pointer interface{}) {}
