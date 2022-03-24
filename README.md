## Go implementation of B+ tree

This is a go implementation of B plus tree. The algorithm is from Database System Concepts 7th Edition Chapter 14.

You can submit issues if you have any questions.

Welcome to submit mr too.

## Getting Started
### Installing
```$ go get github.com/vision9527/btree```

### Using
```go
package main

import (
	"fmt"

	"github.com/vision9527/btree"
)

func main() {
	tree, _ := btree.StartDefaultNewTree()
	k := "a"
	v := "ack"
	tree.Insert(k, v)
	value, ok := tree.Find(k)
	fmt.Println("value: ", value, ok)
	tree.Insert("b", "bson")
	tree.Insert("c", "canal")
	tree.Insert("d", "django")
	values := tree.FindRange("a", "c")
	fmt.Println("values: ", values)
}
```

### Viewing the process of insert
```go
package main

import (
	"github.com/vision9527/btree"
)

func main() {
	tree, _ := btree.StartNewTree(3, 3)
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
}

```
![insert](./doc/the_process_of_insert.jpg)

### Viewing the process of delete
```go
package main

import (
	"github.com/vision9527/btree"
)

func main() {
	tree, _ := btree.StartNewTree(3, 3)
	tree.Insert("a", "a")
	tree.Insert("b", "b")
	tree.Insert("c", "c")
	tree.Insert("d", "d")
	tree.Insert("e", "e")
	tree.Insert("f", "f")
	tree.Insert("g", "g")
	tree.Insert("h", "h")
	tree.Insert("i", "i")
	tree.Insert("j", "j")
	tree.Insert("k", "k")
	tree.Insert("l", "l")
	tree.Insert("m", "m")
	tree.Insert("n", "n")
	tree.Print()
	tree.Delete("a")
	tree.Print()
	tree.Delete("h")
	tree.Print()
	tree.Delete("g")
	tree.Print()
	tree.Delete("b")
	tree.Print()
	tree.Delete("c")
	tree.Print()
	tree.Delete("a")
	tree.Print()
	tree.Delete("l")
	tree.Print()
	tree.Delete("g")
	tree.Print()
	tree.Delete("n")
	tree.Print()
	tree.Delete("k")
	tree.Print()
}
```
![delete](./doc/the_process_of_delete.jpg)