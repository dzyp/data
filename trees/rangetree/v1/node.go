package v1

import (
	r "github.com/dzyp/data/trees/rangetree"
)

type itree interface {
	// takes a list of entries and returns a value indicating the number of entries
	// added
	insert(entries ...r.Entry) int
	copy() itree
	query(r.Query, *result)
}

type node struct {
	parent      *node
	left        *node
	right       *node
	value       int
	p           itree
	numChildren int
}

func newNode() *node {
	return &node{}
}

func (self *node) isLeaf() bool {
	return self.left == nil
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

func (self *node) splitLeft(tree *tree, count *int, value int, leftNodes, rightNodes []r.Entry) {
	left := self.copy()
	left.parent = self
	self.left = left
	right := newNode()
	right.value = value
	if tree.isSecondToLastDimension() {
		right.p = newOrderedList(tree.dimension + 1)
	} else {
		right.p = newTree(tree.maxDimensions, tree.dimension+1)
	}
	self.right = right
	right.parent = self

	self.left.insert(tree, count, leftNodes...)
	self.right.insert(tree, count, rightNodes...)

	self.numChildren += *count
}

func (self *node) splitRight(tree *tree, count *int, value int, leftNodes, rightNodes []r.Entry) {
	right := self.copy()
	right.parent = self
	self.right = right
	left := newNode()
	left.value = value
	left.parent = self
	self.left = left

	if tree.isSecondToLastDimension() {
		left.p = newOrderedList(tree.dimension + 1)
	} else {
		left.p = newTree(tree.maxDimensions, tree.dimension+1)
	}

	self.left.insert(tree, count, leftNodes...)
	self.right.insert(tree, count, rightNodes...)

	self.value = value
	self.numChildren += *count
}

func (self *node) query(tree *tree, query r.Query, result *result, left, right bool) {
	bounds := query.GetDimensionalBounds(tree.dimension)
	if self.isLeaf() {
		self.p.query(query, result)
		return
	}

	if bounds.High() <= self.value {
		self.left.query(tree, query, result, left, right) //left right should be false here
		return
	}

	if bounds.Low() > self.value {
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

func (self *node) insert(tree *tree, count *int, entries ...r.Entry) {
	if len(entries) == 0 {
		return
	}

	var value int

	if self.isLeaf() {
		leftNodes, rightNodes := Entries(entries).SplitAtValue(
			self.value, tree.dimension,
		)

		if len(leftNodes) > 0 { // there are more than 1 value and we need to split
			value = Entries(leftNodes).MedianEntry().GetDimensionalValue(
				tree.dimension,
			)

			self.splitRight(tree, count, value, leftNodes, rightNodes)
			self.p.insert(entries...)
		} else {
			index := 0
			nodesToInsert := make([]r.Entry, 0, len(rightNodes))
			for _, node := range rightNodes {
				if node.GetDimensionalValue(tree.dimension) == self.value {
					nodesToInsert = append(nodesToInsert, node)
					index++
				} else {
					break
				}
			}

			if index > 0 {
				added := self.p.insert(nodesToInsert...)
				*count += added
				self.numChildren += added
				_, rightNodes = Entries(rightNodes).SplitAtIndex(index)
			}

			if len(rightNodes) == 0 {
				return
			}

			value = Entries(rightNodes).MedianEntry().GetDimensionalValue(
				tree.dimension,
			)

			self.splitLeft(tree, count, value, leftNodes, rightNodes)
			self.p.insert(rightNodes...)
		}

		return // all leaves work are done
	}

	left, right := Entries(entries).SplitAtValue(self.value, tree.dimension)
	self.left.insert(tree, count, left...)
	self.right.insert(tree, count, right...)

	self.numChildren += *count
	self.p.insert(entries...)
}
