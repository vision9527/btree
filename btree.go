package btree

import (
	"errors"
	"fmt"
	"strings"
)

const defaultLeafMaxSize = 100
const defaultInternalMaxSize = 100

type BPlusTree struct {
	// 叶子节点key最多数量:leafMaxSize, 叶子节点指针最多数量(degree): leafMaxSize+1
	// key半满条件: leafMaxSize/2
	// 指针的数量半满条件: leafMaxSize/2+1
	// 最后一个指针在node的lastOrNextNode
	leafMaxSize int
	// 内部节点key最多数量:internalMaxSize, 内部节点指针最多数量(degree): internalMaxSize+1
	// key半满条件: internalMaxSize/2
	// 指针的数量半满条件: internalMaxSize/2+1
	// 最后一个指针在node的lastOrNextNode
	internalMaxSize int
	// 当根节点是非叶子节点时，key的数量可以是1-internalMaxSize，不需要满足半满条件
	root *node
	// 测试使用
	*stat
}

func StartNewTree(leafMaxSize, internalMaxSize int) (*BPlusTree, error) {
	if leafMaxSize < 3 || internalMaxSize < 3 {
		return nil, errors.New("need more than 2")
	}
	tree := &BPlusTree{
		leafMaxSize:     leafMaxSize,
		internalMaxSize: internalMaxSize,
		stat:            new(stat),
	}
	tree.root = tree.makeEmptyLeafNode()
	return tree, nil
}

func StartDefaultNewTree() (*BPlusTree, error) {
	tree := &BPlusTree{
		leafMaxSize:     defaultLeafMaxSize,
		internalMaxSize: defaultInternalMaxSize,
		stat:            new(stat),
	}
	tree.root = tree.makeEmptyLeafNode()
	return tree, nil
}

// 功能接口
func (b *BPlusTree) Insert(ky string, value interface{}) {
	pointer := &entry{
		value: value,
	}
	b.insert(key(ky), pointer)
}

func (b *BPlusTree) Delete(ky string) (interface{}, bool) {
	return b.delete(key(ky))
}

func (b *BPlusTree) Find(targetKey string) (interface{}, bool) {
	b.resetCount()
	leafNode := b.findLeafNode(key(targetKey))
	et, ok := leafNode.findRecord(key(targetKey))
	b.incrCount()
	if !ok {
		return nil, false
	}

	if value, ok := et.(*entry); ok {
		return value.value, true
	}
	return nil, false
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
	b.kvCount = 0
	if b.root == nil {
		return
	}
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
				b.kvCount = b.kvCount + len(nd.keys)
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

// 层序打印树结构
func (b *BPlusTree) Print() {
	b.CountNode()
	fmt.Println("----------------------------------------------------------------------------------------------------start print tree")
	fmt.Printf("Internal Size: min=%d max=%d, Leaf Size: min=%d max=%d, Total Node Size:%d, Total K/V Size:%d\n", b.internalMaxSize/2, b.internalMaxSize, b.leafMaxSize/2,
		b.leafMaxSize, b.GetNodeCount(), b.GetKeyCount())
	queue := make([]interface{}, 0)
	if b.root != nil {
		queue = append(queue, b.root)
	}

	level := 0
	for len(queue) != 0 {
		level++
		size := len(queue)
		str := ""
		for i := 0; i < size; i++ {
			nodeI := queue[i]
			if nodeI == nil {
				str = strings.Trim(str, " &&")
				str = str + " - "
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
					if strings.HasSuffix(str, " - ") {
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
		str = strings.Trim(str, "-")
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
		pointers: make([]interface{}, 0, b.leafMaxSize),
		maxSize:  b.leafMaxSize,
	}
}

func (b *BPlusTree) makeEmptyInternalNode() *node {
	return &node{
		isLeaf:   false,
		keys:     make([]key, 0, b.internalMaxSize),
		pointers: make([]interface{}, 0, b.internalMaxSize),
		maxSize:  b.internalMaxSize,
	}
}

func print(prt bool, a ...interface{}) {
	if prt {
		fmt.Println(a...)
	}
}

// 检查是否满足B+树(仅测试用，并且有点小问题，不想改了)
func (b *BPlusTree) check(prt bool) {
	defer func() {
		if err := recover(); err != nil {
			b.Print()
			panic(err)
		}
	}()
	if b.root == nil {
		return
	}
	if b.root.isLeaf {
		if len(b.root.keys) > b.leafMaxSize {
			panic("b plus tree wrong max size")
		}
		for i, ky := range b.root.keys {
			if i == 0 {
				// fmt.Println("key: ", ky)
				continue
			}
			if ky.compare(b.root.keys[i-1]) != 1 {
				panic(fmt.Sprintf("b plus tree error i=%d ky=%s", i, ky))
			}
			// fmt.Println("key: ", ky)
		}
		return
	}
	if len(b.root.keys) > b.internalMaxSize {
		panic("b plus tree wrong max size")
	}
	for i, ky := range b.root.keys {
		n := b.root.pointers[i].(*node)
		if i == 0 {
			b.checkTree(n, ky, prt)
			print(prt, "key: ", ky)
			continue
		}
		if ky.compare(b.root.keys[i-1]) == 1 {
			b.checkTree(n, ky, prt)
			print(prt, "key: ", ky)
		} else {
			panic(fmt.Sprintf("b plus tree error i=%d ky=%s", i, ky))
		}
	}
	b.checkTree(b.root.lastOrNextNode, "", prt)
}

func (b *BPlusTree) checkTree(nd *node, lastKey key, prt bool) {
	if !nd.isHalf() {
		panic(fmt.Sprintf("nd should be half but: %d", len(nd.keys)))
	}
	if len(nd.keys) > nd.maxSize || len(nd.keys) < nd.getHalf() {
		panic("b plus tree node wrong size")
	}
	if nd.isLeaf {
		for i, ky := range nd.keys {
			if lastKey.compare("") == 0 {
				if i == 0 {
					print(prt, "key: ", ky)
					continue
				}
				if ky.compare(nd.keys[i-1]) != 1 {
					panic("b plus tree error")
				}
				print(prt, "key: ", ky)
			} else {
				if i == 0 && lastKey.compare(ky) != 1 {
					panic("b plus tree error")
				}
				if i == 0 {
					print(prt, "key: ", ky)
					continue
				}
				if lastKey.compare(ky) != 1 || ky.compare(nd.keys[i-1]) != 1 {
					panic("b plus tree error")
				}
				print(prt, "key: ", ky)
			}

		}
		return
	}
	for i, ky := range nd.keys {
		n := nd.pointers[i].(*node)
		if i == 0 {
			b.checkTree(n, ky, prt)
			print(prt, "key: ", ky)
			continue
		}
		if lastKey.compare("") == 0 && ky.compare(nd.keys[i-1]) == 1 {
			b.checkTree(n, ky, prt)
			print(prt, "key: ", ky)
		} else {
			if lastKey.compare(ky) == 1 && ky.compare(nd.keys[i-1]) == 1 {
				b.checkTree(n, ky, prt)
				print(prt, "key: ", ky)
			} else {
				panic(fmt.Sprintf("b plus tree error lastKey:%s ky:%s i:%d", lastKey, ky, i))
			}
		}

	}
	b.checkTree(nd.lastOrNextNode, lastKey, prt)

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
	if len(leafNode.keys) < leafNode.maxSize {
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
		leafNode.keys = append(leafNode.keys, tempNode.keys[0:leafNode.getHalf()+1]...)
		leafNode.pointers = append(leafNode.pointers, tempNode.pointers[0:leafNode.getHalf()+1]...)
		siblingNode.keys = append(siblingNode.keys, tempNode.keys[leafNode.getHalf()+1:]...)
		siblingNode.pointers = append(siblingNode.pointers, tempNode.pointers[leafNode.getHalf()+1:]...)

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
	if len(parentNode.keys) < parentNode.maxSize {
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
		parentNode.keys = append(parentNode.keys, tempKeys[0:parentNode.getHalf()]...)
		for i := 0; i < parentNode.getHalf(); i++ {
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
		lst, ok := tempPointers[parentNode.getHalf()].(*node)
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

		siblingParentNode.keys = append(siblingParentNode.keys, tempKeys[parentNode.getHalf()+1:]...)
		for i := parentNode.getHalf() + 1; i < parentNode.maxSize+1; i++ {
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
		lst, ok = tempPointers[parentNode.maxSize+1].(*node)
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

		childKeyTwo := tempKeys[parentNode.getHalf()]
		b.insertIntoParent(parentNode, siblingParentNode, childKeyTwo)
	}

}

func (b *BPlusTree) delete(ky key) (interface{}, bool) {
	leafNode := b.findLeafNode(ky)
	ent, ok := leafNode.findRecord(ky)
	if !ok {
		return nil, false
	}
	b.deleteNode(leafNode, ky, ent)
	v, ok := ent.(*entry)
	if !ok {
		panic("should be entry")
	}
	return v.value, true
}
func (b *BPlusTree) deleteNode(nd *node, ky key, p interface{}) {
	nd.delete(ky, p)
	if nd.parent == nil && len(nd.keys) == 0 {
		b.root = nd.lastOrNextNode
		if b.root != nil {
			b.root.parent = nil
		} else {
			b.root = b.makeEmptyLeafNode()
		}
		return
	}
	if nd.parent == nil {
		return
	}
	if !nd.isHalf() {
		sibling, index, ky, isPrev := nd.lookupSibling()
		if len(sibling.keys)+len(nd.keys) < nd.maxSize {
			// Coalesce
			if !isPrev {
				sibling, nd = nd, sibling
			}
			if !sibling.isLeaf {
				sibling.keys = append(sibling.keys, ky)
				sibling.keys = append(sibling.keys, nd.keys...)
				sibling.pointers = append(sibling.pointers, sibling.lastOrNextNode)
				sibling.pointers = append(sibling.pointers, nd.pointers...)
				sibling.lastOrNextNode = nd.lastOrNextNode
				for _, p := range sibling.pointers {
					p.(*node).parent = sibling
				}
				if sibling.lastOrNextNode != nil {
					sibling.lastOrNextNode.parent = sibling
				}
			} else {
				sibling.keys = append(sibling.keys, nd.keys...)
				sibling.pointers = append(sibling.pointers, nd.pointers...)
				sibling.lastOrNextNode = nd.lastOrNextNode
			}
			b.deleteNode(sibling.parent, ky, nd)
		} else {
			// Redistribution
			if isPrev {
				var lastKey key
				if !nd.isLeaf {
					lastKey = sibling.keys[len(sibling.keys)-1]
					lastPointer := sibling.lastOrNextNode
					tempKeys := []key{ky}
					tempPointers := []interface{}{lastPointer}
					tempKeys = append(tempKeys, nd.keys...)
					tempPointers = append(tempPointers, nd.pointers...)
					nd.keys = tempKeys
					nd.pointers = tempPointers
					lastPointer.parent = nd
					sibling.keys = sibling.keys[0 : len(sibling.keys)-1]
					sibling.lastOrNextNode = sibling.pointers[len(sibling.pointers)-1].(*node)
					sibling.pointers = sibling.pointers[0 : len(sibling.pointers)-1]
					nd.parent.keys[index] = lastKey
				} else {
					lastKey = sibling.keys[len(sibling.keys)-1]
					lastPointer := sibling.pointers[len(sibling.pointers)-1]
					tempKeys := []key{lastKey}
					tempPointers := []interface{}{lastPointer}
					tempKeys = append(tempKeys, nd.keys...)
					tempPointers = append(tempPointers, nd.pointers...)
					nd.keys = tempKeys
					nd.pointers = tempPointers
					sibling.keys = sibling.keys[0 : len(sibling.keys)-1]
					sibling.pointers = sibling.pointers[0 : len(sibling.pointers)-1]
					nd.parent.keys[index] = lastKey
				}
			} else {
				var firstKey key
				if !nd.isLeaf {
					firstKey = sibling.keys[0]
					firstPointer := sibling.pointers[0]
					nd.keys = append(nd.keys, ky)
					nd.pointers = append(nd.pointers, nd.lastOrNextNode)
					nd.lastOrNextNode = firstPointer.(*node)
					firstPointer.(*node).parent = nd
					sibling.keys = sibling.keys[1:]
					sibling.pointers = sibling.pointers[1:]
					nd.parent.keys[index] = firstKey
				} else {
					firstKey = sibling.keys[0]
					firstPointer := sibling.pointers[0]
					nd.keys = append(nd.keys, firstKey)
					nd.pointers = append(nd.pointers, firstPointer)
					sibling.keys = sibling.keys[1:]
					sibling.pointers = sibling.pointers[1:]
					nd.parent.keys[index] = sibling.keys[0]
				}
			}
		}
	}

}
