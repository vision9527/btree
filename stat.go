package btree

// 查询统计（测试使用）
type stat struct {
	// 查询遍历到的节点数
	count int64
	// 树的节点总数
	nodeCount int64
	// 数的高度
	level int64
}

func (b *stat) incrCount() {
	b.count++
}

func (b *stat) resetCount() {
	b.count = 0
}

func (b *stat) GetCount() int64 {
	return b.count
}

func (b *stat) GetLevel() int64 {
	return b.level
}

func (b *stat) GetNodeCount() int64 {
	return b.nodeCount
}
