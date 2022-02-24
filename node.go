package btree

type node struct {
	// 是否是叶子节点
	isLeaf bool
	keys   []key
	// 非叶子节点是*Node，叶子节点是*entry, 最后一个指针挪到了lastOrNextNode
	// 所以len(keys)=len(pointers)
	pointers []interface{}
	parent   *node
	// 最后一个指针
	lastOrNextNode *node
}

type entry struct {
	value []byte
}

func (r *entry) toValue() string {
	return string(r.value)
}

func (n *node) findRecord(targetKey key) ([]byte, bool) {
	if !n.isLeaf {
		panic("should be leaf nd")
	}
	if len(n.keys) == 0 {
		return []byte{}, false
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
	return []byte{}, false
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
