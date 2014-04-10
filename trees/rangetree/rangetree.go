package rangetree

import (
	"sort"
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
	left *node
	right *node
	value int
}

type result struct {
	ints []int
}

func (self *node) isLeaf() bool {
	return self.left == nil
}

func (self *node) getRange(start, stop int, ints *result) {
	if self.isLeaf() {
		if self.value >= start && self.value < stop {
			ints.ints = append(ints.ints, self.value)
			return
		} else {
			return
		}
	}

	if start <= self.value {
		self.left.getRange(start, stop, ints)
	}
	if stop > self.value {
		self.right.getRange(start, stop, ints)
	}
}

func newNode(values []int) *node {
	if len(values) == 1 {
		return &node{
			value: values[0],
		}
	} else if len(values) == 2 {
		return &node{
			left: newNode(values[:1]),
			right: newNode(values[1:2]),
			value: values[1],
		}
	}

	median := len(values) / 2

	return &node{
		left: newNode(values[:median+1]),
		right: newNode(values[median+1:len(values)]),
		value: values[median],
	}
}

type tree struct {
	root *node
}

func (self *tree) GetRange(start, stop int) []int {
	ints := &result{
		ints: make([]int, 0),
	}
	self.root.getRange(start, stop, ints)
	return ints.ints
}

func New(values []int) *tree {
	sort.Ints(values)
	return &tree{
		root: newNode(values),
	}
}