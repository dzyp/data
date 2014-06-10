package v1

import (
	r "github.com/dzyp/data/trees/rangetree"
)

var (
	REBALANCE_RATIO float64 = .3 // performance tuning will be required to change this
	// .5 would be perfectly balanced
)

type node struct {
	left        *node
	right       *node
	parent      *node
	entry       r.Entry
	value       int
	numChildren int
	rt          *tree
}

func newNode(tree *tree, entries *entriesWrapper) *node {
	if entries.len() == 0 {
		return nil
	}

	if tree.isLastDimension() && entries.isLastValue() {
		entry := entries.lastValue()
		return &node{
			entry: entry,
			value: entry.GetDimensionalValue(tree.dimension),
		}
	}

	if entries.isLastValue() { // we need to add another tree
		return &node{
			value: entries.median(),
			rt: new(
				tree.maxDimensions,
				tree.dimension+1,
				entries.getEntriesAtValue(entries.median())...,
			),
		}
	}
	left, right := entries.split(-1)

	n := &node{
		left:        newNode(tree, left),
		right:       newNode(tree, right),
		value:       entries.median(),
		numChildren: entries.len(),
	}

	n.left.parent = n
	n.right.parent = n

	return n
}

type queryResult struct {
	entries []r.Entry
	index   int
}

func (self *queryResult) addEntry(entry r.Entry) {
	self.entries[self.index] = entry
	self.index++
}

func (self *queryResult) results() []r.Entry {
	return self.entries[0:self.index]
}

func newResult(numChildren int) *queryResult {
	return &queryResult{
		entries: make([]r.Entry, numChildren),
	}
}

func (self *node) isLeaf() bool {
	return self.left == nil
}

func (self *node) isLastDimension() bool {
	return self.rt == nil
}

func (self *node) needsRebalancing() bool {
	if self.isLeaf() {
		return false
	}

	total := float64(self.left.numChildren + self.right.numChildren)

	if float64(self.left.numChildren)/total < REBALANCE_RATIO {
		return true
	} else if float64(self.right.numChildren)/total < REBALANCE_RATIO {
		return true
	}

	return false
}

func (self *node) sibling() *node {
	if self.isRoot() {
		return nil
	}

	if self.isRight() {
		return self.parent.left
	} else {
		return self.parent.right
	}
}

func (self *node) all(results *queryResult) {
	if self.isLeaf() {
		if self.isLastDimension() {
			results.addEntry(self.entry)
		} else {
			self.rt.all(results)
		}

		return
	}

	self.left.all(results)
	self.right.all(results)
}

func (self *node) getRange(query r.Query, dimension int, results *queryResult, left, right bool) {
	bounds := query.GetDimensionalBounds(dimension)
	if self.isLeaf() {
		if self.value >= bounds.Low() && self.value < bounds.High() {
			if self.rt == nil { // i am a true leaf, last dimension
				results.addEntry(self.entry)
				return
			} else { // i am not the last dimension
				self.rt.getRange(query, results)
				return
			}
		} else { // we should hopefully not get here
			return
		}
	}

	if bounds.High() <= self.value {
		self.left.getRange(query, dimension, results, left, right) //left right should be false here
		return
	}

	if bounds.Low() > self.value {
		self.right.getRange(query, dimension, results, left, right) //left right should be false here
		return
	}

	if bounds.Low() <= self.value && left { // we can safely grab all of right here
		self.left.getRange(query, dimension, results, true, false)
		self.right.flatten(query, dimension, results)
	} else if bounds.High() > self.value && right {
		self.left.flatten(query, dimension, results)
		self.right.getRange(query, dimension, results, false, true)
	} else {
		self.left.getRange(query, dimension, results, true, false)
		self.right.getRange(query, dimension, results, false, true)
	}
}

func (self *node) grandParent() *node {
	if self.parent == nil {
		return nil
	}

	return self.parent.parent
}

func (self *node) flatten(query r.Query, dimension int, results *queryResult) {
	if self.isLeaf() {
		if self.rt == nil { // i am a true leaf
			results.addEntry(self.entry)
		} else {
			self.rt.getRange(query, results)
		}
		return
	}

	self.left.flatten(query, dimension, results)
	self.right.flatten(query, dimension, results)
}

func (self *node) rebalance(tree *tree) {
	if self.isLeaf() { // i can't be rebalanced
		if self.rt == nil { // i am last dimension
			return
		} else {
			self.rt.rebalance()
			return
		}
	}

	if self.needsRebalancing() {
		results := newResult(tree.numChildren)
		self.left.all(results)
		self.right.all(results)

		entries := results.results()

		n := newNode(tree, newEntries(entries, tree.dimension, false))
		n.parent = self.parent
		if self.isRoot() {
			tree.root = n
		} else {
			if self.isLeft() {
				self.parent.left = n
			} else {
				self.parent.right = n
			}
		}

		return // we don't need to rebalance our children now
	} else {
		self.left.rebalance(tree)
		self.right.rebalance(tree)
	}
}

/*
returns the inserted entry, returns nil if nothing was inserted
*/
func (self *node) insert(tree *tree, entries *entriesWrapper) {
	if len(entries.entries) == 0 {
		return
	}

	if self.isLeaf() {
		median := entries.median()

		if median == self.value && len(entries.sortedDimensionalValues) == 1 && !tree.isLastDimension() {
			self.rt.insert(entries.getEntriesAtValue(self.value)...)
		} else if median == self.value { // we now go to the right
			if !tree.isLastDimension() {
				self.rt.insert(entries.getEntriesAtValue(self.value)...) // insert entries
			}

			left, right := entries.split(-1)

			leftN := newNode(tree, left)
			if leftN == nil { // no values left, we need to take the new value
				if tree.isLastDimension() {
					self.entry = entries.lastValue()
				} else {
					self.rt.insert(entries.getEntriesAtValue(self.value)...)
				}

				return
			}
			leftN.parent = self
			self.left = leftN

			rightN := &node{
				value:  self.value,
				rt:     self.rt,
				parent: self,
			}

			rightN.entry = self.entry
			self.entry = nil

			self.right = rightN
			self.rt = nil
			self.right.insert(tree, right)
			self.numChildren = 2
		} else if median > self.value {
			leftN := &node{
				value:  self.value,
				rt:     self.rt,
				parent: self,
				entry:  self.entry,
			}

			self.entry = nil

			self.left = leftN
			self.rt = nil
			self.value = median

			left, right := entries.split(-1)

			rightN := newNode(tree, right)
			rightN.parent = self
			self.right = rightN
			self.left.insert(tree, left)
			self.numChildren = 2
		} else if median < self.value {
			rightN := &node{
				value:  self.value,
				rt:     self.rt,
				parent: self,
			}

			rightN.entry = self.entry
			self.entry = nil

			left, right := entries.split(
				len(entries.sortedDimensionalValues)/2 + 1,
			)

			self.right = rightN
			self.rt = nil
			leftN := newNode(tree, left)
			self.left = leftN
			leftN.parent = self
			self.right.insert(tree, right)
			self.numChildren = 2
		}

		return
	}

	index := entries.find(self.value)

	left, right := entries.split(index)

	self.left.insert(tree, left)
	self.right.insert(tree, right)
}

func (self *node) copy() *node {
	newNode := &node{
		numChildren: self.numChildren,
		value:       self.value,
	}

	if self.rt != nil {
		newNode.rt = self.rt.copy()
	}

	if self.isLeaf() {
		return newNode
	}

	newNode.left = self.left.copy()
	newNode.left.parent = newNode
	newNode.right = self.right.copy()
	newNode.right.parent = newNode
	return newNode
}

func (self *node) isRoot() bool {
	return self.parent == nil
}

func (self *node) isRight() bool {
	if self.isRoot() {
		return false
	}

	return self.parent.right == self
}

func (self *node) isLeft() bool {
	if self.isRoot() {
		return false
	}

	return self.parent.left == self
}

func (self *node) removeSelf(tree *tree) {
	if self.isRoot() { // we are the root
		tree.root = nil
		return
	}

	if self.grandParent() == nil { // map parent is the root
		sibling := self.sibling()
		tree.root = sibling
		sibling.parent = nil
		return
	}

	sibling := self.sibling()
	gp := self.grandParent()
	sibling.parent = gp
	if self.parent.isRight() {
		gp.right = sibling
	} else {
		gp.left = sibling
	}
}

/*
Returns nil if entry wasn't found
*/
func (self *node) remove(tree *tree, entry r.Entry) r.Entry {
	if self.isLeaf() {
		if self.value == entry.GetDimensionalValue(tree.dimension) {
			if self.rt == nil { // we are the last dimension
				self.removeSelf(tree)

				return entry
			} else {
				entry = self.rt.remove(entry)
				if entry == nil { // nothing was removed
					return entry
				}

				if self.rt.numChildren == 0 {
					self.removeSelf(tree)
				}

				return entry
			}
		} else {
			return nil
		}
	}

	if entry.GetDimensionalValue(tree.dimension) >= self.value {
		entry = self.right.remove(tree, entry)
	} else {
		entry = self.left.remove(tree, entry)
	}

	if entry != nil {
		self.numChildren--
	}

	return entry
}

type tree struct {
	root          *node
	dimension     int
	maxDimensions int
	numChildren   int
}

func (self *tree) remove(entry r.Entry) r.Entry {
	if self.root == nil {
		return nil
	}

	entry = self.root.remove(self, entry)
	if entry != nil {
		self.numChildren--
	}

	return entry
}

func (self *tree) rebalance() {
	if self.root == nil {
		return
	}

	self.root.rebalance(self)
}

func (self *tree) isLastDimension() bool {
	return self.dimension >= self.maxDimensions
}

func (self *tree) all(results *queryResult) {
	if self.root == nil {
		return
	}

	self.root.all(results)
}

func (self *tree) All() []r.Entry {
	results := newResult(self.numChildren)
	self.all(results)

	return results.entries[0:results.index]
}

func (self *tree) Remove(entries ...r.Entry) {
	for _, entry := range entries {
		self.remove(entry)
	}
}

func (self *tree) Len() int {
	return self.numChildren
}

func (self *tree) getRange(query r.Query, results *queryResult) {
	if self.root == nil {
		return
	}

	self.root.getRange(query, self.dimension, results, false, false)
}

func (self *tree) GetRange(query r.Query) []r.Entry {
	if self.root == nil {
		return []r.Entry{}
	}

	results := newResult(self.numChildren)

	self.getRange(query, results)
	return results.entries[0:results.index]
}

func (self *tree) insert(entries ...r.Entry) r.Entry {
	ew := newEntries(entries, self.dimension, false)
	if self.root == nil {
		self.root = newNode(self, ew)
		self.numChildren += len(entries)
		return nil
	}

	self.root.insert(self, ew)

	self.numChildren += len(entries)

	return nil
}

func (self *tree) Insert(values ...r.Entry) {
	byDimension(self.dimension).Sort(values)

	self.insert(values...)
}

func (self *tree) copy() *tree {
	cp := &tree{
		dimension:     self.dimension,
		maxDimensions: self.maxDimensions,
	}

	if self.root == nil {
		return cp
	}

	cp.root = self.root.copy()
	return cp
}

func (self *tree) Copy() r.RangeTree {
	return self.copy()
}

func (self *tree) Clear() {
	self.root = nil
}

func new(maxDimensions, dimension int, entries ...r.Entry) *tree {
	t := &tree{
		maxDimensions: maxDimensions,
		dimension:     dimension,
		numChildren:   len(entries),
	}

	t.root = newNode(t, newEntries(entries, dimension, false))
	return t
}

func New(maxDimensions int, entries ...r.Entry) *tree {
	byDimension(1).Sort(entries)
	return new(maxDimensions, 1, entries...)
}
