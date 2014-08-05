package v1

import (
	"runtime"
	"sync"

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

func (self *tree) insert(items ...r.Entry) {
	if len(items) == 0 {
		return
	}

	entries := make([]r.Entry, len(items))
	copy(entries, items)

	if self.root == nil {
		i := len(entries) / 2
		entry := entries[i]
		root := newNode()
		root.value = entry.GetDimensionalValue(self.dimension)
		if self.isSecondToLastDimension() {
			root.p = newOrderedList(self.dimension + 1)
			root.p.insert(entry)
		} else {
			root.p = newTree(self.maxDimensions, self.dimension+1, entry)
		}

		copy(entries[i:], entries[i+1:])
		entries[len(entries)-1] = nil // or the zero value of T
		entries = entries[:len(entries)-1]

		self.root = root
	}

	path := make([]*node, 0, len(entries))
	nextDimensionMap := make(map[*node][]r.Entry, len(entries))

	for _, entry := range entries {
		parent := self.root
		path = nil
		value := entry.GetDimensionalValue(self.dimension)
		for {
			if parent.isLeaf() {
				if value == parent.value { // add to next dimension
					path = append(path, parent)
				} else if value > parent.value {
					path = append(path, parent)
					leftNode := newNode()
					leftNode.p = parent.p.copy()
					leftNode.value = parent.value
					leftNode.parent = parent
					parent.left = leftNode

					toAdd := nextDimensionMap[parent]
					if toAdd != nil {
						ns := make([]r.Entry, len(toAdd))
						copy(ns, toAdd)
						nextDimensionMap[leftNode] = ns
					}

					rightNode := newNode()
					if self.isSecondToLastDimension() {
						rightNode.p = newOrderedList(self.dimension + 1)
					} else {
						rightNode.p = newTree(
							self.maxDimensions, self.dimension+1,
						)
					}
					rightNode.value = value
					rightNode.parent = parent
					parent.right = rightNode
					path = append(path, rightNode)

					parent.value = value
				} else {
					path = append(path, parent)
					rightNode := newNode()
					rightNode.value = parent.value
					rightNode.parent = parent
					rightNode.p = parent.p.copy()
					parent.right = rightNode

					toAdd := nextDimensionMap[parent]
					if toAdd != nil {
						ns := make([]r.Entry, len(toAdd))
						copy(ns, toAdd)
						nextDimensionMap[rightNode] = ns
					}

					leftNode := newNode()
					if self.isSecondToLastDimension() {
						leftNode.p = newOrderedList(self.dimension + 1)
					} else {
						leftNode.p = newTree(
							self.maxDimensions, self.dimension+1,
						)
					}
					leftNode.value = value
					leftNode.parent = parent
					parent.left = leftNode
					path = append(path, leftNode)
				}

				break
			}

			path = append(path, parent)
			if value < parent.value {
				parent = parent.left
			} else {
				parent = parent.right
			}
		}

		for _, node := range path {
			if _, ok := nextDimensionMap[node]; !ok {
				nextDimensionMap[node] = make([]r.Entry, 0, len(entries))
			}
			nextDimensionMap[node] = append(nextDimensionMap[node], entry)
		}
	}

	path = make([]*node, 0, len(nextDimensionMap))

	for node, _ := range nextDimensionMap {
		path = append(path, node)
	}

	chunks := splitNodes(path, runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for _, chunk := range chunks {
		go func(nodes []*node) {
			for _, node := range nodes {
				entries := nextDimensionMap[node]
				node.p.insert(entries...)
			}

			wg.Done()
		}(chunk)
	}

	wg.Wait()
}

func (self *tree) Insert(entries ...r.Entry) {
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

func (self *tree) len() int {
	if self.root == nil {
		return 0
	}

	return self.root.len()
}

func (self *tree) All() []r.Entry {
	results := newResults(self.Len())
	self.all(results)
	return results.entries
}

//returnParents will take a node and return a list representing the path
//to the node with the node coming last
func (self *tree) returnParents(n *node) []*node {
	nodes := make([]*node, 0, self.Len())
	for n != nil {
		nodes = append(nodes, n)
		n = n.parent
	}

	nodes = Nodes(nodes).Reverse()

	return nodes
}

//findPath will return a slice representing the path to that node.  The first
//item in the slice will be closest to the root and the last item closest to
//the node.  Approximate indicates if this function should return a leaf
//node "closest" to the given value or should return only an exact match.
func (self *tree) findPath(entry r.Entry, approximate bool) []*node {
	if self.root == nil {
		return nil
	}

	parent := self.root
	for {
		if parent.isLeaf() {

		}
	}
}

func (self *tree) remove(entries ...r.Entry) {

}

func (self *tree) Remove(entries ...r.Entry) {
	if self.root == nil {
		return
	}
}

func (self *tree) Copy() r.RangeTree {
	return self.copy().(*tree)
}

func (self *tree) Clear() {
	self.root = nil
}

func (self *tree) Len() int {
	return self.len()
}

func (self *tree) GetRange(query r.Query) []r.Entry {
	result := newResults(self.Len())
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
