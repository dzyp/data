package v1

import (
	"fmt"
	"log"
	"testing"
	"time"

	r "github.com/dzyp/data/trees/rangetree"
)

type point struct {
	coordinates [2]int
}

func (self *point) GetDimensionalValue(dimension int) int {
	return self.coordinates[dimension-1]
}

func (self *point) MaxDimensions() int {
	return len(self.coordinates)
}

func (self *point) EqualAtDimension(entry r.Entry, dimension int) bool {
	return self.GetDimensionalValue(dimension) ==
		entry.GetDimensionalValue(dimension)
}

func (self *point) LessThan(entry r.Entry, dimension int) bool {
	return self.coordinates[dimension-1] < entry.GetDimensionalValue(dimension)
}

func (self *point) String() string {
	return fmt.Sprintf(`X: %d, Y: %d`, self.coordinates[0], self.coordinates[1])
}

func (self *point) Less(other r.Entry, dimension int) bool {
	for i := self.MaxDimensions(); i >= dimension; i-- {
		if other.LessThan(self, i) {
			return false
		}
	}

	return true
}

func newPoint(x, y int) *point {
	return &point{[2]int{x, y}}
}

type bound struct {
	high int
	low  int
}

func (self *bound) High() int {
	return self.high
}

func (self *bound) Low() int {
	return self.low
}

func newBound(low, high int) *bound {
	return &bound{low: low, high: high}
}

type query struct {
	coordinates [2]*bound
}

func (self *query) GetDimensionalBounds(dimension int) r.Bounds {
	return self.coordinates[dimension-1]
}

func newQuery(startRow, stopRow, startColumn, stopColumn int) *query {
	log.Print(``)
	return &query{
		[2]*bound{
			newBound(startRow, stopRow),
			newBound(startColumn, stopColumn),
		},
	}
}

type coordinate struct {
	x int
	y int
}

func newCoordinate(x, y int) *coordinate {
	return &coordinate{x, y}
}

func checkCoordinates(t *testing.T, entry r.Entry, x, y int) {
	p := entry.(*point)

	if p.coordinates[0] != x {
		t.Errorf(`X coordinate expected: %d, received: %d`, x, p.coordinates[0])
	}

	if p.coordinates[1] != y {
		t.Errorf(`Y coordinate expected: %d, received: %d`, y, p.coordinates[1])
	}
}

func checkNumChildren(t *testing.T, n *node, numChildren int) {
	if numChildren != n.numChildren {
		t.Errorf(
			`Expected num children: %d, received: %d`,
			numChildren, n.numChildren,
		)
	}
}

func checkLen(t *testing.T, entries []r.Entry, expected int) {
	if len(entries) != expected {
		t.Errorf(`Expected len: %d, received: %d`, expected, len(entries))
	}
}

func checkEntries(t *testing.T, entries []r.Entry, expected ...*coordinate) {
	checkLen(t, entries, len(expected))

	// this is inefficient, i know, just for testing
	for _, coord := range expected {
		found := false
		for _, entry := range entries {
			p := entry.(*point)
			if coord.x == p.coordinates[0] && coord.y == p.coordinates[1] {
				found = true
			}
		}

		if !found {
			t.Errorf(`Expected: %+v, not found.`, coord)
		}
	}
}

func checkNode(t *testing.T, n *node, c *coordinate) {
	checkEntries(t, []r.Entry{n.value}, c)
}

func TestInsertTwoPoints(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)

	tree := New(2)

	tree.Insert(p1, p2)

	checkCoordinates(t, tree.root.value, 1, 1)
	checkCoordinates(t, tree.root.left.value, 0, 0)
	checkCoordinates(t, tree.root.right.value, 1, 1)
	checkNumChildren(t, tree.root, 2)
	checkNumChildren(t, tree.root.right, 0)
	checkNumChildren(t, tree.root.left, 0)
}

func TestQueryTwoPoints(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)

	tree := New(2)

	tree.Insert(p1, p2)

	entries := tree.GetRange(newQuery(0, 1, 0, 1))

	checkLen(t, entries, 1)
	checkCoordinates(t, entries[0], 0, 0)

	entries = tree.GetRange(newQuery(1, 2, 1, 2))

	checkLen(t, entries, 1)
	checkCoordinates(t, entries[0], 1, 1)

	entries = tree.GetRange(newQuery(10, 11, 10, 11))
	checkLen(t, entries, 0)
}

func TestQueryAfterEditLow(t *testing.T) {
	tree := New(2)

	p := newPoint(5, 5)
	tree.Insert(p)

	entries := tree.GetRange(newQuery(0, 10, 0, 10))

	checkLen(t, entries, 1)
	checkCoordinates(t, entries[0], 5, 5)

	p = newPoint(1, 1)
	tree.Insert(p)

	checkCoordinates(t, tree.root.left.value, 1, 1)

	entries = tree.GetRange(newQuery(0, 10, 0, 10))

	checkLen(t, entries, 2)

}

func TestQueryAfterEditHigh(t *testing.T) {
	p := newPoint(5, 5)

	tree := New(2)
	tree.Insert(p)

	p = newPoint(9, 9)
	tree.Insert(p)

	log.Printf(`rt: %+v`, tree.root.left.value)

	checkCoordinates(t, tree.root.right.value, 9, 9)

	entries := tree.GetRange(newQuery(0, 10, 0, 10))

	checkLen(t, entries, 2)
}

func TestQueryMultipleLevels(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)
	p3 := newPoint(5, 5)
	p4 := newPoint(9, 9)
	p5 := newPoint(10, 10)

	tree := New(2)

	tree.Insert(p1, p2, p3, p4, p5)

	entries := tree.GetRange(newQuery(1, 10, 1, 10))

	checkEntries(
		t, entries,
		newCoordinate(1, 1),
		newCoordinate(5, 5),
		newCoordinate(9, 9),
	)
}

func TestQueryMultipleLevlsRandomInsertion(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)
	p3 := newPoint(5, 5)
	p4 := newPoint(9, 9)
	p5 := newPoint(10, 10)

	tree := New(2)

	tree.Insert(p3, p2, p5, p4, p1)

	entries := tree.GetRange(newQuery(1, 10, 1, 10))

	checkEntries(
		t, entries,
		newCoordinate(1, 1),
		newCoordinate(5, 5),
		newCoordinate(9, 9),
	)
}

func TestQueryIdenticalFirstDimension(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(0, 1)

	tree := New(2)

	tree.Insert(p2, p1)

	entries := tree.GetRange(newQuery(0, 1, 0, 1))

	checkEntries(t, entries, newCoordinate(0, 0))
}

func TestQueryIdenticalSecondDimension(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 0)

	tree := New(2)

	tree.Insert(p1, p2)

	entries := tree.GetRange(newQuery(0, 1, 0, 1))

	checkEntries(t, entries, newCoordinate(0, 0))
}

func TestQueryIdenticalAllDimensions(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(0, 0)

	tree := New(2)

	tree.Insert(p1, p2)

	checkNumChildren(t, tree.root, 0)

	entries := tree.GetRange(newQuery(0, 10, 0, 10))

	checkEntries(t, entries, newCoordinate(0, 0))
}

func TestMiddleOfMultiDimensionalRange(t *testing.T) {
	p1 := newPoint(3, 3)
	p2 := newPoint(4, 3)
	p3 := newPoint(3, 4)
	p4 := newPoint(4, 4)

	tree := New(2)

	tree.Insert(p4, p2, p3, p1)

	entries := tree.GetRange(newQuery(3, 4, 3, 4))

	checkEntries(t, entries, newCoordinate(3, 3))
}

func TestLargeDenseMatrix(t *testing.T) {
	maxRange := 10

	tree := New(2)

	for i := 0; i < maxRange; i++ {
		for j := 0; j < maxRange; j++ {
			p := newPoint(i, j)
			tree.Insert(p)
		}
	}

	t0 := time.Now()
	entries := tree.GetRange(newQuery(0, maxRange, 0, maxRange))
	log.Printf(`time to query: %d ms`, time.Since(t0).Nanoseconds()/int64(time.Millisecond))

	checkLen(t, entries, maxRange*maxRange)
}

func TestRemoveRootNode(t *testing.T) {
	tree := New(2)

	p := newPoint(0, 0)
	tree.Insert(p)

	tree.Remove(newPoint(0, 0))

	if tree.numChildren != 0 {
		t.Errorf(`Expected num children: %d, received: %d`, 0, tree.numChildren)
	}
}

func TestRemoveFirstLevelNode(t *testing.T) {
	tree := New(2)

	p1 := newPoint(0, 0)
	p2 := newPoint(1, 0)
	tree.Insert(p1, p2)

	tree.Remove(newPoint(1, 0))

	if tree.numChildren != 1 {
		t.Errorf(`Expected num children: %d, received: %d`, 1, tree.numChildren)
	}

	entries := tree.GetRange(newQuery(0, 2, 0, 2))

	checkEntries(t, entries, newCoordinate(0, 0))
}

func TestRemoveDeepLevelNode(t *testing.T) {
	tree := New(2)

	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)
	p3 := newPoint(2, 2)
	p4 := newPoint(3, 3)

	tree.Insert(p1, p2, p3, p4)

	tree.Remove(newPoint(2, 2))

	entries := tree.GetRange(newQuery(0, 5, 0, 5))

	checkEntries(
		t, entries,
		newCoordinate(0, 0),
		newCoordinate(1, 1),
		newCoordinate(3, 3),
	)
}

func TestRemoveNodesSecondDimension(t *testing.T) {
	tree := New(2)

	p1 := newPoint(0, 0)
	p2 := newPoint(0, 1)
	p3 := newPoint(0, 3)

	tree.Insert(p1, p2, p3)

	tree.Remove(newPoint(0, 1))

	entries := tree.GetRange(newQuery(0, 1, 0, 5))

	checkEntries(t, entries, newCoordinate(0, 0), newCoordinate(0, 3))
}

func TestRemoveAllSecondDimensionNodes(t *testing.T) {
	tree := New(2)

	p1 := newPoint(0, 0)
	p2 := newPoint(0, 1)

	tree.Insert(p1, p2)

	tree.Remove(newPoint(0, 0), newPoint(0, 1))

	if tree.numChildren != 0 {
		t.Errorf(`Expected num children: %d, received: %d`, 0, tree.numChildren)
	}

	if tree.root != nil {
		t.Errorf(`Expected nil root, received: %+v`, tree.root)
	}
}

func TestFetchAll(t *testing.T) {
	tree := New(2)

	p1 := newPoint(0, 0)
	p2 := newPoint(0, 1)
	p3 := newPoint(1, 1)
	p4 := newPoint(1, 2)

	tree.Insert(p1, p2, p3, p4)

	entries := tree.All()

	checkEntries(
		t, entries,
		newCoordinate(0, 0),
		newCoordinate(0, 1),
		newCoordinate(1, 1),
		newCoordinate(1, 2),
	)
}

func TestNewNode(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)
	p3 := newPoint(1, 3)
	p4 := newPoint(2, 0)

	tree := new(2, 1, p1, p2, p3, p4)

	entries := tree.GetRange(newQuery(0, 4, 0, 4))

	checkEntries(
		t, entries,
		newCoordinate(0, 0),
		newCoordinate(1, 1),
		newCoordinate(1, 3),
		newCoordinate(2, 0),
	)

	if tree.numChildren != 4 {
		t.Errorf(`Expected num children: %d, received: %d`, 4, tree.numChildren)
	}

	checkNumChildren(t, tree.root, 3)
	checkNumChildren(t, tree.root.right, 2)
}

func TestRebalanceSingleDimension(t *testing.T) {
	p1 := newPoint(1, 0)
	p2 := newPoint(2, 0)
	p3 := newPoint(3, 0)
	p4 := newPoint(4, 0)

	tree := New(2)

	tree.Insert(p1)
	tree.Insert(p2)
	tree.Insert(p3)
	tree.Insert(p4)

	checkNumChildren(t, tree.root.right.right.right, 0)
	tree.root.rebalance(tree)

	entries := tree.GetRange(newQuery(0, 5, 0, 5))

	checkEntries(
		t, entries,
		newCoordinate(1, 0),
		newCoordinate(2, 0),
		newCoordinate(3, 0),
		newCoordinate(4, 0),
	)

	checkNode(t, tree.root, newCoordinate(3, 0))
}

func TestRebalanceSecondDimension(t *testing.T) {
	p1 := newPoint(0, 1)
	p2 := newPoint(0, 2)
	p3 := newPoint(0, 3)
	p4 := newPoint(0, 4)

	tree := New(2)

	tree.Insert(p1)
	tree.Insert(p2)
	tree.Insert(p3)
	tree.Insert(p4)

	tree = tree.root.rt
	checkNumChildren(t, tree.root.right.right.right, 0)

	tree.rebalance()

	entries := tree.GetRange(newQuery(0, 5, 0, 5))

	checkEntries(
		t, entries,
		newCoordinate(0, 1),
		newCoordinate(0, 2),
		newCoordinate(0, 3),
		newCoordinate(0, 4),
	)

	checkNode(t, tree.root, newCoordinate(0, 3))
}
