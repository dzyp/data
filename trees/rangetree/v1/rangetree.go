package v1

import (
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
	if len(entries) == 0 {
		return 0
	}

	if self.root == nil {
		self.root = newNodesFromEntries(
			self,
			Entries(entries).GetOrderedUniqueAtDimension(self.dimension),
			entries,
		)
		self.numChildren = self.root.numChildren
		return self.numChildren
	}

	var count int

	self.root.insert(
		self,
		&count,
		Entries(entries).GetOrderedUniqueAtDimension(self.dimension),
		entries,
	)
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

func (self *tree) all(result *result) {
	if self.root == nil {
		return
	}

	self.root.all(result)
}

func (self *tree) All() []r.Entry {
	results := newResults(self.numChildren)
	self.all(results)
	return results.entries
}

func (self *tree) Remove(entries ...r.Entry) {

}

func (self *tree) Copy() r.RangeTree {
	return self.copy().(*tree)
}

func (self *tree) Clear() {
	self.root = nil
	self.numChildren = 0
}

func (self *tree) Len() int {
	return self.numChildren
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
