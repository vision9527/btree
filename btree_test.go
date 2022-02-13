package main

import "testing"

func makeLeafNode(keys []string, values []string) *Node {
	n := &Node{
		isLeaf: true,
	}
	for _, v := range keys {
		k := Key(v)
		n.keys = append(n.keys, k)
	}
	for _, v := range values {
		record := &Record{
			value: []byte(v),
		}
		n.pointers = append(n.pointers, record)
	}
	return n
}

func TestBPlusTreeFind_rootDontHaveChild(t *testing.T) {
	tree := StartNewTree(4, 4)
	root := &Node{
		isLeaf:   true,
		keys:     []Key{Key("test_key1"), Key("test_key2"), Key("test_key3")},
		pointers: make([]interface{}, 0),
	}
	tree.root = root
	p1 := &Record{
		value: []byte("test_value1"),
	}
	root.pointers = append(root.pointers, p1)
	p2 := &Record{
		value: []byte("test_value2"),
	}
	root.pointers = append(root.pointers, p2)
	p3 := &Record{
		value: []byte("test_value3"),
	}
	root.pointers = append(root.pointers, p3)
	v, ok := tree.Find("test_key1")
	if !ok {
		t.Fatalf("not find test_key1 value")
	}
	if v != "test_value1" {
		t.Fatalf("find test_key1 value, but not correct")
	}
	v, ok = tree.Find("test_key2")
	if !ok {
		t.Fatalf("not find test_key2 value")
	}
	if v != "test_value2" {
		t.Fatalf("find test_key2 value, but not correct")
	}
}

func TestBPlusTreeFind_rootHaveChild(t *testing.T) {
	tree := StartNewTree(6, 6)
	root := &Node{
		isLeaf:   false,
		keys:     []Key{Key("El Said"), Key("Mozart")},
		pointers: make([]interface{}, 2),
	}
	tree.root = root
	nodeOne := makeLeafNode([]string{"Brandt", "Califieri", "Crick", "Einstein"},
		[]string{"Brandt_value", "Califieri_value", "Crick_value", "Einstein_value"})
	nodeTwo := makeLeafNode([]string{"El Said", "Gold", "Katz", "Kim"},
		[]string{"El Said_value", "Gold_value", "Katz_value", "Kim_value"})
	nodeThree := makeLeafNode([]string{"Mozart", "Singh", "Srinivasan", "Wu"},
		[]string{"Mozart_value", "Singh_value", "Srinivasan_value", "Wu_value"})
	tree.root.pointers[0] = nodeOne
	tree.root.pointers[1] = nodeTwo
	tree.root.lastOrNextNode = nodeThree
	nodeOne.parent = tree.root
	nodeTwo.parent = tree.root
	nodeThree.parent = tree.root
	nodeOne.lastOrNextNode = nodeTwo
	nodeTwo.lastOrNextNode = nodeThree
	v, ok := tree.Find("Mozart")
	if !ok {
		t.Fatalf("not find Mozart value")
	}
	if v != "Mozart_value" {
		t.Fatalf("find Mozart_value, but not correct")
	}
	v, ok = tree.Find("Wu")
	if !ok {
		t.Fatalf("not find Wu value")
	}
	if v != "Wu_value" {
		t.Fatalf("find Wu_value, but not correct")
	}

}

func TestBPlusTreeFind_findRecord(t *testing.T) {
	leafNode := makeLeafNode([]string{"Mozart", "Singh", "Srinivasan", "Wu"},
		[]string{"Mozart_value", "Singh_value", "Srinivasan_value", "Wu_value"})
	v, ok := leafNode.findRecord(Key("Singh"))
	if !ok {
		t.Fatalf("should find Singh_value")
	}
	if v != "Singh_value" {
		t.Fatalf("find Singh_value, but not correct")
	}
}
