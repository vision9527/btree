package btree

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var gaintTestKv []string

func TestMain(m *testing.M) {
	gaintTestKv = GenTestKeyAndValue(10000)
	m.Run()
}

func toString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func makeTestLeafNode(keys []string, values []string) *node {
	n := &node{
		isLeaf: true,
	}
	for _, v := range keys {
		k := key(v)
		n.keys = append(n.keys, k)
	}
	for _, v := range values {
		et := &entry{
			value: v,
		}
		n.pointers = append(n.pointers, et)
	}
	return n
}

func GenTestKeyAndValue(repeatNum int) []string {
	str := "abcdefghijklmnopqrstuvwxyz"
	result := make([]string, 0)
	for _, i := range str {
		buf := bytes.Buffer{}
		for j := 0; j < repeatNum; j++ {
			buf.WriteRune(i)
			result = append(result, buf.String())
		}

	}
	return result
}

func GenTestRandomKeyAndValue(amount int, n int) []string {
	result := make([]string, 0)
	for i := 0; i < amount; i++ {
		result = append(result, RandStringRunes(n))
	}
	return result
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()_+{}:~.,")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ShuffleTestkv(kv []string) {
	rand.Seed(time.Now().Unix())
	for len(kv) > 0 {
		size := len(kv)
		r := rand.Intn(size)
		kv[r], kv[size-1] = kv[size-1], kv[r]
		kv = kv[0 : size-1]
	}
}

func TestBPlusTree_check(t *testing.T) {
	tree, _ := StartNewTree(3, 3)
	testkv := GenTestKeyAndValue(3)
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
	}
	tree.check(false)
	num := 10
	for n := 1; n < num; n++ {
		tree, _ := StartNewTree(4, 4)
		testkv := GenTestKeyAndValue(n)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			tree.Insert(testkv[i], testkv[i])
		}
		tree.check(false)
	}
}

func TestGenTestKeyAndValue_Print(t *testing.T) {
	testkv := GenTestKeyAndValue(3)
	t.Log("test_keys:", testkv)
}

func TestGenTestRandomKeyAndValue(t *testing.T) {
	testkv := GenTestRandomKeyAndValue(100, 30)
	t.Log("TestGenTestRandomKeyAndValue_keyvalue:", testkv)
}

func TestBPlusTreeFind_rootDontHaveChild(t *testing.T) {
	tree, _ := StartNewTree(4, 4)
	root := &node{
		isLeaf:   true,
		keys:     []key{key("test_key1"), key("test_key2"), key("test_key3")},
		pointers: make([]interface{}, 0),
	}
	tree.root = root
	p1 := &entry{
		value: "test_value1",
	}
	root.pointers = append(root.pointers, p1)
	p2 := &entry{
		value: "test_value2",
	}
	root.pointers = append(root.pointers, p2)
	p3 := &entry{
		value: "test_value3",
	}
	root.pointers = append(root.pointers, p3)
	v, ok := tree.Find("test_key1")
	if !ok {
		t.Fatalf("not find test_key1 value")
	}
	if toString(v) != "test_value1" {
		t.Fatalf("find test_key1 value, but not correct")
	}
	v, ok = tree.Find("test_key2")
	if !ok {
		t.Fatalf("not find test_key2 value")
	}
	if toString(v) != "test_value2" {
		t.Fatalf("find test_key2 value, but not correct")
	}
	tree.check(false)
}

func TestBPlusTreeFind_rootHaveChild(t *testing.T) {
	tree, _ := StartNewTree(6, 6)
	root := &node{
		isLeaf:   false,
		keys:     []key{key("El Said"), key("Mozart")},
		pointers: make([]interface{}, 2),
	}
	tree.root = root
	nodeOne := makeTestLeafNode([]string{"Brandt", "Califieri", "Crick", "Einstein"},
		[]string{"Brandt_value", "Califieri_value", "Crick_value", "Einstein_value"})
	nodeTwo := makeTestLeafNode([]string{"El Said", "Gold", "Katz", "Kim"},
		[]string{"El Said_value", "Gold_value", "Katz_value", "Kim_value"})
	nodeThree := makeTestLeafNode([]string{"Mozart", "Singh", "Srinivasan", "Wu"},
		[]string{"Mozart_value", "Singh_value", "Srinivasan_value", "Wu_value"})
	nodeOne.maxSize = 6
	nodeTwo.maxSize = 6
	nodeThree.maxSize = 6
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
	if toString(v) != "Mozart_value" {
		t.Fatalf("find Mozart_value, but not correct")
	}
	v, ok = tree.Find("Wu")
	if !ok {
		t.Fatalf("not find Wu value")
	}
	if toString(v) != "Wu_value" {
		t.Fatalf("find Wu_value, but not correct")
	}
	tree.check(false)

}

func TestBPlusTree_insertIntoLeaf(t *testing.T) {
	tree, _ := StartNewTree(6, 6)
	leafNode := makeTestLeafNode([]string{"key1", "key2", "key4", "key5"},
		[]string{"key1_value", "key2_value", "key4_value", "key5_value"})
	leafNode.maxSize = 6
	targetKey := key("key3")
	et := &entry{
		value: "key3_value",
	}
	tree.insertIntoLeaf(leafNode, targetKey, et)
	if targetKey.compare(leafNode.keys[2]) != 0 {
		t.Fatalf("should be key3")
	}
	r, ok := leafNode.pointers[2].(*entry)
	if !ok {
		t.Fatalf("should be et")
	}
	if r.toValue() != "key3_value" {
		t.Fatalf("should be key3_value")
	}
	t.Logf("keys: %v", leafNode.keys)
	values := make([]string, 0)
	for _, i := range leafNode.pointers {
		r, _ := i.(*entry)
		values = append(values, r.toValue())
	}
	t.Logf("et: %v", values)
	targetKey = key("key0")
	et = &entry{
		value: "key0_value",
	}
	tree.insertIntoLeaf(leafNode, targetKey, et)
	if targetKey.compare(leafNode.keys[0]) != 0 {
		t.Fatalf("should be key0")
	}
	t.Logf("keys: %v", leafNode.keys)
	targetKey = key("key6")
	et = &entry{
		value: "key6_value",
	}
	tree.insertIntoLeaf(leafNode, targetKey, et)
	if targetKey.compare(leafNode.keys[6]) != 0 {
		t.Fatalf("should be key6")
	}
	r, ok = leafNode.pointers[6].(*entry)
	if !ok {
		t.Fatalf("should be et")
	}
	if r.toValue() != "key6_value" {
		t.Fatalf("should be key6_value")
	}
	t.Logf("keys: %v", leafNode.keys)
	tree.check(false)
}

func TestInsertCaseOne(t *testing.T) {
	tree, _ := StartNewTree(3, 3)
	tree.Insert("a", "a")
	tree.Print()
	tree.Insert("b", "b")
	tree.Print()
	tree.Insert("c", "c")
	tree.Print()
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
	if toString(v) != "a" {
		t.Fatalf("value should be a, but value:%v", v)
	}
	v, _ = tree.Find("g")
	if toString(v) != "g" {
		t.Fatalf("value should be g, but value:%v", v)
	}
	v, _ = tree.Find("i")
	if toString(v) != "i" {
		t.Fatalf("value should be i, but value:%v", v)
	}
	v, _ = tree.Find("m")
	if toString(v) != "m" {
		t.Fatalf("value should be m, but value:%v", v)
	}
	v, _ = tree.Find("n")
	if toString(v) != "n" {
		t.Fatalf("value should be n, but value:%v", v)
	}
	v, ok := tree.Find("xxxxxx")
	if ok {
		t.Fatalf("value should not exsit")
	}
	if v != nil {
		t.Fatalf("value should be empty")
	}
	tree.check(false)
	firstLeafNode := tree.findFirstLeafNode()
	firstLeafNode.checkOrder()
}

func TestInsertCaseTwo(t *testing.T) {
	tree, _ := StartNewTree(3, 3)
	testkv := GenTestKeyAndValue(3)
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
	}
	tree.Print()
	for i := 0; i < len(testkv); i++ {
		ky := testkv[i]
		value := testkv[i]
		v, ok := tree.Find(ky)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if toString(v) != value {
			t.Fatalf("value should be %s, but value:%v", ky, v)
		}
	}
	tree.check(false)
	firstLeafNode := tree.findFirstLeafNode()
	firstLeafNode.checkOrder()

}

func TestInsertCaseShuffleTestkv1(t *testing.T) {
	for j := 3; j <= 10; j++ {
		tree, _ := StartNewTree(j, j)
		shuffleTestKv := make([]string, len(gaintTestKv))
		copy(shuffleTestKv, gaintTestKv)
		ShuffleTestkv(shuffleTestKv)
		for i := 0; i < len(shuffleTestKv); i++ {
			tree.Insert(shuffleTestKv[i], shuffleTestKv[i])
		}
		for i := 0; i < len(shuffleTestKv); i++ {
			ky := shuffleTestKv[i]
			value := shuffleTestKv[i]
			v, ok := tree.Find(ky)
			if !ok {
				t.Fatalf("value:%s, should exsit", ky)
			}
			if toString(v) != value {
				t.Fatalf("value should be %s, but value:%v", ky, v)
			}
		}
		firstLeafNode := tree.findFirstLeafNode()
		currentNode := firstLeafNode
		count := 0
		if firstLeafNode == nil {
			panic("should not be nil")
		}
		for currentNode != nil {
			for _, k := range currentNode.keys {
				tk := gaintTestKv[count]
				if k.compare(key(tk)) != 0 {
					t.Fatalf("should be:%s", tk)
				}
				count++
			}
			currentNode = currentNode.lastOrNextNode
		}
		t.Logf("node size:%d key total:%d find total:%d\n", j, len(gaintTestKv), count)
	}
}

func TestInsertCaseShuffleTestkv2(t *testing.T) {
	tree, _ := StartDefaultNewTree()
	shuffleTestKv := make([]string, len(gaintTestKv))
	copy(shuffleTestKv, gaintTestKv)
	ShuffleTestkv(shuffleTestKv)
	for i := 0; i < len(shuffleTestKv); i++ {
		tree.Insert(shuffleTestKv[i], shuffleTestKv[i])
	}
	for i := 0; i < len(shuffleTestKv); i++ {
		ky := shuffleTestKv[i]
		value := shuffleTestKv[i]
		v, ok := tree.Find(ky)
		if !ok {
			t.Fatalf("value:%s, should exsit", ky)
		}
		if toString(v) != value {
			t.Fatalf("value should be %s, but value:%v", ky, v)
		}
	}
	firstLeafNode := tree.findFirstLeafNode()
	currentNode := firstLeafNode
	count := 0
	if firstLeafNode == nil {
		panic("should not be nil")
	}
	for currentNode != nil {
		for _, k := range currentNode.keys {
			tk := gaintTestKv[count]
			if k.compare(key(tk)) != 0 {
				t.Fatalf("should be:%s", tk)
			}
			count++
		}
		currentNode = currentNode.lastOrNextNode
	}
	t.Logf("key total:%d find total:%d\n", len(gaintTestKv), count)
}

func TestInsertCaseShuffleTestkv3(t *testing.T) {
	for n := 0; n < 5; n++ {
		tree, _ := StartNewTree(10, 10)
		testkv := GenTestRandomKeyAndValue(100000, 10)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			ky := testkv[i]
			value := ky + "_" + "v"
			tree.Insert(ky, value)
			v, ok := tree.Find(ky)
			if !ok {
				t.Fatalf("value:%s, should exsit", ky)
			}
			if toString(v) != value {
				t.Fatalf("value should be %s, but value:%v", ky, v)
			}
		}
		tree.check(false)
		firstLeafNode := tree.findFirstLeafNode()
		firstLeafNode.checkOrder()
	}
}

func TestInsertCaseShuffleTestkv4(t *testing.T) {
	testkv := GenTestRandomKeyAndValue(10000, 10)
	for n := 3; n < 100; n++ {
		tree, _ := StartNewTree(n+10, n)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			tree.Insert(testkv[i], testkv[i])
		}
		for i := 0; i < len(testkv); i++ {
			ky := testkv[i]
			value := testkv[i]
			v, ok := tree.Find(ky)
			if !ok {
				t.Fatalf("value:%s, should exsit", ky)
			}
			if toString(v) != value {
				t.Fatalf("value should be %s, but value:%v", ky, v)
			}
		}
		tree.check(false)
		firstLeafNode := tree.findFirstLeafNode()
		firstLeafNode.checkOrder()
	}

	for n := 3; n < 100; n++ {
		tree, _ := StartNewTree(n, n+10)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			tree.Insert(testkv[i], testkv[i])
		}
		for i := 0; i < len(testkv); i++ {
			ky := testkv[i]
			value := testkv[i]
			v, ok := tree.Find(ky)
			if !ok {
				t.Fatalf("value:%s, should exsit", ky)
			}
			if toString(v) != value {
				t.Fatalf("value should be %s, but value:%v", ky, v)
			}
		}
		tree.check(false)
		firstLeafNode := tree.findFirstLeafNode()
		firstLeafNode.checkOrder()
	}
}

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
		ky := testkv[i]
		value := testkv[i]
		v, ok := tree.Find(ky)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if toString(v) != value {
			t.Fatalf("value should be %s, but value:%v", ky, v)
		}
	}
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
	}
	tree.Print()
	for i := 0; i < len(testkv); i++ {
		ky := testkv[i]
		value := testkv[i]
		v, ok := tree.Find(ky)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if toString(v) != value {
			t.Fatalf("value should be %s, but value:%v", ky, v)
		}
	}
	tree.check(false)
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
		tree.check(false)
		tree.Insert(testkv[i], testkv[i])
		tree.check(false)
		tree.Insert(testkv[i], testkv[i])
		tree.check(false)
	}
	tree.Print()
	for i := 0; i < len(testkv); i++ {
		ky := testkv[i]
		value := testkv[i]
		v, ok := tree.Find(ky)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if toString(v) != value {
			t.Fatalf("value should be %s, but value:%v", ky, v)
		}
	}

}

func TestBPlusTree_findFirstLeafNode(t *testing.T) {
	tree, _ := StartNewTree(5, 5)
	testkv := GenTestKeyAndValue(1000)
	for i := 0; i < len(testkv); i++ {
		ky := testkv[i]
		value := ky
		tree.Insert(ky, value)
	}
	firstLeafNode := tree.findFirstLeafNode()
	currentNode := firstLeafNode
	count := 0
	for currentNode != nil {
		for _, k := range currentNode.keys {
			tk := testkv[count]
			if k.compare(key(tk)) != 0 {
				t.Fatalf("should be:%s", tk)
			}
			count++
		}
		currentNode = currentNode.lastOrNextNode
	}
}

func TestBPlusTree_FindRangeOrder(t *testing.T) {
	tree, _ := StartNewTree(5, 5)
	testkv := GenTestKeyAndValue(10)
	for i := 0; i < len(testkv); i++ {
		ky := testkv[i]
		value := ky
		tree.Insert(ky, value)
	}
	startIndex := 1
	endIndex := 15
	result := tree.FindRange(testkv[startIndex], testkv[endIndex])
	t.Logf("testkv[startIndex]=%s, testkv[endIndex]=%s \n", testkv[startIndex], testkv[endIndex])
	t.Logf("result=%s \n", result)
	t.Logf("range:%d, result length=%d \n", endIndex-startIndex+1, len(result))
	if testkv[startIndex] != result[0] {
		t.Fatalf("result first should be: %s", testkv[startIndex])
	}
	if testkv[endIndex] != result[len(result)-1] {
		t.Fatalf("result last should be: %s", testkv[endIndex])
	}
	if len(result) != (endIndex - startIndex + 1) {
		t.Fatalf("result length should be: %s", testkv[endIndex])
	}
	for i := startIndex; i <= endIndex; i++ {
		if testkv[i] != result[i-startIndex] {
			t.Fatalf("result index:%d should be: %s", i-startIndex, testkv[i])
		}
	}
	result = tree.FindRange(testkv[endIndex], testkv[startIndex])
	if len(result) != 0 {
		t.Fatalf("result length should be:0, result:%v", result)
	}
	tree.check(false)
}

func TestBPlusTree_FindRangeShuffle(t *testing.T) {
	tree, _ := StartNewTree(3, 3)
	testkv := GenTestKeyAndValue(1000)
	originTestKv := make([]string, len(testkv))
	copy(originTestKv, testkv)
	start := "aa"
	end := "cccccc"
	startIndex := -1
	endIndex := -1
	for i := 0; i < len(originTestKv); i++ {
		kv := originTestKv[i]
		if kv == start {
			startIndex = i
		}
		if kv == end {
			endIndex = i
		}
		if startIndex != -1 && endIndex != -1 {
			break
		}
	}
	ShuffleTestkv(testkv)
	for i := 0; i < len(testkv); i++ {
		ky := testkv[i]
		value := ky
		tree.Insert(ky, value)
	}

	result := tree.FindRange(start, end)
	t.Logf("startIndex=%d, endIndex=%d \n", startIndex, endIndex)
	t.Logf("start=%s, end=%s \n", start, end)
	t.Logf("length:%d \n", len(result))
	if start != toString(result[0]) {
		t.Fatalf("result first should be: %s", start)
	}
	if end != toString(result[len(result)-1]) {
		t.Fatalf("result last should be: %s", end)
	}

	if originTestKv[startIndex] != toString(result[0]) {
		t.Fatalf("result first should be: %s", originTestKv[startIndex])
	}
	if originTestKv[endIndex] != toString(result[len(result)-1]) {
		t.Fatalf("result last should be: %s", testkv[endIndex])
	}
	if len(result) != (endIndex - startIndex + 1) {
		t.Fatalf("result length should be: %s", originTestKv[endIndex])
	}
	for i := startIndex; i <= endIndex; i++ {
		if originTestKv[i] != toString(result[i-startIndex]) {
			t.Fatalf("result index:%d should be: %s", i-startIndex, originTestKv[i])
		}
	}
	t.Logf("load node count:%d", tree.GetCount())
	tree.check(false)
}

func TestBPlusTree_DeleteOne(t *testing.T) {
	tree, _ := StartNewTree(3, 3)
	tree.Insert("a", "a")
	tree.Print()
	tree.Insert("b", "b")
	tree.Print()
	tree.Insert("c", "c")
	tree.Print()
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
	a, ok := tree.Delete("a")
	if !ok {
		t.Fatal("shoul have a")
	}
	fmt.Println(a.(string))
	tree.Print()
	b, ok := tree.Delete("b")
	if !ok {
		t.Fatal("shoul have b")
	}
	fmt.Println(b.(string))
	g, ok := tree.Delete("g")
	if !ok {
		t.Fatal("shoul have g")
	}
	fmt.Println(g.(string))
	tree.check(false)
	firstLeafNode := tree.findFirstLeafNode()
	firstLeafNode.checkOrder()
}

func TestBPlusTree_DeleteTwo(t *testing.T) {
	tree, _ := StartNewTree(3, 3)
	testkv := GenTestKeyAndValue(3)
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
		tree.check(false)
	}
	tree.Print()
	for i := 0; i < len(testkv); i++ {
		ky := testkv[i]
		value := testkv[i]
		v, ok := tree.Delete(ky)
		tree.check(false)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if toString(v) != value {
			t.Fatalf("value should be %s, but value:%v", ky, v)
		}
	}

}

func TestBPlusTree_DeleteThree(t *testing.T) {
	num := 100
	for n := 4; n < num; n++ {
		tree, _ := StartNewTree(n+10, n)
		testkv := GenTestKeyAndValue(n)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			ky := testkv[i]
			value := testkv[i]
			tree.Insert(ky, value)
			v, ok := tree.Delete(ky)
			tree.check(false)
			if !ok {
				t.Fatalf("value should exsit")
			}
			if toString(v) != value {
				t.Fatalf("value should be %s, but value:%v", ky, v)
			}
		}
	}
}

func TestBPlusTree_DeleteFour(t *testing.T) {
	tree, _ := StartNewTree(4, 4)
	testkv := GenTestKeyAndValue(4)
	ShuffleTestkv(testkv)
	for i := 0; i < len(testkv); i++ {
		tree.Insert(testkv[i], testkv[i])
	}
	tree.Print()
	for i := 0; i < len(testkv); i++ {
		ky := testkv[i]
		value := testkv[i]
		v, ok := tree.Delete(ky)
		tree.check(false)
		if !ok {
			t.Fatalf("value should exsit")
		}
		if toString(v) != value {
			t.Fatalf("value should be %s, but value:%v", ky, v)
		}
	}

	for n := 4; n < 100; n++ {
		tree, _ := StartNewTree(4, 4)
		testkv := GenTestKeyAndValue(n)
		ShuffleTestkv(testkv)
		for i := 0; i < len(testkv); i++ {
			ky := testkv[i]
			value := testkv[i]
			tree.Insert(ky, value)
		}
		for i := 0; i < len(testkv); i++ {
			ky := testkv[i]
			value := testkv[i]
			v, ok := tree.Delete(ky)
			tree.check(false)
			if !ok {
				tree.Print()
				t.Fatalf("value should exsit")
			}
			if toString(v) != value {
				t.Fatalf("value should be %s, but value:%v", ky, v)
			}
			firstLeafNode := tree.findFirstLeafNode()
			if firstLeafNode == nil {
				continue
			}
			firstLeafNode.checkOrder()
		}
	}

}

func TestBPlusTree_DeleteFive(t *testing.T) {
	testkv := GenTestKeyAndValue(100)
	for n := 3; n <= 10; n++ {
		tree, _ := StartNewTree(n, n)
		shuffleTestKv := make([]string, len(testkv))
		copy(shuffleTestKv, testkv)
		ShuffleTestkv(shuffleTestKv)
		for i := 0; i < len(shuffleTestKv); i++ {
			ky := shuffleTestKv[i]
			value := shuffleTestKv[i]
			tree.Insert(ky, value)
		}

		for i := 0; i < len(shuffleTestKv); i++ {
			ky := shuffleTestKv[i]
			value := shuffleTestKv[i]
			v, ok := tree.Delete(ky)
			if !ok {
				tree.Print()
				t.Fatalf("value should exsit")
			}
			if toString(v) != value {
				t.Fatalf("value should be %s, but value:%v", ky, v)
			}
			firstLeafNode := tree.findFirstLeafNode()
			if firstLeafNode == nil {
				continue
			}
			firstLeafNode.checkOrder()
		}

	}
}
