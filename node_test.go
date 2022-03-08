package btree

import (
	"testing"
)

func TestNode_findRecord(t *testing.T) {
	leafNode := makeTestLeafNode([]string{"Mozart", "Singh", "Srinivasan", "Wu"},
		[]string{"Mozart_value", "Singh_value", "Srinivasan_value", "Wu_value"})
	v, ok := leafNode.findRecord(key("Singh"))
	if !ok {
		t.Fatalf("should find Singh_value")
	}
	value, ok := v.(*entry)
	if !ok {
		t.Fatalf("shuold be *entry")
	}
	if value.toValue() != "Singh_value" {
		t.Fatalf("find Singh_value, but not correct")
	}
	_, ok = leafNode.findRecord(key("trump"))
	if ok {
		t.Fatalf("shuold be *entry")
	}
}

func TestNode_updateRecord(t *testing.T) {
	leafNode := makeTestLeafNode([]string{"key1", "key2", "key4", "key5"},
		[]string{"key1_value", "key2_value", "key4_value", "key5_value"})
	targetKey := key("key2")
	et := &entry{
		value: "key2222_value",
	}
	leafNode.updateRecord(targetKey, et)
	r, ok := leafNode.pointers[1].(*entry)
	if !ok {
		t.Fatalf("should be entry")
	}
	if r.toValue() != "key2222_value" {
		t.Fatalf("should be key2222_value, actully value: %v", r.toValue())
	}
	t.Logf("keys: %v", leafNode.keys)
}

func TestNode_lookupSiblingOne(t *testing.T) {
	parentNode := &node{
		isLeaf:   false,
		keys:     make([]key, 0, 3),
		pointers: make([]interface{}, 0, 3),
		maxSize:  3,
	}
	parentNode.keys = append(parentNode.keys, "c")
	leafNode1 := &node{
		isLeaf:   true,
		keys:     []key{"a", "b"},
		pointers: []interface{}{entry{}, entry{}},
	}
	leafNode2 := &node{
		isLeaf:   true,
		keys:     []key{"c", "d"},
		pointers: []interface{}{entry{}, entry{}},
	}
	leafNode1.lastOrNextNode = leafNode2
	leafNode1.parent = parentNode
	leafNode2.parent = parentNode
	parentNode.pointers = append(parentNode.pointers, leafNode1)
	parentNode.lastOrNextNode = leafNode2
	sibling, index, ky, isPrev := leafNode1.lookupSibling()
	if sibling != leafNode2 {
		t.Fatal("should be leafNode2")
	}
	if index != 0 {
		t.Fatal("should be 0")
	}
	if ky.toString() != "c" {
		t.Fatal("should be c")
	}
	if isPrev {
		t.Fatal("should be not prev")
	}

	sibling, index, ky, isPrev = leafNode2.lookupSibling()
	if sibling != leafNode1 {
		t.Fatal("should be leafNode1")
	}
	if index != 0 {
		t.Fatal("should be 0")
	}
	if ky.toString() != "c" {
		t.Fatal("should be c")
	}
	if !isPrev {
		t.Fatal("should be prev")
	}
}

func TestNode_lookupSiblingTwo(t *testing.T) {
	parentNode := &node{
		isLeaf:   false,
		keys:     make([]key, 0, 4),
		pointers: make([]interface{}, 0, 4),
		maxSize:  4,
	}
	parentNode.keys = append(parentNode.keys, "c", "e", "d")
	leafNode1 := &node{
		isLeaf:   true,
		keys:     []key{"a", "b"},
		pointers: []interface{}{entry{}, entry{}},
	}
	leafNode2 := &node{
		isLeaf:   true,
		keys:     []key{"c", "d"},
		pointers: []interface{}{entry{}, entry{}},
	}
	leafNode3 := &node{
		isLeaf:   true,
		keys:     []key{"e", "f", "g"},
		pointers: []interface{}{entry{}, entry{}, entry{}},
	}
	leafNode4 := &node{
		isLeaf:   true,
		keys:     []key{"e", "f", "g"},
		pointers: []interface{}{entry{}, entry{}, entry{}},
	}
	leafNode1.lastOrNextNode = leafNode2
	leafNode2.lastOrNextNode = leafNode3
	leafNode3.lastOrNextNode = leafNode4
	leafNode1.parent = parentNode
	leafNode2.parent = parentNode
	leafNode3.parent = parentNode
	leafNode4.parent = parentNode
	parentNode.pointers = append(parentNode.pointers, leafNode1)
	parentNode.pointers = append(parentNode.pointers, leafNode2)
	parentNode.pointers = append(parentNode.pointers, leafNode3)
	parentNode.lastOrNextNode = leafNode4
	sibling, index, ky, isPrev := leafNode1.lookupSibling()
	if sibling != leafNode2 {
		t.Fatal("should be leafNode2")
	}
	if index != 0 {
		t.Fatal("should be 0")
	}
	if ky.toString() != "c" {
		t.Fatal("should be c")
	}
	if isPrev {
		t.Fatal("should be not prev")
	}

	sibling, index, ky, isPrev = leafNode2.lookupSibling()
	if sibling != leafNode1 {
		t.Fatal("should be leafNode1")
	}
	if index != 0 {
		t.Fatal("should be 0")
	}
	if ky.toString() != "c" {
		t.Fatal("should be e")
	}
	if !isPrev {
		t.Fatal("should be prev")
	}

	sibling, index, ky, isPrev = leafNode3.lookupSibling()
	if sibling != leafNode4 {
		t.Fatal("should be leafNode4")
	}
	if index != 2 {
		t.Fatal("should be 2")
	}
	if ky.toString() != "d" {
		t.Fatal("should be e")
	}
	if isPrev {
		t.Fatal("should be prev")
	}

	sibling, index, ky, isPrev = leafNode4.lookupSibling()
	if sibling != leafNode3 {
		t.Fatal("should be leafNode3")
	}
	if index != 2 {
		t.Fatal("should be 2")
	}
	if ky.toString() != "d" {
		t.Fatal("should be d")
	}
	if !isPrev {
		t.Fatal("should be prev")
	}

}

func TestNode_deleteLeaf(t *testing.T) {
	keys := []key{key("a"), key("b"), key("c"), key("d")}
	next := new(node)
	a := &entry{}
	pointers := []interface{}{a, &entry{}, &entry{}, &entry{}}
	nd := &node{
		isLeaf:         true,
		keys:           keys,
		pointers:       pointers,
		maxSize:        4,
		lastOrNextNode: next,
	}
	nd.delete(key("a"), a)
	for i, k := range nd.keys {
		if i == 0 && k.toString() == "b" {
			continue
		}
		if i == 1 && k.toString() == "c" {
			continue
		}
		if i == 2 && k.toString() == "d" {
			continue
		}
		t.Fatalf("node delete wrong, node keys:%v\n", nd.keys)
	}
	if len(nd.pointers) != 3 {
		t.Fatalf("node delete wrong, node pointers:%v\n", nd.pointers)
	}

}

func TestNode_deleteInternal(t *testing.T) {
	keys := []key{key("a"), key("b"), key("c"), key("d")}
	next := new(node)
	a := &node{}
	b := &node{}
	c := &node{}
	d := &node{}
	pointers := []interface{}{a, b, c, d}
	nd := &node{
		isLeaf:         false,
		keys:           keys,
		pointers:       pointers,
		maxSize:        4,
		lastOrNextNode: next,
	}
	nd.delete(key("d"), next)
	for i, k := range nd.keys {
		if i == 0 && k.toString() == "a" {
			continue
		}
		if i == 1 && k.toString() == "b" {
			continue
		}
		if i == 2 && k.toString() == "c" {
			continue
		}
		t.Fatalf("node delete wrong, node keys:%v\n", nd.keys)
	}
	for i, k := range nd.pointers {
		if i == 0 && k == a {
			continue
		}
		if i == 1 && k == b {
			continue
		}
		if i == 2 && k == c {
			continue
		}
		if i == 3 && k == d {
			continue
		}
		t.Fatalf("node delete wrong, node pointer:%v\n", nd.pointers)
	}
	if nd.lastOrNextNode != d {
		t.Fatal("node delete wrong lastOrNextNode")
	}

}
