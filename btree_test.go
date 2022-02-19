package main

import (
	"math/rand"
	"testing"
)

func makeTestLeafNode(keys []string, values []string) *Node {
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

func GenTestKeyAndValue(repeatNum int) []string {
	str := "abcdefghijklmnopqrstuvwxyz"
	result := make([]string, 0)
	for _, i := range str {
		s := ""
		for j := 0; j < repeatNum; j++ {
			s = s + string(i)
			result = append(result, s)
		}

	}
	return result
}

func ShuffleTestkv(kv []string) {
	for len(kv) > 0 {
		size := len(kv)
		r := rand.Intn(size)
		kv[r], kv[size-1] = kv[size-1], kv[r]
		kv = kv[0 : size-1]
	}
}

func TestGenTestKeyAndValue_Print(t *testing.T) {
	testkv := GenTestKeyAndValue(3)
	t.Log("test_keys:", testkv)
}

func TestBPlusTreeFind_rootDontHaveChild(t *testing.T) {
	tree, _ := StartNewTree(4, 4)
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
	tree, _ := StartNewTree(6, 6)
	root := &Node{
		isLeaf:   false,
		keys:     []Key{Key("El Said"), Key("Mozart")},
		pointers: make([]interface{}, 2),
	}
	tree.root = root
	nodeOne := makeTestLeafNode([]string{"Brandt", "Califieri", "Crick", "Einstein"},
		[]string{"Brandt_value", "Califieri_value", "Crick_value", "Einstein_value"})
	nodeTwo := makeTestLeafNode([]string{"El Said", "Gold", "Katz", "Kim"},
		[]string{"El Said_value", "Gold_value", "Katz_value", "Kim_value"})
	nodeThree := makeTestLeafNode([]string{"Mozart", "Singh", "Srinivasan", "Wu"},
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
	leafNode := makeTestLeafNode([]string{"Mozart", "Singh", "Srinivasan", "Wu"},
		[]string{"Mozart_value", "Singh_value", "Srinivasan_value", "Wu_value"})
	v, ok := leafNode.findRecord(Key("Singh"))
	if !ok {
		t.Fatalf("should find Singh_value")
	}
	if v != "Singh_value" {
		t.Fatalf("find Singh_value, but not correct")
	}
	_, ok = leafNode.findRecord(Key("trump"))
	if ok {
		t.Fatalf("should not find trump value")
	}
}

func TestBPlusTree_insertIntoLeaf(t *testing.T) {
	tree, _ := StartNewTree(6, 6)
	leafNode := makeTestLeafNode([]string{"key1", "key2", "key4", "key5"},
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
	leafNode := makeTestLeafNode([]string{"key1", "key2", "key4", "key5"},
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

func TestInsertCaseOne(t *testing.T) {
	tree, _ := StartNewTree(3, 3)
	tree.Insert("a", "a")
	tree.Insert("b", "b")
	tree.Insert("c", "c")
	tree.Insert("d", "d")
	tree.Print()
	tree.Insert("e", "e")
	tree.Print()
	tree.Insert("f", "f")
	tree.Print()
	tree.Insert("g", "g")
	tree.Print()
	tree.Insert("h", "h")
	tree.Print()
	tree.Insert("i", "i")
	tree.Print()
	tree.Insert("j", "j")
	tree.Print()
	tree.Insert("k", "k")
	tree.Print()
	tree.Insert("l", "l")
	tree.Print()
	tree.Insert("m", "m")
	tree.Print()
	tree.Insert("n", "n")
	tree.Print()
	v, _ := tree.Find("a")
	if v != "a" {
		t.Fatalf("value should be a, but value:%s", v)
	}
	v, _ = tree.Find("g")
	if v != "g" {
		t.Fatalf("value should be g, but value:%s", v)
	}
	v, _ = tree.Find("i")
	if v != "i" {
		t.Fatalf("value should be i, but value:%s", v)
	}
	v, _ = tree.Find("m")
	if v != "m" {
		t.Fatalf("value should be m, but value:%s", v)
	}
	v, _ = tree.Find("n")
	if v != "n" {
		t.Fatalf("value should be n, but value:%s", v)
	}
	v, ok := tree.Find("xxxxxx")
	if ok {
		t.Fatalf("value should not exsit")
	}
	if v != "" {
		t.Fatalf("value should be empty")
	}

}

// 小fanout， 多key，value
func TestInsertCaseTwo(t *testing.T) {
	tree, _ := StartNewTree(3, 3)
	testkv := GenTestKeyAndValue(3)
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
	}
	tree.Print()
	for i := 0; i < len(testkv); i++ {
		key := testkv[i]
		value := testkv[i]
		v, ok := tree.Find(key)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if v != value {
			t.Fatalf("value should be %s, but value:%s", key, v)
		}
	}

}

// 小fanout，乱序插入
func TestInsertCaseShuffleTestkv(t *testing.T) {
	tree, _ := StartNewTree(4, 4)
	testkv := GenTestKeyAndValue(1)
	ShuffleTestkv(testkv)
	t.Log(testkv)
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
	}
	tree.Print()

	for i := 0; i < len(testkv); i++ {
		key := testkv[i]
		value := testkv[i]
		v, ok := tree.Find(key)
		if !ok {
			t.Fatalf("value:%s, should exsit", key)
		}
		if v != value {
			t.Fatalf("value should be %s, but value:%s", key, v)
		}
	}
}

// 大量，乱序插入
func TestInsertCaseShuffleTestkv2(t *testing.T) {
	num := 100
	for n := 1; n < num; n++ {
		tree, _ := StartNewTree(4, 4)
		testkv := GenTestKeyAndValue(n)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			// fmt.Printf("insert: %s ...\n", testkv[i])
			tree.Insert(testkv[i], testkv[i])
			// tree.Print()
		}
		for i := 0; i < len(testkv); i++ {
			key := testkv[i]
			value := testkv[i]
			v, ok := tree.Find(key)
			if !ok {
				t.Fatalf("value:%s, should exsit", key)
			}
			if v != value {
				t.Fatalf("value should be %s, but value:%s", key, v)
			}
		}
	}

	num = 100
	for n := 1; n < num; n++ {
		tree, _ := StartNewTree(num, num)
		testkv := GenTestKeyAndValue(n)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			// fmt.Printf("insert: %s ...\n", testkv[i])
			tree.Insert(testkv[i], testkv[i])
			// tree.Print()
		}
		for i := 0; i < len(testkv); i++ {
			key := testkv[i]
			value := testkv[i]
			v, ok := tree.Find(key)
			if !ok {
				t.Fatalf("value:%s, should exsit", key)
			}
			if v != value {
				t.Fatalf("value should be %s, but value:%s", key, v)
			}
		}
	}
	num = 100
	for n := 1; n < num; n++ {
		tree, _ := StartNewTree(num, num)
		testkv := GenTestKeyAndValue(n)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			// fmt.Printf("insert: %s ...\n", testkv[i])
			tree.Insert(testkv[i], testkv[i])
			// tree.Print()
		}
		for i := 0; i < len(testkv); i++ {
			key := testkv[i]
			value := testkv[i]
			v, ok := tree.Find(key)
			if !ok {
				t.Fatalf("value:%s, should exsit", key)
			}
			if v != value {
				t.Fatalf("value should be %s, but value:%s", key, v)
			}
		}
	}
}

func TestInsertCaseShuffleTestkv3(t *testing.T) {
	num := 100
	for n := 4; n < num; n++ {
		tree, _ := StartNewTree(n+10, n)
		testkv := GenTestKeyAndValue(n)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			// fmt.Printf("insert: %s ...\n", testkv[i])
			tree.Insert(testkv[i], testkv[i])
			// tree.Print()
		}
		for i := 0; i < len(testkv); i++ {
			key := testkv[i]
			value := testkv[i]
			v, ok := tree.Find(key)
			if !ok {
				t.Fatalf("value:%s, should exsit", key)
			}
			if v != value {
				t.Fatalf("value should be %s, but value:%s", key, v)
			}
		}
	}

	for n := 4; n < num; n++ {
		tree, _ := StartNewTree(n, n+10)
		testkv := GenTestKeyAndValue(n)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			// fmt.Printf("insert: %s ...\n", testkv[i])
			tree.Insert(testkv[i], testkv[i])
			// tree.Print()
		}
		for i := 0; i < len(testkv); i++ {
			key := testkv[i]
			value := testkv[i]
			v, ok := tree.Find(key)
			if !ok {
				t.Fatalf("value:%s, should exsit", key)
			}
			if v != value {
				t.Fatalf("value should be %s, but value:%s", key, v)
			}
		}
	}
}

// 重复插入
func TestInsertCaseDuplicated(t *testing.T) {
	tree, _ := StartNewTree(4, 4)
	testkv := GenTestKeyAndValue(3)
	ShuffleTestkv(testkv)
	t.Log(testkv)
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
	}
	tree.Print()
	for i := 0; i < len(testkv); i++ {
		key := testkv[i]
		value := testkv[i]
		v, ok := tree.Find(key)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if v != value {
			t.Fatalf("value should be %s, but value:%s", key, v)
		}
	}
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
	}
	tree.Print()
	for i := 0; i < len(testkv); i++ {
		key := testkv[i]
		value := testkv[i]
		v, ok := tree.Find(key)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if v != value {
			t.Fatalf("value should be %s, but value:%s", key, v)
		}
	}
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
		tree.Insert(testkv[i], testkv[i])
		tree.Insert(testkv[i], testkv[i])
	}
	tree.Print()
	for i := 0; i < len(testkv); i++ {
		key := testkv[i]
		value := testkv[i]
		v, ok := tree.Find(key)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if v != value {
			t.Fatalf("value should be %s, but value:%s", key, v)
		}
	}

}
