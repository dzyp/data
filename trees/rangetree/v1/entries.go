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

func (self Entries) RemoveAt(i int) []r.Entry {
	if i >= len(self) { // this can't happen
		return self
	}

	copy(self[i:], self[i+1:])
	self[len(self)-1] = nil
	self = self[:len(self)-1]

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
