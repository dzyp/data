package rangetree

import (
	"sort"
	"sync"
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
	left  *node
	right *node
	value int
}

type result struct {
	ints []int
}

func newResult() *result {
	return &result{
		ints: make([]int, 0),
	}
}

func (self *node) isLeaf() bool {
	return self.left == nil
}

func (self *node) getRange(start, stop int, ints *result, left, right bool) {
	if self.isLeaf() {
		if self.value >= start && self.value < stop {
			ints.ints = append(ints.ints, self.value)
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
		res := newResult()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			self.right.flatten(res)
			println(`go routine done`)
			wg.Done()
		}()
		self.left.getRange(start, stop, ints, true, false)
		wg.Wait()

		ints.ints = append(ints.ints, res.ints...)
	} else if stop > self.value && right {
		res := newResult()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			self.left.flatten(res)
			wg.Done()
		}()
		self.right.getRange(start, stop, ints, false, true)
		wg.Wait()

		ints.ints = append(ints.ints, res.ints...)
	} else {
		self.left.getRange(start, stop, ints, true, false)
		self.right.getRange(start, stop, ints, false, true)
	}
}

func (self *node) flatten(ints *result) {
	if self.isLeaf() {
		ints.ints = append(ints.ints, self.value)
		return
	}

	self.left.flatten(ints)
	self.right.flatten(ints)
}

func newNode(values []int) *node {
	if len(values) == 1 {
		return &node{
			value: values[0],
		}
	} else if len(values) == 2 {
		return &node{
			left:  newNode(values[:1]),
			right: newNode(values[1:2]),
			value: values[1],
		}
	}

	median := len(values) / 2

	return &node{
		left:  newNode(values[:median+1]),
		right: newNode(values[median+1 : len(values)]),
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
	self.root.getRange(start, stop, ints, false, false)
	return ints.ints
}

func New(values []int) *tree {
	sort.Ints(values)
	return &tree{
		root: newNode(values),
	}
}
