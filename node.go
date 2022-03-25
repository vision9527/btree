package btree

import "fmt"

type node struct {
	// 是否是叶子节点
	isLeaf bool
	keys   []key
	// 非叶子节点是*node，叶子节点是*entry, 最后一个指针挪到了lastOrNextNode
	// 所以len(keys)=len(pointers)
	pointers []interface{}
	parent   *node
	maxSize  int
	// 最后一个指针
	lastOrNextNode *node
}

type entry struct {
	value interface{}
}

func (r *entry) toValue() string {
	return fmt.Sprintf("%v", r.value)
}

// 删除key的同时删除pointer
func (n *node) delete(targetKey key, p interface{}) {
	if len(n.keys) == 0 {
		return
	}
	index := -1
	for i, ky := range n.keys {
		if ky.compare(targetKey) == 0 {
			index = i
			break
		}
	}
	tempKeys := n.keys[0:index]
	if index < len(n.keys)-1 {
		tempKeys = append(tempKeys, n.keys[index+1:]...)
	}
	n.keys = tempKeys
	index = -1
	for i, pt := range n.pointers {
		if pt == p {
			index = i
			break
		}
	}
	if index != -1 {
		tempPointers := n.pointers[0:index]
		if index < len(n.pointers)-1 {
			tempPointers = append(tempPointers, n.pointers[index+1:]...)
		}
		n.pointers = tempPointers
	} else {
		n.lastOrNextNode = n.pointers[len(n.pointers)-1].(*node)
		n.pointers = n.pointers[0 : len(n.pointers)-1]
	}
}

func (n *node) findRecord(targetKey key) (interface{}, bool) {
	if !n.isLeaf {
		panic("should be leaf nd")
	}
	if len(n.keys) == 0 {
		return nil, false
	}
	// 可使用二分查找，待优化
	for i, ky := range n.keys {
		if ky.compare(targetKey) == 0 {
			et := n.pointers[i]
			if _, ok := et.(*entry); ok {
				return et, true
			}
			panic("should be entry")
		}
	}
	return nil, false
}

func (n *node) updateRecord(targetKey key, et *entry) bool {
	// 如果值已经存在则更新
	for i, k := range n.keys {
		if targetKey.compare(k) == 0 {
			r, ok := n.pointers[i].(*entry)
			if !ok {
				panic("should be *entry")
			}
			r.value = et.value
			n.pointers[i] = r
			return true
		}
	}
	return false
}

func (n *node) insertNextAfterPrev(childKey key, prev, next *node) {
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
		n.keys = append(n.keys[:number], append([]key{childKey}, n.keys[number:]...)...)
		n.pointers = append(n.pointers[:number+1], append([]interface{}{next}, n.pointers[number+1:]...)...)
	}

}

// 找到nd的兄弟节点
func (nd *node) lookupSibling() (sibling *node, index int, ky key, isPrev bool) {
	if nd.parent != nil {
		index = -1
		for i, pointer := range nd.parent.pointers {
			n, ok := pointer.(*node)
			if !ok {
				panic("should be node")
			}
			if n == nd {
				index = i
				break
			}
		}
		if index == -1 {
			index = len(nd.parent.pointers) - 1
			sibling = nd.parent.pointers[index].(*node)
			isPrev = true
			ky = nd.parent.keys[index]
			return
		}
		// pointers里最后一个
		if index == len(nd.parent.pointers)-1 {
			sibling = nd.parent.lastOrNextNode
			isPrev = false
			ky = nd.parent.keys[index]
			return
		} else if index == 0 {
			// pointers里第一个
			// index = index + 1
			sibling = nd.parent.pointers[index+1].(*node)
			isPrev = false
			ky = nd.parent.keys[0]
			return
		} else {
			// 默认用前一个
			index = index - 1
			ky = nd.parent.keys[index]
			sibling = nd.parent.pointers[index].(*node)
			isPrev = true
			return
		}
	}
	return
}

func (n *node) copy() *node {
	nd := &node{
		isLeaf:         n.isLeaf,
		keys:           make([]key, 0),
		pointers:       make([]interface{}, 0),
		parent:         n.parent,
		lastOrNextNode: n.lastOrNextNode,
	}
	nd.keys = append(nd.keys, n.keys...)
	nd.pointers = append(nd.pointers, n.pointers...)
	return nd
}

// 是否半满
func (n *node) isHalf() bool {
	return len(n.keys) >= n.getHalf()
}

func (n *node) getHalf() int {
	return n.maxSize / 2
}

// 从此节点检查key顺序
func (n *node) checkOrder() {
	if !n.isLeaf {
		panic("only leaf node for use")
	}
	currentNode := n
	lastKey := key("")
	for currentNode != nil {
		for _, k := range currentNode.keys {
			if k.compare(lastKey) < 1 {
				panic("wrong key order")
			}
			lastKey = k
		}
		currentNode = currentNode.lastOrNextNode
	}
}
