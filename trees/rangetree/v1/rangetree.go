package v1

import (
	r "github.com/dzyp/data/trees/rangetree"
)

var (
	REBALANCE_RATIO float64 = .3 // performance tuning will be required to change this
)

type node struct {
	left        *node
	right       *node
	parent      *node
	value       r.Entry
	numChildren int
	rt          *tree
}

type queryResult struct {
	entries []r.Entry
	index   int
}

func (self *queryResult) addEntry(entry r.Entry) {
	self.entries[self.index] = entry
	self.index++
}

func newResult(numChildren int) *queryResult {
	return &queryResult{
		entries: make([]r.Entry, numChildren),
	}
}

func (self *node) isLeaf() bool {
	return self.left == nil
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

func (self *node) getRange(query r.Query, dimension int, results *queryResult, left, right bool) {
	bounds := query.GetDimensionalBounds(dimension)
	value := self.value.GetDimensionalValue(dimension)
	if self.isLeaf() {
		if value >= bounds.Low() && value < bounds.High() {
			if self.rt == nil { // i am a true leaf, last dimension
				results.addEntry(self.value)
				return
			} else { // i am not the last dimension
				self.rt.getRange(query, results)
				return
			}
		} else { // we should hopefully not get here
			return
		}
	}

	if bounds.High() <= value {
		self.left.getRange(query, dimension, results, left, right) //left right should be false here
		return
	}

	if bounds.Low() > value {
		self.right.getRange(query, dimension, results, left, right) //left right should be false here
		return
	}

	if bounds.Low() <= value && left { // we can safely grab all of right here
		self.left.getRange(query, dimension, results, true, false)
		self.right.flatten(query, dimension, results)
	} else if bounds.High() > value && right {
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
			results.addEntry(self.value)
		} else {
			self.rt.getRange(query, results)
		}
		return
	}

	self.left.flatten(query, dimension, results)
	self.right.flatten(query, dimension, results)
}

/*
returns the inserted entry, returns nil if nothing was inserted
*/
func (self *node) insert(entry r.Entry, dimension int) r.Entry {
	if self.isLeaf() {
		highestDimension := false
		if self.rt == nil {
			highestDimension = true
		}
		if self.value.LessThan(entry, dimension) {
			newLeftNode := &node{}
			newRightNode := &node{}
			newRightNode.value = entry
			newLeftNode.value = self.value
			newLeftNode.rt = self.rt
			self.rt = nil
			newLeftNode.parent = self
			newRightNode.parent = self
			self.left = newLeftNode
			self.right = newRightNode
			self.value = entry
			self.numChildren = 2
			if !highestDimension {
				self.right.rt = new(self.left.rt.maxDimensions, dimension)
				return self.right.rt.insert(entry)
			}

			return entry
		} else if entry.EqualAtDimension(self.value, dimension) {
			if self.rt == nil { // duplicate :(
				return nil
			}

			return self.rt.insert(entry)
		} else {
			newLeftNode := &node{}
			newRightNode := &node{}
			newRightNode.value = self.value
			newRightNode.rt = self.rt
			self.rt = nil
			newLeftNode.parent = self
			newRightNode.parent = self
			newLeftNode.value = entry
			self.left = newLeftNode
			self.right = newRightNode
			self.numChildren = 2
			if !highestDimension {
				self.left.rt = new(self.right.rt.maxDimensions, dimension)
				return self.left.rt.insert(entry)
			}
			return entry
		}

		return nil
	}

	var newEntry r.Entry

	if entry.LessThan(self.value, dimension) {
		newEntry = self.left.insert(entry, dimension)
	} else {
		newEntry = self.right.insert(entry, dimension)
	}

	if newEntry != nil {
		self.numChildren++
	}

	return newEntry
}

func (self *node) copy() *node {
	newNode := &node{
		numChildren: self.numChildren,
		value:       self.value,
	}

	if self.rt != nil {
		newNode.rt = self.rt.Copy()
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
		if self.value.EqualAtDimension(entry, tree.dimension) {
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

	if entry.LessThan(self.value, tree.dimension) {
		entry = self.left.remove(tree, entry)
	} else {
		entry = self.right.remove(tree, entry)
	}

	if entry != nil {
		self.numChildren--
	}

	return entry
}

/*
func (self *node) needsRebalancing() bool {
	if self.isLeaf() {
		return false
	}

	return float64(self.left.numChildren)/float64(self.numChildren) < REBALANCE_RATIO ||
		float64(self.right.numChildren)/float64(self.numChildren) < REBALANCE_RATIO
}
*/
/*
func (self *node) rebalance(tree *tree) {
	var n *node
	if self.isLeaf() {
		if self.parent == nil {
			return
		}

		n = self.parent
	} else {
		n = self
	}
	for n != nil {
		if n.parent != nil && n.parent.needsRebalancing() { // we don't want to duplicate rebalances
			n = n.parent
			continue
		} else if n.needsRebalancing() {
			result := newResult(n.numChildren)
			n.flatten(result)
			newNode := newNode(result.ints, n.parent)
			if n.parent == nil { //rebalanced at root
				tree.root = newNode
			} else {
				if n.parent.left == n {
					n.parent.left = newNode
				} else {
					n.parent.right = newNode
				}
			}
			return
		} else {
			n = n.parent
		}
	}
}
*/

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

func (self *tree) Remove(entries ...r.Entry) {
	for _, entry := range entries {
		self.remove(entry)
	}
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

func (self *tree) insert(value r.Entry) r.Entry {
	if self.root == nil {
		self.root = &node{
			value: value,
		}

		self.numChildren = 1

		if self.dimension < value.MaxDimensions() {
			self.root.rt = new(self.maxDimensions, self.dimension+1)
			self.root.rt.insert(value)
		}

		return value
	}

	value = self.root.insert(value, self.dimension)
	if value != nil {
		self.numChildren++
	}

	return value
}

func (self *tree) Insert(values ...r.Entry) {
	for _, value := range values {
		self.insert(value)
	}
}

func (self *tree) Copy() *tree {
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

func (self *tree) Clear() {
	self.root = nil
}

func new(maxDimensions, dimension int) *tree {
	return &tree{
		maxDimensions: maxDimensions,
		dimension:     dimension,
	}
}

func New(maxDimensions int) *tree {
	return &tree{
		dimension:     1,
		maxDimensions: maxDimensions,
	}
}
