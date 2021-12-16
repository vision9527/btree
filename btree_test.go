package main

import "testing"

func TestFind_RootIsLeaf(t *testing.T) {
	tree := NewTree()
	key1 := "key3"
	value1 := "value3"
	leafNode := &LeafNode{
		Pairs: [SizeLeafNode]*Element{
			{Key: key1, Value: value1},
			{Key: "key2", Value: "value2"},
		},
		Size: 2,
	}
	tree.Root = leafNode
	value, flag := tree.Find(key1)
	if !flag {
		t.Fatalf("cant find %s", key1)
		return
	}
	t.Logf("find key:%s , value:%s", key1, value)
	key100 := "key100"
	value, flag = tree.Find(key100)
	if flag {
		t.Fatalf("find key:%s, value:%s, but dont want", key100, value)
		return
	}
	t.Logf("success not find key:%s", key100)
}

func TestFind_RootIsInternalNode(t *testing.T) {
	// TODO
}

func TestInsert(t *testing.T) {
	tree := new(Tree)
	tree.Insert("key1", "value1")
	tree.Insert("key2", "value2")
	tree.Insert("key3", "value3")
	tree.Insert("key4", "value4")
	tree.Insert("key5", "value5")
	value, flag := tree.Find("key1")
	if !flag {
		t.Fatalf("find key:%s, value:%s", "key1", value)
		return
	}

	value, flag = tree.Find("key12")
	if flag {
		t.Fatalf("find key:%s, value:%s, but dont want", "key12", value)
		return
	}
	tree.Print()
}

func TestDelete(t *testing.T) {
	// TODO
}
