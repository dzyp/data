package v1

import (
	"sort"

	r "github.com/dzyp/data/trees/rangetree"
)

type Entries []r.Entry

func (self Entries) MedianEntry() r.Entry {
	return self[self.Median()]
}

func (self Entries) Median() int {
	return len(self) / 2
}

func (self Entries) Sort(dimension int) []r.Entry {
	es := &entrySorter{
		entries:   self,
		dimension: dimension,
	}

	sort.Sort(es)
	return self
}

func (self Entries) Index(value, dimension int) int {
	return sort.Search(len(self), func(i int) bool {
		return self[i].GetDimensionalValue(dimension) >= value
	})
}

func (self Entries) SplitAtValue(value, dimension int) ([]r.Entry, []r.Entry) {
	i := self.Index(value, dimension)

	return self.SplitAtIndex(i)
}

func (self Entries) SplitAtMedian() ([]r.Entry, []r.Entry) {
	return self.SplitAtIndex(self.Median())
}

func (self Entries) SplitAtIndex(i int) ([]r.Entry, []r.Entry) {
	left := make([]r.Entry, i)
	right := make([]r.Entry, len(self)-i)

	copy(left, self[0:i])
	copy(right, self[i:len(self)])

	return left, right
}

func (self Entries) GetEntriesAtValue(value, dimension int) []r.Entry {
	i := self.Index(value, dimension)

	entries := make([]r.Entry, 0, len(self))

	for i := i; i < len(self); i++ {
		if self[i].GetDimensionalValue(dimension) == value {
			entries = append(entries, self[i])
		} else {
			break
		}
	}

	return entries
}

func (self Entries) RemoveAt(i int) []r.Entry {
	if i >= len(self) { // this can't happen
		return self
	}

	copy(self[i:], self[i+1:])
	self[len(self)-1] = nil
	self = self[:len(self)-1]

	return self
}

func (self Entries) ShiftLeft(value, dimension int, entries []r.Entry) {
	for _, entry := range entries {
		if entry.GetDimensionalValue(dimension) == value {
			self = append(self, entry)
		} else {
			break
		}
	}
}

/*
The return value is the number of new values entered
*/
func (self Entries) Merge(entries ...r.Entry) ([]r.Entry, int) {
	if len(entries) == 0 {
		return self, 0
	}

	nodesAdded := 0

	i := sort.Search(len(self), func(i int) bool {
		return !self[i].Less(entries[0], 1)
	})

	if len(entries) == 1 {
		if i == len(self) {
			self = append(self, entries[0])
			return self, 1
		}
		entry := entries[0]
		if !entry.Less(self[i], 1) && !self[i].Less(entry, 1) {
			self[i] = entry
			return self, 0
		}
		self = append(self, nil)
		copy(self[i+1:], self[i:])
		self[i] = entries[0]
		return self, 1
	}

	newNodes := make([]r.Entry, 0, len(entries)+len(self))

	newNodes = append(newNodes, self[:i]...)

	var newNode r.Entry
	var oldNode r.Entry
	var newIndex int
	var oldIndex int = i

	for {
		if len(self) == oldIndex && len(entries) == newIndex {
			break
		} else if len(self) == oldIndex {
			newNodes = append(newNodes, entries...)
			nodesAdded += len(entries)
			break
		} else if len(entries) == newIndex {
			newNodes = append(newNodes, self...)
			break
		}

		newNode = entries[newIndex]
		oldNode = self[oldIndex]

		if newNode.Less(oldNode, 1) {
			newNodes = append(newNodes, newNode)
			newIndex++
			nodesAdded++
		} else if oldNode.Less(newNode, 1) {
			newNodes = append(newNodes, oldNode)
			oldIndex++
		} else { //equal
			newNodes = append(newNodes, newNode) // we override the old value
			oldIndex++
			newIndex++
			// we don't add to nodes added here as this number of rebalancing purposes
		}
	}

	self = newNodes

	return newNodes, nodesAdded
}

/*
First value is items removed and second value is result slice
*/
func (self Entries) RemovePrefix(value, dimension int) ([]r.Entry, []r.Entry) {
	index := 0
	removed := make([]r.Entry, 0, len(self))
	for _, entry := range self {
		if entry.GetDimensionalValue(dimension) == value {
			removed = append(removed, entry)
			index++
		} else {
			break
		}
	}

	return removed, self[:index]
}

func (self Entries) GetOrderedUniqueAtDimension(dimension int) []int {
	if len(self) == 0 {
		return []int{}
	}

	uniques := make([]int, 0, len(self))
	index := 0
	uniques = append(uniques, self[0].GetDimensionalValue(dimension))

	for i := 1; i < len(self); i++ {
		if self[i].GetDimensionalValue(dimension) != uniques[index] {
			uniques = append(uniques, self[i].GetDimensionalValue(dimension))
			index++
		}
	}

	return uniques
}

type Values []int

func (self Values) Median() int {
	return len(self) / 2
}

func (self Values) ValueAtMedian() int {
	return self[self.Median()]
}

func (self Values) Index(value int) int {
	return sort.SearchInts(self, value)
}

func (self Values) SplitAtIndex(i int) ([]int, []int) {
	return self[:i], self[i:]
}

func (self Values) SplitAtMedian() ([]int, []int) {
	return self.SplitAtIndex(self.Median())
}

func (self Values) SplitAtValue(value int) ([]int, []int) {
	return self.SplitAtIndex(self.Index(value))
}

func (self Values) Add(value int) []int {
	i := self.Index(value)
	if i >= len(self) {
		self = append(self, value)
		return self
	}

	if self[i] == value {
		return self
	}

	self = append(self, 0)
	copy(self[i+1:], self[i:])
	self[i] = value

	return self
}

type entrySorter struct {
	entries   []r.Entry
	dimension int
}

func (self *entrySorter) Len() int {
	return len(self.entries)
}

func (self *entrySorter) Swap(i, j int) {
	self.entries[i], self.entries[j] = self.entries[j], self.entries[i]
}

func (self *entrySorter) Less(i, j int) bool {
	return self.entries[i].Less(self.entries[j], self.dimension)
}
