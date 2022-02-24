package btree

import (
	"fmt"
	"testing"
)

func TestBPlusTreeFind_findRecord(t *testing.T) {
	leafNode := makeTestLeafNode([]string{"Mozart", "Singh", "Srinivasan", "Wu"},
		[]string{"Mozart_value", "Singh_value", "Srinivasan_value", "Wu_value"})
	v, ok := leafNode.findRecord(key("Singh"))
	if !ok {
		t.Fatalf("should find Singh_value")
	}
	if fmt.Sprintf("%v", v) != "Singh_value" {
		t.Fatalf("find Singh_value, but not correct")
	}
	_, ok = leafNode.findRecord(key("trump"))
	if ok {
		t.Fatalf("should not find trump value")
	}
}

func TestBPlusTree_updateRecord(t *testing.T) {
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
