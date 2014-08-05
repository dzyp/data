package v1

import (
	"fmt"
	"log"
	"reflect"
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

	if tree.root.len() != 1 {
		t.Errorf(
			`Expected num children: %d, received: %d`, 1, tree.root.len(),
		)
	}

	if tree.root.p.(*orderedList).nodes[0] != point {
		t.Errorf(
			`Expected: %+v, received: %+v`,
			point,
			tree.root.p.(*orderedList).nodes[0],
		)
	}

	if tree.Len() != 1 {
		t.Errorf(`Expected num children: %d, received: %d`, 1, tree.Len())
	}

	log.Printf(`tree.root.p: %+v`, tree.root.p.(*orderedList).nodes)
}

func TestInsertMultipleValues(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 0)
	tree := newTree(2, 1, p1, p2)

	if tree.Len() != 2 {
		t.Errorf(`Expected numchildren: %d, received: %d`, 2, tree.Len())
	}

	if tree.root.value != 1 {
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

	log.Printf(`NODE: %+v`, tree.root.p.(*orderedList).nodes)
}

func TestInsertMultipleDuplicateValues(t *testing.T) {
	p1 := newPoint(0, 1)
	p2 := newPoint(0, 2)
	p3 := newPoint(0, 3)

	tree := newTree(2, 1, p1, p2, p3)

	if tree.root.len() != 3 {
		t.Errorf(
			`Expected num children: %d, received: %d`, 3, tree.root.len(),
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

func TestInsertAfterCreation(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 0)

	tree := newTree(2, 1, p1, p2)

	p3 := newPoint(2, 0)
	tree.Insert(p3)

	log.Printf(`ROOT: %+v, %p`, tree.root.left, tree.root.left)

	results := tree.GetRange(newQuery(0, 3, 0, 3))

	checkEntries(t, results, p1, p2, p3)
}

func TestCopyOverwritesRoot(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(0, 0)

	tree := newTree(2, 1, p1)
	tree.Insert(p2)

	results := tree.GetRange(newQuery(0, 1, 0, 1))

	checkEntries(t, results, p2)
}

func TestCopyOverwritesBranch(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 0)

	tree := newTree(2, 1, p1, p2)

	p3 := newPoint(1, 0)
	tree.Insert(p3)

	results := tree.GetRange(newQuery(0, 3, 0, 3))

	checkEntries(t, results, p1, p3)
}

func TestOverwriteMultipleDimensions(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)
	p3 := newPoint(2, 2)
	p4 := newPoint(0, 0)
	p5 := newPoint(1, 1)
	p6 := newPoint(2, 2)

	tree := newTree(2, 1, p1, p2, p3)

	tree.Insert(p4, p5, p6)

	results := tree.GetRange(newQuery(0, 3, 0, 3))

	checkEntries(t, results, p4, p5, p6)
}

func TestAll(t *testing.T) {
	numItems := 3
	entries := make([]r.Entry, numItems)

	for i := 0; i < numItems; i++ {
		entries[i] = newPoint(i, i)
	}

	tree := New(2, entries...)

	results := tree.All()

	checkEntries(t, results, entries...)
}

func TestGetSingleItem(t *testing.T) {
	p1 := newPoint(0, 0)
	p2 := newPoint(1, 0)
	p3 := newPoint(0, 1)
	p4 := newPoint(1, 1)

	tree := newTree(2, 1, p1, p2, p3, p4)

	entries := tree.GetRange(newQuery(1, 2, 1, 2))

	checkEntries(t, entries, p4)
}

func TestQueryNonExistentValue(t *testing.T) {
	p1 := newPoint(0, 0)

	tree := newTree(2, 1, p1)

	entries := tree.GetRange(newQuery(1, 2, 1, 2))

	checkEntries(t, entries)
}

func BenchmarkGetRange(b *testing.B) {
	points := make([]r.Entry, 100)
	tree := newTree(2, 1)
	for i := 0; i < 100; i++ {
		points[i] = newPoint(i, i)
	}

	tree.Insert(points...)

	var results []r.Entry

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results = tree.GetRange(newQuery(0, 100, 0, 100))
	}

	b.StopTimer()

	if len(results) != 100 {
		b.Errorf(`Expected len: %d, received: %d`, 100, len(results))
	}
}

func BenchmarkInsertNodes(b *testing.B) {
	numItems := 100
	points := make([]r.Entry, numItems)

	for i := 0; i < numItems; i++ {
		points[i] = newPoint(i, i)
	}

	b.ResetTimer()

	tree := New(2)

	for i := 0; i < b.N; i++ {
		tree.Insert(points...)
	}
	b.StopTimer()

	result := tree.GetRange(newQuery(0, numItems, 0, numItems))

	if len(result) != numItems {
		b.Errorf(`Expected len: %d, received: %d`, numItems, len(result))
	}
}

func BenchmarkAll(b *testing.B) {
	numItems := 100
	points := make([]r.Entry, numItems)

	for i := 0; i < numItems; i++ {
		points[i] = newPoint(i, i)
	}

	tree := New(2, points...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.All()
	}
}

func TestDenseMatrixQuery(t *testing.T) {
	maxRange := 5

	tree := New(2)

	points := make([]r.Entry, maxRange*maxRange)
	index := 0

	for i := 0; i < maxRange; i++ {
		for j := 0; j < maxRange; j++ {
			if i%2 == 0 && j%2 == 0 {
				continue
			}

			p := newPoint(i, j)
			points[index] = p
			index++
		}
	}
	tz := time.Now()

	tree.Insert(points[:index]...)

	log.Printf(`INSERT TOOK: %d ms`, time.Since(tz).Nanoseconds()/int64(time.Millisecond))

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

func TestReturnParents(t *testing.T) {
	n1 := &node{value: 1}
	n2 := &node{value: 2}
	n1.parent = n2

	tree := New(2)

	path := tree.returnParents(n1)

	expected := []*node{n2, n1}
	if !reflect.DeepEqual(path, expected) {
		t.Errorf(`Expected: %+v, received: %+v`, expected, path)
	}
}
