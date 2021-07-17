# B+树理论与实现
#### 1、B+树出现的背景（大量数据->磁盘搜索）
* 应用数据极速变多，在内存中已经不能经济实惠的保存大量搜索的数据，所以考虑通过磁盘来帮助搜索，需要将全部或部分数据保存在磁盘。

#### 2、B+树特点

1. 平衡多叉树，从根节点到所有叶子节点的高度相同
1. 叶子节点之间通过指针互连，方便遍历叶子节点
#### 3、B+树结构
1. 节点：主要是搜索key和值value的集合，一般一个节点都保存在一个磁盘的block上，一次IO访问一个block
1. 内部节点（根节点）
    内部节点的元素为key和关联的子节点指针
    根节点是特殊的内部节点
1. 叶子节点
    叶子节点元素为key和值或指向值的指针
    叶子节点之间通过指针互连
#### 4、节点大小（边界条件）
1. fanout：节点可包含的子节点数量
如何计算fanout？
一半3层高度的B+树即可保存千万级别的数据，如何计算一颗B+树能保存的数据量？
#### 5、查询
#### 6、更新
#### 7、删除
#### 8、增删改的时间复杂度、空间复杂度、填充因子、fanout
#### 9、其它：
1. 聚簇索引、非聚簇索引
1. 前缀压缩
1. 范围查询
1. 重复key
1. 批量加载
1. 并发读写
1. 合并、分裂策略

#### 10. mysql选择B+树而不是B树原因
1. B+树内部节点不保存具体数据，只保存在叶子节点
1. B+树查询在范围查询效率更高，B+树在范围查询可以根据叶子节点的链接直接顺序遍历，B树需要遍历完子树才能完成范围查找，访问的磁盘IO次数更多
1. B+树查询效率更稳定，B+树每次访问数据都会到达叶子节点，查询时间稳定，而B树访问的数据不确定在第几层所以查询效率不太稳定
1. B+树访问磁盘数更少，由于结构原因，B+树更矮胖访问的磁盘树更少，B树更高访问的磁盘树更多