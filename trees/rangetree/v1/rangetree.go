package v1

import (
	"log"

	r "github.com/dzyp/data/trees/rangetree"
)

var (
	REBALANCE_RATIO float64 = .3 // performance tuning will be required to change this
	// .5 would be perfectly balanced
)

type result struct {
	entries []r.Entry
}

func (self *result) AddEntry(entries ...r.Entry) {
	self.entries = append(self.entries, entries...)
}

func newResults(num int) *result {
	return &result{
		entries: make([]r.Entry, 0, num),
	}
}

type tree struct {
	root          *node
	dimension     int
	maxDimensions int
	numChildren   int
}

func (self *tree) copy() itree {
	t := &tree{
		dimension:     self.dimension,
		maxDimensions: self.maxDimensions,
	}

	if self.root != nil {
		t.root = self.root.copy()
	}

	return t
}

func (self *tree) insert(entries ...r.Entry) int {
	var n *node
	if self.root == nil {
		n = newNode()
		n.value = Entries(entries).MedianEntry().GetDimensionalValue(
			self.dimension,
		)
		if self.isSecondToLastDimension() {
			n.p = newOrderedList(self.dimension + 1)
		} else {
			n.p = newTree(self.maxDimensions, self.dimension+1)
		}
		self.root = n
	} else {
		n = self.root
	}

	log.Printf(`ROOT VALUE: %+v`, n.value)

	count := 0

	self.root.insert(self, &count, entries...)
	self.numChildren += count

	return count
}

func (self *tree) Insert(entries ...r.Entry) {
	Entries(entries).Sort(self.dimension)
	self.insert(entries...)
}

func (self *tree) isSecondToLastDimension() bool {
	return self.dimension == self.maxDimensions-1
}

func (self *tree) query(query r.Query, result *result) {
	if self.root == nil {
		return
	}

	self.root.query(self, query, result, false, false)
}

func (self *tree) GetRange(query r.Query) []r.Entry {
	result := newResults(self.numChildren)
	self.query(query, result)
	return result.entries
}

func newTree(maxDimensions, dimension int, entries ...r.Entry) *tree {
	t := &tree{
		maxDimensions: maxDimensions,
		dimension:     dimension,
	}

	if len(entries) == 0 {
		return t
	}

	t.Insert(entries...)

	return t
}

func New(maxDimensions int, entries ...r.Entry) *tree {
	return newTree(maxDimensions, 1, entries...)
}
