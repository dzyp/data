package v1

import (
	r "github.com/dzyp/data/trees/rangetree"
)

type itree interface {
	// takes a list of entries and returns a value indicating the number of entries
	// added
	insert(entries ...r.Entry)
	remove(entries ...r.Entry)
	copy() itree
	query(r.Query, *result)
	all(*result)
	len() int
}

type node struct {
	parent *node
	left   *node
	right  *node
	value  int
	p      itree
}

func newNode() *node {
	return &node{}
}

func newNodesFromEntries(tree *tree, values []int, entries []r.Entry) *node {
	if len(values) == 0 {
		return nil
	}

	n := newNode()

	if tree.isSecondToLastDimension() {
		n.p = newOrderedList(tree.dimension + 1)
	} else {
		n.p = newTree(tree.maxDimensions, tree.dimension+1)
	}

	if len(values) == 1 {
		n.value = values[0]
		entries = Entries(entries).GetEntriesAtValue(n.value, tree.dimension)
		n.p.insert(entries...)
		return n
	}

	value := Values(values).ValueAtMedian()
	leftValues, rightValues := Values(values).SplitAtMedian()

	leftEntries, rightEntries := Entries(entries).SplitAtValue(
		value, tree.dimension,
	)

	n.value = value

	left := newNodesFromEntries(tree, leftValues, leftEntries)
	right := newNodesFromEntries(tree, rightValues, rightEntries)
	left.parent = n
	right.parent = n

	n.left = left
	n.right = right
	n.p.insert(entries...)

	return n
}

func (self *node) isLeaf() bool {
	return self.left == nil
}

func (self *node) isLeft() bool {
	if self.parent == nil {
		return false
	}

	return self.parent.left == self
}

func (self *node) isRoot() bool {
	return self.parent == nil
}

func (self *node) len() int {
	return self.p.len()
}

func (self *node) copy() *node {
	cp := newNode()
	cp.value = self.value
	cp.p = self.p.copy()

	if !self.isLeaf() {
		left := self.left.copy()
		right := self.right.copy()
		cp.left = left
		cp.right = right
		right.parent = cp
		left.parent = cp
	}

	return cp
}

func (self *node) query(tree *tree, query r.Query, result *result, left, right bool) {
	bounds := query.GetDimensionalBounds(tree.dimension)
	if self.isLeaf() {
		if self.value >= bounds.Low() && self.value < bounds.High() {
			self.p.query(query, result)
		}
		return
	}

	if bounds.High() < self.value {
		self.left.query(tree, query, result, left, right) //left right should be false here
		return
	}

	if bounds.Low() >= self.value {
		self.right.query(tree, query, result, left, right) //left right should be false here
		return
	}

	if bounds.Low() <= self.value && left { // we can safely grab all of right here
		self.left.query(tree, query, result, true, false)
		self.right.flatten(tree, query, result)
	} else if bounds.High() > self.value && right {
		self.left.flatten(tree, query, result)
		self.right.query(tree, query, result, false, true)
	} else {
		self.left.query(tree, query, result, true, false)
		self.right.query(tree, query, result, false, true)
	}
}

func (self *node) flatten(tree *tree, query r.Query, result *result) {
	self.p.query(query, result)
}

func (self *node) all(result *result) {
	self.p.all(result)
}

type Nodes []*node

func (self Nodes) Reverse() Nodes {
	for i, j := 0, len(self)-1; i < j; i, j = i+1, j-1 {
		self[i], self[j] = self[j], self[i]
	}

	self = self
	return self
}
