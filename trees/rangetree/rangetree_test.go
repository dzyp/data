package rangetree

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func TestConstructTree(t *testing.T) {
	ints := []int{3, 5, 1, 2, 9, 0, 13}
	New(ints)
}

func TestGetRange(t *testing.T) {
	ints := []int{3, 5, 1, 2, 9, 0, 13}
	tree := New(ints)

	result := tree.GetRange(0, 4)
	if !reflect.DeepEqual(ints[0:4], result) {
		t.Errorf(`Expected result: %+v, received: %+v`, ints[0:5], result)
	}
}

func TestGetRangeMatchingEndpoints(t *testing.T) {
	ints := []int{3, 5, 1, 2, 9, 0, 13}
	tree := New(ints)

	result := tree.GetRange(1, 9)
	if !reflect.DeepEqual(ints[1:5], result) {
		t.Errorf(`Expected result: %+v, received: %+v`, ints[1:5], result)
	}
}

func TestCreateTreeProperly(t *testing.T) {
	ints := []int{1, 3, 5, 9, 10}
	tree := New(ints)

	left := tree.root.left
	if left.value != 3 {
		t.Errorf(`Expected value: %d, received: %d`, 3, left.value)
	}
}

func TestGetRightLeaf(t *testing.T) {
	ints := []int{1, 3, 5, 9, 10}
	tree := New(ints)

	result := tree.GetLeaf(6)
	if result.value != 9 {
		t.Errorf(`Expected value: %d, received: %d`, 9, result.value)
	}
}

func TestGetLeftLeaf(t *testing.T) {
	ints := []int{1, 3, 5, 9, 10}
	tree := New(ints)

	result := tree.GetLeaf(2)
	if result.value != 3 {
		t.Errorf(`Expected value: %d, received: %d`, 3, result.value)
	}
}

func TestGetExactLeaf(t *testing.T) {
	ints := []int{1, 3, 5, 9, 10}
	tree := New(ints)

	result := tree.GetLeaf(3)
	if result.value != 3 {
		t.Errorf(`Expected value: %d, received: %d`, 3, result.value)
	}

	if !result.isLeaf() {
		t.Errorf(`Expected leaf, received otherwise.`)
	}
}

func TestInsertLeft(t *testing.T) {
	ints := []int{1, 3, 5, 9, 10}
	tree := New(ints)

	tree.Insert(4)

	n := tree.root.left.left.right
	if n.value != 4 {
		t.Errorf(`Expected value: %d, received: %d`, 4, n.value)
	}

	if tree.root.numChildren != 6 {
		t.Errorf(
			`Expected children: %d, received: %d`, 6, tree.root.numChildren,
		)
	}

	result := tree.GetRange(3, 6)
	if !reflect.DeepEqual([]int{3, 4, 5}, result) {
		t.Errorf(`Wrong range received: %+v`, result)
	}
}

func TestInsertRight(t *testing.T) {
	ints := []int{1, 3, 5, 9, 10}
	tree := New(ints)

	tree.Insert(6)

	n := tree.root.right.left.left
	if n.value != 6 {
		t.Errorf(`Expected value: %d, received: %d`, 6, n.value)
	}

	if tree.root.numChildren != 6 {
		t.Errorf(`Expected value: %d, received: %d`, 6, tree.root.numChildren)
	}

	result := tree.GetRange(5, 10)
	if !reflect.DeepEqual([]int{5, 6, 9}, result) {
		t.Errorf(`Wrong range received: %+v`, result)
	}
}

func TestInsertDuplicate(t *testing.T) {
	ints := []int{1, 3, 5, 9, 10}
	tree := New(ints)

	tree.Insert(3)

	if tree.root.numChildren != 5 {
		t.Errorf(
			`Expected numchildren: %d, received: %d`, 5, tree.root.numChildren,
		)
	}
}

func TestRebalancing(t *testing.T) {
	ints := []int{1, 3, 5, 7}
	tree := New(ints)

	oldRebalanceRatio := REBALANCE_RATIO
	REBALANCE_RATIO = .5

	tree.Insert(2)

	if tree.root.value != 3 {
		t.Errorf(`Expected value: %d, received: %d`, 3, tree.root.value)
	}

	REBALANCE_RATIO = oldRebalanceRatio
}

func TestNonRootRebalancing(t *testing.T) {
	ints := []int{1, 3, 5, 7, 9}
	tree := New(ints)

	oldRebalanceRatio := REBALANCE_RATIO
	REBALANCE_RATIO = .25

	tree.Insert(2) // this should cause a rebalance on the root's left node

	n := tree.root.left.left
	if n.value != 2 {
		t.Errorf(`Expected value: %d, received: %d`, 2, n.value)
	}

	REBALANCE_RATIO = oldRebalanceRatio
}

func TestDeleteLeft(t *testing.T) {
	ints := []int{1, 3, 5, 7, 9}
	tree := New(ints)

	tree.Delete(3)

	r := tree.GetRange(1, 6)
	if !reflect.DeepEqual(r, []int{1, 5}) {
		t.Errorf(`Received incorrect result: %+v`, r)
	}

	n := tree.GetLeaf(3)
	if n.value == 3 {
		t.Errorf(`Received deleted value: %+v`, n.value)
	}

	if tree.root.numChildren != 4 {
		t.Errorf(
			`Children not decremented.  Received: %d.`, tree.root.numChildren,
		)
	}
}

func TestDeleteRight(t *testing.T) {
	ints := []int{1, 3, 5, 7, 9}
	tree := New(ints)

	tree.Delete(7)

	r := tree.GetRange(5, 10)
	if !reflect.DeepEqual(r, []int{5, 9}) {
		t.Errorf(`Received incorrect result: %+v`, r)
	}

	n := tree.GetLeaf(7)
	if n.value == 7 {
		t.Errorf(`Received deleted value: %+v`, n.value)
	}

	if tree.root.numChildren != 4 {
		t.Errorf(`Received numChildren: %d`, tree.root.numChildren)
	}
}

func TestDeleteCenter(t *testing.T) {
	ints := []int{1, 3, 5, 7, 9}
	tree := New(ints)

	tree.Delete(5)

	r := tree.GetRange(3, 8)
	if !reflect.DeepEqual(r, []int{3, 7}) {
		t.Errorf(`Received incorrect result: %+v`, r)
	}
}

func TestBenchmark(t *testing.T) {
	numInts := 1000000

	ints := make([]int, numInts)
	intMap := make(map[int]bool)
	for i := 0; i < numInts; i++ {
		ints[i] = i
		intMap[i] = true
	}

	tree := New(ints)

	t1 := time.Now()
	result := tree.GetRange(250000, 750000)

	log.Printf(
		`It took %d ms to fetch %d items.`,
		time.Since(t1).Nanoseconds()/int64(time.Millisecond),
		numInts,
	)

	if !reflect.DeepEqual(ints[250000:750000], result) {
		t.Errorf(`Expected result: %+v, received: %+v`, ints, result)
	}

	t2 := time.Now()
	rangeInts := make([]int, 0)
	for i := 0; i < numInts; i++ {
		if _, ok := intMap[i]; ok {
			rangeInts = append(rangeInts, i)
		}
	}

	log.Printf(
		`It took %d ms to fetch %d items.`,
		time.Since(t2).Nanoseconds()/int64(time.Millisecond),
		numInts,
	)
}
