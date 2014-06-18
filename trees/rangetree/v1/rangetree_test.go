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

func (self *point) x() int {
	return self.coordinates[0]
}

func (self *point) y() int {
	return self.coordinates[1]
}

func (self *point) GetDimensionalValue(dimension int) int {
	return self.coordinates[dimension-1]
}

func (self *point) MaxDimensions() int {
	return len(self.coordinates)
}

func (self *point) String() string {
	return fmt.Sprintf(`X: %d, Y: %d`, self.coordinates[0], self.coordinates[1])
}

func (self *point) Less(other r.Entry, dimension int) bool {
	var selfValue int
	var otherValue int

	for i := 1; i <= self.MaxDimensions(); i++ {
		selfValue = self.GetDimensionalValue(i)
		otherValue = other.GetDimensionalValue(i)
		if selfValue > otherValue {
			return false
		} else if selfValue == otherValue {
			continue
		} else {
			return true
		}
	}

	return false
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

func TestCreateTree(t *testing.T) {
	tree := newTree(2, 1)

	if tree.maxDimensions != 2 {
		t.Errorf(
			`Expected max dimensions: %d, received: %d`, 2, tree.maxDimensions,
		)
	}

	if tree.dimension != 1 {
		t.Errorf(`Expected dimension: %d, received: %d`, 1, tree.dimension)
	}
}

func TestInsertSingleValue(t *testing.T) {
	point := newPoint(0, 0)
	tree := newTree(2, 1, point)

	if tree.root.value != 0 {
		t.Errorf(`Expected value: %d, received: %d`, 0, tree.root.value)
	}

	if tree.root.numChildren != 1 {
		t.Errorf(
			`Expected num children: %d, received: %d`, 1, tree.root.numChildren,
		)
	}

	if tree.root.p.(*orderedList).nodes[0] != point {
		t.Errorf(
			`Expected: %+v, received: %+v`,
			point,
			tree.root.p.(*orderedList).nodes[0],
		)
	}

	if tree.numChildren != 1 {
		t.Errorf(`Expected num children: %d, received: %d`, 1, tree.numChildren)
	}

	log.Printf(`tree.root.p: %+v`, tree.root.p.(*orderedList).nodes)
}

func TestInsertMultipleValues(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 0)
	tree := newTree(2, 1, p1, p2)

	if tree.numChildren != 2 {
		t.Errorf(`Expected numchildren: %d, received: %d`, 2, tree.numChildren)
	}

	if tree.root.value != 0 {
		t.Errorf(`Expected value: %d, received: %d`, 0, tree.root.value)
	}

	if tree.root.left.value != 0 {
		t.Errorf(`Expected value: %d, received: %d`, 0, tree.root.left.value)
	}

	if tree.root.right.value != 1 {
		t.Errorf(`Expected value: %d, received: %d`, 1, tree.root.right.value)
	}

	if len(tree.root.p.(*orderedList).nodes) != 2 {
		t.Errorf(
			`Expected len: %d, received: %d`,
			2,
			len(tree.root.p.(*orderedList).nodes),
		)
	}

	if len(tree.root.right.p.(*orderedList).nodes) != 1 {
		t.Errorf(
			`Expected len: %d, received: %d`,
			1,
			len(tree.root.right.p.(*orderedList).nodes),
		)
	}

	if len(tree.root.left.p.(*orderedList).nodes) != 1 {
		t.Errorf(
			`Expected len: %d, received: %d`,
			1,
			len(tree.root.left.p.(*orderedList).nodes),
		)
	}

	log.Printf(`NODE: %+v`, tree.root.left.p.(*orderedList).nodes)
}

func TestInsertMultipleDuplicateValues(t *testing.T) {
	p1 := newPoint(0, 1)
	p2 := newPoint(0, 2)
	p3 := newPoint(0, 3)

	tree := newTree(2, 1, p1, p2, p3)

	if tree.root.numChildren != 3 {
		t.Errorf(
			`Expected num children: %d, received: %d`, 3, tree.root.numChildren,
		)
	}

	if !tree.root.isLeaf() {
		t.Errorf(`Root should be leaf.`)
	}

	if len(tree.root.p.(*orderedList).nodes) != 3 {
		t.Errorf(
			`Expected len: %d, received: %d`,
			3,
			len(tree.root.p.(*orderedList).nodes),
		)
	}
}

func TestOneDimensionalQuery(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 0)

	tree := newTree(2, 1, p1, p2)

	results := tree.GetRange(newQuery(0, 2, 0, 2))

	checkEntries(t, results, p1, p2)
}

func TestTwoDimensionalQuery(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(0, 1)

	tree := newTree(2, 1, p1, p2)

	results := tree.GetRange(newQuery(0, 2, 0, 2))

	checkEntries(t, results, p1, p2)
}

func TestDenseMatrixQuery(t *testing.T) {
	maxRange := 1000

	tree := New(2)

	points := make([]r.Entry, maxRange*maxRange)
	index := 0

	for i := 0; i < maxRange; i++ {
		for j := 0; j < maxRange; j++ {
			p := newPoint(i, j)
			points[index] = p
			index++
		}
	}

	//tree.rebalance()

	tree.Insert(points...)
	//log.Printf(`NODE: %+v`, tree.root.left.left)

	points = points[0:index]

	t0 := time.Now()
	entries := tree.GetRange(newQuery(0, maxRange, 0, maxRange))
	log.Printf(`time to query: %d ms`, time.Since(t0).Nanoseconds()/int64(time.Millisecond))

	checkLen(t, entries, index)

	mp := make(map[int]map[int]*point)

	for _, entry := range points {
		p := entry.(*point)
		if _, ok := mp[p.x()]; !ok {
			mp[p.x()] = make(map[int]*point)
		}

		mp[p.x()][p.y()] = p
	}

	t0 = time.Now()
	index = 0
	for i := 0; i < maxRange; i++ {
		if _, ok := mp[i]; !ok {
			continue
		}

		for j := 0; j < maxRange; j++ {
			if p, ok := mp[i][j]; ok {
				points[index] = p
				index++
			}
		}
	}

	log.Printf(`time to query map: %d`, time.Since(t0).Nanoseconds()/int64(time.Millisecond))
}

/*
func checkCoordinates(t *testing.T, entry r.Entry, x, y int) {
	p, ok := entry.(*point)
	if !ok {
		t.Errorf(`Entry does not exist.`)
		return
	}

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

func checkNode(t *testing.T, n *node, c *coordinate) {
	checkEntries(t, []r.Entry{n.entry}, c)
}

func checkValue(t *testing.T, n *node, value int) {
	if n.value != value {
		t.Errorf(`Expected value: %d, received: %d`, value, n.value)
	}
}

func TestInsertTwoPoints(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)

	tree := New(2)

	tree.Insert(p1)
	tree.Insert(p2)

	checkValue(t, tree.root.right, 1)
	checkValue(t, tree.root.left, 0)
	checkValue(t, tree.root.left.rt.root, 0)
	checkValue(t, tree.root.right.rt.root, 1)
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

	checkValue(t, tree.root, 5)
	checkValue(t, tree.root.right, 5)
	checkValue(t, tree.root.left, 1)

	entries = tree.GetRange(newQuery(0, 10, 0, 10))

	checkLen(t, entries, 2)

}

func TestQueryAfterEditHigh(t *testing.T) {
	p := newPoint(5, 5)

	tree := New(2)
	tree.Insert(p)

	p = newPoint(9, 9)
	tree.Insert(p)

	checkValue(t, tree.root, 9)
	checkValue(t, tree.root.left, 5)
	checkValue(t, tree.root.right, 9)

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

	tree.Insert(p1)
	tree.Insert(p2)
	tree.Insert(p3)
	tree.Insert(p4)
	tree.Insert(p5)

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
	maxRange := 9

	tree := New(2)

	points := make([]r.Entry, maxRange*maxRange)
	index := 0

	for i := 0; i < maxRange; i++ {
		for j := 0; j < maxRange; j++ {
			p := newPoint(i, j)
			//tree.Insert(p)
			points[index] = p
			index++
		}
	}

	//tree.rebalance()

	tree.Insert(points...)
	log.Printf(`NODE: %+v`, tree.root.left.left)

	points = points[0:index]

	t0 := time.Now()
	println(`SHIT STARTS HERE`)
	entries := tree.GetRange(newQuery(0, maxRange, 0, maxRange))
	log.Printf(`time to query: %d ms`, time.Since(t0).Nanoseconds()/int64(time.Millisecond))

	checkLen(t, entries, index)

	mp := make(map[int]map[int]*point)

	for _, entry := range points {
		p := entry.(*point)
		if _, ok := mp[p.x()]; !ok {
			mp[p.x()] = make(map[int]*point)
		}

		mp[p.x()][p.y()] = p
	}

	t0 = time.Now()
	index = 0
	for i := 0; i < maxRange; i++ {
		if _, ok := mp[i]; !ok {
			continue
		}

		for j := 0; j < maxRange; j++ {
			if p, ok := mp[i][j]; ok {
				points[index] = p
				index++
			}
		}
	}

	log.Printf(`time to query map: %d`, time.Since(t0).Nanoseconds()/int64(time.Millisecond))
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

	checkValue(t, tree.root, 3)
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

	tree.rebalance()

	entries := tree.GetRange(newQuery(0, 5, 0, 5))

	checkEntries(
		t, entries,
		newCoordinate(0, 1),
		newCoordinate(0, 2),
		newCoordinate(0, 3),
		newCoordinate(0, 4),
	)

	checkValue(t, tree.root, 3)
}

func TestInsertOverwritesRoot(t *testing.T) {
	p1 := newPoint(0, 0)

	tree := New(2)

	tree.Insert(p1)

	p1 = newPoint(0, 0)

	tree.Insert(p1)

	if tree.root.rt.root.entry != p1 {
		t.Errorf(
			`Expected entry: %+v, received: %+v`, p1, tree.root.rt.root.entry,
		)
	}
}

func TestInsertOverwritesRight(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)

	tree := New(2)

	tree.Insert(p1, p2)

	p2 = newPoint(1, 1)

	tree.Insert(p2)

	if tree.root.right.rt.root.entry != p2 {
		t.Errorf(
			`Expected entry: %+v, received: %+v`,
			p2, tree.root.right.rt.root.entry,
		)
	}
}

func TestInsertOverwritesLeft(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)

	tree := New(2)

	tree.Insert(p1, p2)

	p2 = newPoint(0, 0)

	tree.Insert(p2)

	if tree.root.left.rt.root.entry != p2 {
		t.Errorf(
			`Expected entry: %+v, received: %+v`,
			p2, tree.root.left.rt.root.entry,
		)
	}
}

func TestInsertMultipleOverwrites(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)

	tree := New(2)

	tree.Insert(p1, p2)

	p1 = newPoint(0, 0)
	p2 = newPoint(1, 1)

	tree.Insert(p1, p2)

	if tree.root.left.rt.root.entry != p1 {
		t.Errorf(
			`Expected entry: %+v, received: %+v`,
			p2, tree.root.left.rt.root.entry,
		)
	}

	if tree.root.right.rt.root.entry != p2 {
		t.Errorf(
			`Expected entry: %+v, received: %+v`,
			p2, tree.root.right.rt.root.entry,
		)
	}
}

func TestInsertWideRanges(t *testing.T) {
	p1 := newPoint(0, 3)
	p2 := newPoint(1, 0)

	tree := New(2)

	tree.Insert(p1)
	tree.Insert(p2)

	entries := tree.GetRange(newQuery(0, 4, 0, 4))

	checkEntries(t, entries, newCoordinate(0, 3), newCoordinate(1, 0))
}

func BenchmarkFirstDimensionRange(b *testing.B) {
	log.Printf(`N: %d`, b.N)
	numItems := 10
	points := make([]r.Entry, numItems)
	for i := 0; i < numItems; i++ {
		points[i] = newPoint(i, 0)
	}

	tree := New(2, points...)

	var results []r.Entry

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results = tree.GetRange(newQuery(0, numItems, 0, numItems))
	}

	if len(points) != len(results) {
		b.Errorf(`Expected len: %d, received: %d`, len(points), len(results))
	}
}

func BenchmarkCheckSecondDimensionRange(b *testing.B) {
	numItems := 20000

	tree := New(2)

	points := make([]r.Entry, numItems)

	for i := 0; i < numItems; i++ {
		points[i] = newPoint(0, i)
	}

	tree.Insert(points...)

	tree.Insert(newPoint(1, 0))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.GetRange(newQuery(0, numItems, 0, numItems))
	}
}

func BenchmarkInsertCells(b *testing.B) {
	numItems := 100000

	points := make([]r.Entry, numItems)

	for i := 0; i < numItems; i++ {
		points[i] = newPoint(0, i)
	}

	var t *tree

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t = New(2)
		t.Insert(points...)
	}

	b.StopTimer()

	entries := t.GetRange(newQuery(0, numItems, 0, numItems))

	if len(entries) != numItems {
		b.Errorf(`Expected num items: %d, received: %d`, numItems, len(entries))
	}
}

func BenchmarkEditCells(b *testing.B) {
	numItems := 100000

	points := make([]r.Entry, numItems)

	for i := 0; i < numItems; i++ {
		points[i] = newPoint(0, i)
	}

	t := New(2)

	t.Insert(points...)

	points = make([]r.Entry, 100)

	for i := 0; i < 100; i++ {
		points[i] = newPoint(1, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.Insert(points...)
	}

	b.StopTimer()

	entries := t.GetRange(newQuery(0, numItems, 0, numItems))

	if len(entries) != numItems+100 {
		b.Errorf(
			`Expected num items: %d, received: %d`, numItems+100, len(entries),
		)
	}
}

func TestProperQuery(t *testing.T) {
	points := make([]r.Entry, 7)

	for i := 0; i < 7; i++ {
		points[i] = newPoint(i, 0)
	}

	tree := New(2, points...)

	entries := tree.GetRange(newQuery(0, 7, 0, 7))

	log.Printf(`ENTRIES: %+v`, entries)
	log.Printf(`NODE: %+v`, tree.root.left.left)

	if len(entries) != 7 {
		t.Errorf(`fail`)
	}
}*/
