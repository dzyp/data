package rangetree

import (
	"sort"
	//"sync"
)

var (
	REBALANCE_RATIO float64 = .3 // performance tuning will be required to change this
)

type point struct {
	x int
	y int
}

func newPoint(x, y int) *point {
	return &point{
		x: x,
		y: y,
	}
}

type node struct {
	left        *node
	right       *node
	parent      *node
	value       int
	numChildren int
}

type result struct {
	ints  []int
	index int
}

func newResult(length int) *result {
	return &result{
		ints: make([]int, length),
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

func (self *node) getRange(start, stop int, ints *result, left, right bool) {
	if self.isLeaf() {
		if self.value >= start && self.value < stop {
			ints.ints[ints.index] = self.value
			ints.index++
			return
		} else {
			return
		}
	}

	if stop <= self.value {
		self.left.getRange(start, stop, ints, left, right) //left right should be false here
		return
	}

	if start > self.value {
		self.right.getRange(start, stop, ints, left, right) //left right should be false here
		return
	}

	if start <= self.value && left { // we can safely grab all of right here
		self.left.getRange(start, stop, ints, true, false)
		self.right.flatten(ints)
	} else if stop > self.value && right {
		self.left.flatten(ints)
		self.right.getRange(start, stop, ints, false, true)
	} else {
		self.left.getRange(start, stop, ints, true, false)
		self.right.getRange(start, stop, ints, false, true)
	}
}

func (self *node) grandParent() *node {
	if self.parent == nil {
		return nil
	}

	return self.parent.parent
}

func (self *node) flatten(ints *result) {
	if self.isLeaf() {
		ints.ints[ints.index] = self.value
		ints.index++
		return
	}

	self.left.flatten(ints)
	self.right.flatten(ints)
}

func (self *node) insert(value int) {
	if value == self.value {
		return
	}
}

func (self *node) getLeaf(value int) *node {
	if self.isLeaf() {
		return self
	}

	if value <= self.value {
		return self.left.getLeaf(value)
	} else {
		return self.right.getLeaf(value)
	}
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

func (self *node) incrementParents(delta int) {
	n := self.parent
	for n != nil {
		n.numChildren += delta
		n = n.parent
	}
}

func (self *node) needsRebalancing() bool {
	if self.isLeaf() {
		return false
	}

	return float64(self.left.numChildren)/float64(self.numChildren) < REBALANCE_RATIO ||
		float64(self.right.numChildren)/float64(self.numChildren) < REBALANCE_RATIO
}

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

func newNode(values []int, parent *node) *node {
	if len(values) == 1 {
		return &node{
			parent: parent,
			value:  values[0],
		}
	} else if len(values) == 2 {
		n := &node{
			value:       values[0],
			numChildren: len(values),
			parent:      parent,
		}
		n.left = newNode(values[:1], n)
		n.right = newNode(values[1:2], n)
		return n
	}

	median := len(values) / 2

	n := &node{
		value:       values[median],
		numChildren: len(values),
		parent:      parent,
	}
	n.left = newNode(values[:median+1], n)
	n.right = newNode(values[median+1:len(values)], n)
	return n
}

type tree struct {
	root *node
}

func (self *tree) GetRange(start, stop int) []int {
	ints := &result{
		ints: make([]int, stop-start),
	}
	self.root.getRange(start, stop, ints, false, false)
	return ints.ints[0:ints.index]
}

func (self *tree) Insert(value int) {
	n := self.root.getLeaf(value)
	if n.value == value { // we don't need to insert a duplicate
		return
	}

	oldValue := n.value
	n.value = value
	n.left = &node{
		value:  value,
		parent: n,
	}
	n.right = &node{
		value:  oldValue,
		parent: n,
	}

	n.incrementParents(1)
	n.rebalance(self)
}

func (self *tree) Delete(value int) {
	if self.root.isLeaf() && self.root.value == value {
		self.root = nil
		return
	} else if self.root.isLeaf() {
		return
	}

	n := self.root.getLeaf(value)
	if n.value != value { //value doesn't exist in the tree
		return
	}

	if n.parent.isRoot() { //special case to handle this
		sibling := n.sibling()
		sibling.parent = nil
		sibling.left = nil
		sibling.right = nil
		self.root = sibling
		return
	}

	sibling := n.sibling()
	sibling.parent = n.grandParent()
	if n.parent.isRight() {
		n.grandParent().right = sibling
	} else {
		n.grandParent().left = sibling
	}

	sibling.incrementParents(-1)
	sibling.rebalance(self)
}

/*
Retrieves the leaf matching the value *OR THE CLOSEST* approximation.
For inserts, this will be the left where you will need to insert.
*/
func (self *tree) GetLeaf(value int) *node {
	return self.root.getLeaf(value)
}

func New(values []int) *tree {
	sort.Ints(values)
	return &tree{
		root: newNode(values, nil),
	}
}
