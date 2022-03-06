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
	sibling, _, ky, isPrev := leafNode1.lookupSibling()
	if sibling != leafNode2 {
		t.Fatal("should be leafNode2")
	}
	if ky.toString() != "c" {
		t.Fatal("should be c")
	}
	if isPrev {
		t.Fatal("should be not prev")
	}

	sibling, _, ky, isPrev = leafNode2.lookupSibling()
	if sibling != leafNode1 {
		t.Fatal("should be leafNode1")
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
		keys:     make([]key, 0, 3),
		pointers: make([]interface{}, 0, 3),
		maxSize:  3,
	}
	parentNode.keys = append(parentNode.keys, "c", "e")
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
	leafNode1.lastOrNextNode = leafNode2
	leafNode2.lastOrNextNode = leafNode3
	leafNode1.parent = parentNode
	leafNode2.parent = parentNode
	leafNode3.parent = parentNode
	parentNode.pointers = append(parentNode.pointers, leafNode1)
	parentNode.pointers = append(parentNode.pointers, leafNode2)
	parentNode.lastOrNextNode = leafNode3
	sibling, _, ky, isPrev := leafNode1.lookupSibling()
	if sibling != leafNode2 {
		t.Fatal("should be leafNode2")
	}
	if ky.toString() != "c" {
		t.Fatal("should be c")
	}
	if isPrev {
		t.Fatal("should be not prev")
	}

	sibling, _, ky, isPrev = leafNode2.lookupSibling()
	if sibling != leafNode3 {
		t.Fatal("should be leafNode3")
	}
	if ky.toString() != "e" {
		t.Fatal("should be e")
	}
	if isPrev {
		t.Fatal("should be prev")
	}

	sibling, _, ky, isPrev = leafNode3.lookupSibling()
	if sibling != leafNode2 {
		t.Fatal("should be leafNode2")
	}
	if ky.toString() != "e" {
		t.Fatal("should be e")
	}
	if !isPrev {
		t.Fatal("should be prev")
	}
}
