package main

import (
	"fmt"
	"testing"
)

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
	v, ok = leafNode.findRecord(Key("trump"))
	if ok {
		t.Fatalf("should not find trump value")
	}
}

func TestBPlusTree_insertIntoLeaf(t *testing.T) {
	tree := StartNewTree(6, 6)
	leafNode := makeLeafNode([]string{"key1", "key2", "key4", "key5"},
		[]string{"key1_value", "key2_value", "key4_value", "key5_value"})
	targetKey := Key("key3")
	record := &Record{
		value: []byte("key3_value"),
	}
	tree.insertIntoLeaf(leafNode, targetKey, record)
	if targetKey.Compare(leafNode.keys[2]) != 0 {
		t.Fatalf("should be key3")
	}
	r, ok := leafNode.pointers[2].(*Record)
	if !ok {
		t.Fatalf("should be record")
	}
	if string(r.value) != "key3_value" {
		t.Fatalf("should be key3_value")
	}
	t.Logf("keys: %v", leafNode.keys)
	values := make([]string, 0)
	for _, i := range leafNode.pointers {
		r, _ := i.(*Record)
		values = append(values, string(r.value))
	}
	t.Logf("record: %v", values)
	targetKey = Key("key0")
	record = &Record{
		value: []byte("key0_value"),
	}
	tree.insertIntoLeaf(leafNode, targetKey, record)
	if targetKey.Compare(leafNode.keys[0]) != 0 {
		t.Fatalf("should be key0")
	}
	t.Logf("keys: %v", leafNode.keys)
	targetKey = Key("key6")
	record = &Record{
		value: []byte("key6_value"),
	}
	tree.insertIntoLeaf(leafNode, targetKey, record)
	if targetKey.Compare(leafNode.keys[6]) != 0 {
		t.Fatalf("should be key6")
	}
	r, ok = leafNode.pointers[6].(*Record)
	if !ok {
		t.Fatalf("should be record")
	}
	if string(r.value) != "key6_value" {
		t.Fatalf("should be key6_value")
	}
	t.Logf("keys: %v", leafNode.keys)
}

func TestBPlusTree_updateRecord(t *testing.T) {
	leafNode := makeLeafNode([]string{"key1", "key2", "key4", "key5"},
		[]string{"key1_value", "key2_value", "key4_value", "key5_value"})
	targetKey := Key("key2")
	record := &Record{
		value: []byte("key2222_value"),
	}
	leafNode.updateRecord(targetKey, record)
	r, ok := leafNode.pointers[1].(*Record)
	if !ok {
		t.Fatalf("should be record")
	}
	if string(r.value) != "key2222_value" {
		t.Fatalf("should be key2222_value, actully value: %v", string(r.value))
	}
	t.Logf("keys: %v", leafNode.keys)
}

func TestInsert(t *testing.T) {
	tree := StartNewTree(4, 4)
	tree.Insert("key1", "value1")
	tree.Insert("key2", "value2")
	tree.Insert("key3", "value3")
	tree.Insert("key4", "value4")
	fmt.Println("---------44444")
	tree.Print()
	tree.Insert("key5", "value5")
	fmt.Println("---------55555")
	tree.Print()
	tree.Insert("key6", "value6")
	fmt.Println("---------66666")
	tree.Print()
	// tree.Insert("key8", "value8")
	// tree.Insert("key9", "value9")
	// tree.Insert("key10", "value10")
	// tree.Insert("key11", "value11")
	// tree.Insert("key12", "value12")
	// tree.Insert("key13", "value13")
	// tree.Insert("key14", "value14")
	v, _ := tree.Find("key1")
	if v != "value1" {
		t.Fatalf("value should be value1, but value:%s", v)
	}
	v, _ = tree.Find("key2")
	if v != "value2" {
		t.Fatalf("value should be value2, , but value:%s", v)
	}
	v, ok := tree.Find("key0")
	if ok {
		t.Fatalf("value should not exsit")
	}
	if v != "" {
		t.Fatalf("value should be empty")
	}

}
