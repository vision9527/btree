package btree

import "fmt"

type node struct {
	// 是否是叶子节点
	isLeaf bool
	keys   []key
	// 非叶子节点是*Node，叶子节点是*entry, 最后一个指针挪到了lastOrNextNode
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
func (n *node) delete(targetKey key) {
	if len(n.keys) == 0 {
		return
	}
	// 需要删除key的索引,
	// 删除相应的pointer索引：叶子节点index，内部节点index+1（内部节点只有在合并的时候才会有删除的情况）
	var index int
	for i, ky := range n.keys {
		if ky.compare(targetKey) == 0 {
			nd := n.pointers[i]
			if n.isLeaf {
				if _, ok := nd.(*entry); ok {
					index = i
					break
				}
				panic("should be entry")
			} else {
				if _, ok := nd.(*node); ok {
					index = i
					break
				}
				panic("should be node")
			}

		}
	}
	keys := n.keys[0:index]
	if index+1 != len(n.keys) {
		keys = append(keys, n.keys[index+1:]...)
	}
	n.keys = keys
	if n.isLeaf {
		pointers := n.pointers[0:index]
		if index+1 != len(n.pointers) {
			pointers = append(pointers, n.pointers[index+1:]...)
		}
		n.pointers = pointers
	} else {
		if index+1 == len(n.keys) {
			n.lastOrNextNode = n.pointers[index].(*node)
			n.pointers = n.pointers[0:index]
		} else if index+1 == len(n.keys)-1 {
			pointers := n.pointers[0 : index+1]
			if index+1 != len(n.keys) {
				n.pointers = pointers
			}
			n.pointers = pointers
		} else {
			pointers := n.pointers[0 : index+1]
			pointers = append(pointers, n.pointers[index+2:]...)
			n.pointers = pointers
		}
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
			if value, ok := et.(*entry); ok {
				return value.value, true
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
			index = index + 1
			sibling = nd.parent.pointers[index].(*node)
			isPrev = false
			ky = nd.parent.keys[0]
			return
		} else {
			// 默认用前一个
			ky = nd.parent.keys[index]
			index = index - 1
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
	return len(n.keys) >= n.maxSize/2
}

func (n *node) getHalf() int {
	return n.maxSize / 2
}
