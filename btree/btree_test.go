package btree

import (
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	bTree, _ := Create(3)
	bTree.TraversalLevel()
}

func TestInsert(t *testing.T) {
	bTree, _ := Create(3)
	bTree.Insert(40)
	bTree.Insert(30)
	bTree.Insert(50)
	bTree.Insert(10)
	bTree.Insert(20)

	bTree.Insert(5)
	bTree.Insert(35)
	bTree.Insert(55)
	bTree.Insert(60)

	bTree.Insert(70)
	bTree.Insert(80)

	bTree.Insert(90)
	bTree.Insert(100)
	bTree.Insert(63)
	bTree.Insert(65)
	bTree.Insert(52)

	bTree.TraversalLevel()
}

func TestDelete(t *testing.T) {
	bTree, _ := Create(3)
	bTree.Insert(1)
	bTree.Insert(2)
	bTree.Insert(3)
	bTree.Insert(4)
	bTree.Insert(5)
	bTree.Insert(6)

	bTree.Delete(1)
	bTree.TraversalLevel()
	bTree.Delete(4)
	bTree.TraversalLevel()
}

func TestSearch(t *testing.T) {
	bTree, _ := Create(3)
	bTree.Insert(1)
	bTree.Insert(2)
	bTree.Insert(3)
	bTree.Insert(4)
	bTree.Insert(5)
	bTree.Insert(6)
	fmt.Println(bTree.Search(4))
	fmt.Println(bTree.Search(10))
}
